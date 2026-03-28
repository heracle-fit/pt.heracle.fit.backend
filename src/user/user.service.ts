import { Injectable, ForbiddenException, NotFoundException } from '@nestjs/common';

import { PrismaService } from '../prisma/prisma.service';
import { SaveBodyMetricsDto } from './dto/swagger/body-metrics.dto';
import { SaveTargetsDto } from './dto/swagger/targets.dto';

@Injectable()
export class UserService {
	constructor(private readonly prisma: PrismaService) { }

	async getProfile(userId: string) {
		return this.prisma.user.findUnique({
			where: { id: userId },
			select: {
				id: true,
				username: true,
				name: true,
				email: true,
				avatarUrl: true,
				bio: true,
				createdAt: true,
				updatedAt: true,
			},
		});
	}

	async getBodyMetrics(userId: string) {
		return this.prisma.userProfile.findUnique({
			where: { userId },
			select: {
				age: true,
				gender: true,
				heightCm: true,
				heightFt: true,
				weightKg: true,
				weightLbs: true,
				bodyType: true,
				goal: true,
				bmi: true,
				maintenanceCalories: true,
				goalWeightKg: true,
				goalWeightLbs: true,
				targetCalories: true,
				targetProtein: true,
				targetCarbs: true,
				targetFat: true,
				targetFiber: true,
				updatedAt: true,
			},
		});
	}

	async getOnboardingStatus(userId: string) {
		const profile = await this.prisma.userProfile.findUnique({
			where: { userId },
		});

		const bodyMetricsNeeded =
			!profile?.age ||
			!profile?.gender ||
			!profile?.heightCm ||
			!profile?.weightKg ||
			!profile?.goal;

		const dietDataNeeded = !profile?.dietaryPreference || !profile?.mealsPerDay;

		return {
			bodyMetricsNeeded,
			dietDataNeeded,
		};
	}

