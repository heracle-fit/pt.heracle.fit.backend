import { Controller, Get, Req, UseGuards, Post, Body, Query } from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';
import { ApiOkResponse, ApiBody, ApiTags, ApiOperation, ApiSecurity, ApiExcludeEndpoint, ApiQuery } from '@nestjs/swagger';
import { AndroidTokenResponseDto, GoogleCallbackResponseDto } from './dto/swagger/google-callback';
import { AndroidTokenDto } from './dto/authenticate-android';
import { AuthService } from './auth.service';
import { Public } from '../common/decorators/public.decorator';

@ApiTags('Auth')
@Controller('auth')
export class AuthController {
	constructor(private readonly authService: AuthService) { }


	@Post('google/token')
	@Public()
	@ApiOperation({
		summary: 'Authenticate with Google ID Token',
		description: 'For mobile clients to authenticate using Google ID token'
	})
	@ApiBody({ type: AndroidTokenDto })
	@ApiOkResponse({ type: GoogleCallbackResponseDto })
	async authenticateAndroid(@Body() body: AndroidTokenDto) {
		const result = await this.authService.validateIdToken(body.idToken, body.accessToken);
		return result;
	}

	@Post('admin/login')
	@Public()
	@ApiOperation({ summary: 'Admin login with username and password' })
	@ApiBody({
		schema: {
			type: 'object',
			properties: {
				username: { type: 'string' },
				password: { type: 'string' },
			},
		},
	})
	async adminLogin(@Body() body: any) {
		return this.authService.adminLogin(body);
	}

	@Get('dev/token')


	@Public()
	@ApiExcludeEndpoint()
	@ApiOperation({ summary: 'Generate a dev JWT (Local/Dev Only)' })
	@ApiQuery({ name: 'email', required: false, example: 'sanjaysagar.main@gmail.com' })
	async getDevToken(@Query('email') email?: string) {
		const targetEmail = email || 'sanjaysagar.main@gmail.com';
		return this.authService.generateDevToken(targetEmail);
	}


}
