import { Injectable, ForbiddenException } from '@nestjs/common';
import { CanActivate, ExecutionContext } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { IS_TRAINER_KEY } from '../decorators/trainer.decorator';
import { IS_PUBLIC_KEY } from '../decorators/public.decorator';

@Injectable()
export class TrainerGuard implements CanActivate {
	constructor(private reflector: Reflector) { }

	canActivate(context: ExecutionContext): boolean {
		// If route is public, allow
		const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
			context.getHandler(),
			context.getClass(),
		]);
		if (isPublic) return true;

		// If @Trainer not present, allow (normal JWT guard will handle auth)
		const requiresTrainer = this.reflector.getAllAndOverride<boolean>(IS_TRAINER_KEY, [
			context.getHandler(),
			context.getClass(),
		]);
		if (!requiresTrainer) return true;

		const req = context.switchToHttp().getRequest();
		const user = req.user;

		if (!user || (user.role !== 'trainer' && user.role !== 'admin')) {
			throw new ForbiddenException('Trainer privileges required');
		}

		return true;
	}
}