	async saveBodyMetrics(userId: string, dto: SaveBodyMetricsDto) {
		const existingProfile = await this.prisma.userProfile.findUnique({
			where: { userId },
		});

		const age = dto.age ?? existingProfile?.age;
		const gender = dto.gender ?? existingProfile?.gender;
		const heightCm = dto.heightCm ?? existingProfile?.heightCm;
		const weightKg = dto.weightKg ?? existingProfile?.weightKg;
		const goal = dto.goal ?? existingProfile?.goal;

		let bmi: number | undefined;
		let maintenanceCalories: number | undefined;
		let targetCalories: number | undefined;
		let targetProtein: number | undefined;
		let targetCarbs: number | undefined;
		let targetFat: number | undefined;
		let targetFiber: number | undefined;

		if (heightCm && weightKg) {
			bmi = weightKg / Math.pow(heightCm / 100, 2);
		}

		if (age && gender && heightCm && weightKg) {
			// Mifflin-St Jeor Equation
			const bmr =
				10 * weightKg +
				6.25 * heightCm -
				5 * age +
				(gender.toLowerCase() === 'male' ? 5 : -161);

			// Activity Factor (Defaulting to Sedentary = 1.2 for now)
			maintenanceCalories = Math.round(bmr * 1.2);

			// Calculate Target Calories based on Goal
			if (goal === 'weight_loss') {
				targetCalories = maintenanceCalories - 500;
			} else if (goal === 'muscle_gain') {
				targetCalories = maintenanceCalories + 300;
			} else {
				targetCalories = maintenanceCalories;
			}

			// Macro Calculations
			// Protein: 2.0g per kg of body weight
			targetProtein = Number((weightKg * 2.0).toFixed(1));

			// Fat: 0.7g per kg of body weight
			targetFat = Number((weightKg * 0.7).toFixed(1));

			// Carbs: Remaining calories
			// Protein = 4 kcal/g, Fat = 9 kcal/g, Carbs = 4 kcal/g
			const proteinCalories = targetProtein * 4;
			const fatCalories = targetFat * 9;
			const remainingCalories = targetCalories - proteinCalories - fatCalories;
			targetCarbs = Number((Math.max(0, remainingCalories) / 4).toFixed(1));

			// Fiber: 14g per 1000 calories
			targetFiber = Number(((targetCalories / 1000) * 14).toFixed(1));
		}

		return this.prisma.userProfile.upsert({
			where: { userId },
			create: {
				userId,
				age: dto.age,
				gender: dto.gender,
				heightCm: dto.heightCm,
				heightFt: dto.heightFt,
				weightKg: dto.weightKg,
				weightLbs: dto.weightLbs,
				bodyType: dto.bodyType,
				goal: dto.goal,
				goalWeightKg: dto.goalWeightKg,
				goalWeightLbs: dto.goalWeightLbs,
				bmi,
				maintenanceCalories,
				targetCalories,
				targetProtein,
				targetCarbs,
				targetFat,
				targetFiber,
			},
			update: {
				...(dto.age !== undefined && { age: dto.age }),
				...(dto.gender !== undefined && { gender: dto.gender }),
				...(dto.heightCm !== undefined && { heightCm: dto.heightCm }),
				...(dto.heightFt !== undefined && { heightFt: dto.heightFt }),
				...(dto.weightKg !== undefined && { weightKg: dto.weightKg }),
				...(dto.weightLbs !== undefined && { weightLbs: dto.weightLbs }),
				...(dto.bodyType !== undefined && { bodyType: dto.bodyType }),
				...(dto.goal !== undefined && { goal: dto.goal }),
				...(dto.goalWeightKg !== undefined && { goalWeightKg: dto.goalWeightKg }),
				...(dto.goalWeightLbs !== undefined && { goalWeightLbs: dto.goalWeightLbs }),
				...(bmi !== undefined && { bmi }),
				...(maintenanceCalories !== undefined && { maintenanceCalories }),
				...(targetCalories !== undefined && { targetCalories }),
				...(targetProtein !== undefined && { targetProtein }),
				...(targetCarbs !== undefined && { targetCarbs }),
				...(targetFat !== undefined && { targetFat }),
				...(targetFiber !== undefined && { targetFiber }),
				updatedAt: new Date(),
			},
			select: {
				age: true,
				gender: true,
				heightCm: true,
				heightFt: true,
				weightKg: true,
				weightLbs: true,
				bodyType: true,
				goal: true,
				bmi: true,
				maintenanceCalories: true,
				goalWeightKg: true,
				goalWeightLbs: true,
				targetCalories: true,
				targetProtein: true,
				targetCarbs: true,
				targetFat: true,
				targetFiber: true,
				updatedAt: true,
			},
		});
	}

	async saveTargets(userId: string, dto: SaveTargetsDto) {
		return this.prisma.userProfile.upsert({
			where: { userId },
			create: {
				userId,
				targetCalories: dto.targetCalories,
				targetProtein: dto.targetProtein,
				targetCarbs: dto.targetCarbs,
				targetFat: dto.targetFat,
				targetFiber: dto.targetFiber,
			},
			update: {
				...(dto.targetCalories !== undefined && { targetCalories: dto.targetCalories }),
				...(dto.targetProtein !== undefined && { targetProtein: dto.targetProtein }),
				...(dto.targetCarbs !== undefined && { targetCarbs: dto.targetCarbs }),
				...(dto.targetFat !== undefined && { targetFat: dto.targetFat }),
				...(dto.targetFiber !== undefined && { targetFiber: dto.targetFiber }),
				updatedAt: new Date(),
			},
			select: {
				targetCalories: true,
				targetProtein: true,
				targetCarbs: true,
				targetFat: true,
				targetFiber: true,
				updatedAt: true,
			},
		});
	}

	async trainerSaveBodyMetrics(trainerUserId: string, clientId: string, dto: SaveBodyMetricsDto) {

		// 1. Verify trainer-client relationship
		const trainer = await this.prisma.trainer.findUnique({
			where: { userId: trainerUserId },
		});

		if (!trainer) {
			throw new ForbiddenException('Trainer record not found for this user');
		}

		const assignment = await this.prisma.trainerClient.findUnique({
			where: { clientId },
		});

		if (!assignment || assignment.trainerId !== trainer.id) {
			throw new ForbiddenException('You are not assigned to this client');
		}

		// 2. Delegate to existing saveBodyMetrics logic
		return this.saveBodyMetrics(clientId, dto);
	}
}


