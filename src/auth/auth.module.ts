import { Module } from '@nestjs/common';
import { PassportModule } from '@nestjs/passport';
import { JwtModule } from '@nestjs/jwt';
import { APP_GUARD } from '@nestjs/core';
import { AuthService } from './auth.service';
import { TrainerGuard } from '../common/guards/trainer.guard';
import { JwtStrategy } from '../common/strategys/jwt.strategy';
import { JwtAuthGuard } from '../common/guards/jwt-auth.guard';
import { AdminGuard } from '../common/guards/admin.guard';
import { AuthController } from './auth.controller';

@Module({
	imports: [

		PassportModule.register({ defaultStrategy: 'jwt' }),
		JwtModule.register({
			secret: process.env.JWT_SECRET || 'change-me',
			signOptions: { expiresIn: '7d' },
		}),
	],
	providers: [
		AuthService,
		JwtStrategy,

		{ provide: APP_GUARD, useClass: JwtAuthGuard }, // enforces JWT by default except @Public
		{ provide: APP_GUARD, useClass: AdminGuard }, // enforces @Admin where applied
		{ provide: APP_GUARD, useClass: TrainerGuard }, // enforces @Trainer where applied
	],
	controllers: [AuthController],
	exports: [AuthService],
})

export class AuthModule { }
