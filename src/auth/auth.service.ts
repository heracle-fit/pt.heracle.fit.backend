import { Injectable, InternalServerErrorException, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from '../prisma/prisma.service';
import * as admin from 'firebase-admin';

@Injectable()
export class AuthService {
	constructor(
		private readonly jwtService: JwtService,
		private readonly prisma: PrismaService
	) { }

	// Admin login using environment credentials
	async adminLogin(dto: any) {
		const adminUser = process.env.ADMIN_USERNAME || 'admin';
		const adminPass = process.env.ADMIN_PASSWORD || 'admin_password_123';

		if (dto.username === adminUser && dto.password === adminPass) {
			const payload = { sub: 'admin-id', username: adminUser, role: 'admin' };
			const token = this.jwtService.sign(payload);
			return { user: { id: 'admin-id', username: adminUser, role: 'admin' }, token };
		}
		throw new UnauthorizedException('Invalid admin credentials');
	}

	// Find existing user by email or create a new one, then return a jwt + user
	async validateOAuthLogin(profile: any) {
		if (!profile || !profile.emails || !profile.emails.length) {
			throw new InternalServerErrorException('No email found from Google OAuth');
		}

		const email: string = profile.emails[0].value;
		let user = await this.prisma.user.findUnique({ where: { email } });

		if (!user) {
			// create username from local-part of email, ensure uniqueness by appending random suffix if needed
			let baseUsername = email.split('@')[0].replace(/[^a-zA-Z0-9._-]/g, '').slice(0, 30) || 'user';
			let username = baseUsername;

			const existingUsers = await this.prisma.user.findMany({
				where: { username: { startsWith: baseUsername } },
				select: { username: true }
			});
			const existingUsernames = new Set(existingUsers.map(u => u.username));

			let suffix = 0;
			while (existingUsernames.has(username)) {
				suffix += 1;
				username = `${baseUsername}${suffix}`;
			}

			user = await this.prisma.user.create({
				data: {
					username,
					name: profile.displayName || username,
					email,
					avatarUrl: profile.photos?.[0]?.value || null,
					bio: null,
					googleAccessToken: profile.accessToken || null,
				},
			});
		} else if (profile.accessToken && user.googleAccessToken !== profile.accessToken) {
			// Update the access token if it's new or changed
			user = await this.prisma.user.update({
				where: { id: user.id },
				data: { googleAccessToken: profile.accessToken },
			});
		}

		// Determine role based on Trainer table
		const trainer = await this.prisma.trainer.findUnique({ where: { userId: user.id } });
		const role = trainer ? 'trainer' : 'user';

		const payload = { sub: user.id, email: user.email, role: role };
		const token = this.jwtService.sign(payload);

		return { user: { ...user, role }, token };
	}

	// Verify a Firebase ID token (issued by Firebase Auth on the mobile client)
	// and return { user, token } where token is our own JWT
	async validateIdToken(idToken: string, accessToken?: string) {
		if (!idToken) {
			throw new UnauthorizedException('Missing idToken');
		}

		let decodedToken: admin.auth.DecodedIdToken;
		try {
			decodedToken = await admin.auth().verifyIdToken(idToken);
		} catch (err) {
			throw new UnauthorizedException('Invalid Firebase ID token');
		}

		if (!decodedToken.email) {
			throw new UnauthorizedException('Firebase token does not contain an email');
		}

		// Build a profile-like object that validateOAuthLogin expects
		const profileLike = {
			emails: [{ value: decodedToken.email }],
			displayName: decodedToken.name || decodedToken.email.split('@')[0],
			photos: decodedToken.picture ? [{ value: decodedToken.picture }] : [],
			accessToken,
		};

		// Reuse existing logic to find-or-create user and issue our JWT
		return this.validateOAuthLogin(profileLike);
	}


	// Developer-only method to generate a token for an email without Firebase verification
	async generateDevToken(email: string) {
		const profileLike = {
			emails: [{ value: email }],
			displayName: email.split('@')[0],
			photos: [],
		};

		return this.validateOAuthLogin(profileLike);
	}
}
