import { Injectable, NotFoundException, ConflictException, ForbiddenException } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { ClientResponseDto } from './dto/client-response.dto';

@Injectable()
export class TrainerService {
    constructor(private readonly prisma: PrismaService) {}

    async getClients(trainerUserId: string): Promise<ClientResponseDto[]> {
        const trainer = await this.getTrainer(trainerUserId);
        const today = new Date().toISOString().split('T')[0];

        const assignments = await this.prisma.trainerClient.findMany({
            where: { trainerId: trainer.id },
            include: {
                client: {
                    include: {
                        profile: {
                            select: {
                                goal: true,
                                targetCalories: true,
                            },
                        },
                    },
                },
            },
        });

        const clientIds = assignments.map(a => a.clientId);

        // Fetch today's meals for all clients in this trainer's list to calculate progress
        const todayMeals = await this.prisma.meal.findMany({
            where: {
                userId: { in: clientIds },
                date: today,
            },
        });

        // Calculate progress for each client
        return assignments.map(a => {
            const clientMeals = todayMeals.filter(m => m.userId === a.clientId);
            const consumedCalories = clientMeals.reduce((acc, meal) => {
                const foodItems = (meal.data as any) || [];
                return acc + foodItems.reduce((sum: number, item: any) => sum + (item.calories || 0), 0);
            }, 0);

            const targetCalories = a.client.profile?.targetCalories || 0;
            const progress = targetCalories > 0 ? Math.min(1, consumedCalories / targetCalories) : 0;

            return {
                id: a.client.id,
                name: a.client.name,
                email: a.client.email,
                avatarUrl: a.client.avatarUrl,
                assignedAt: a.assignedAt,
                goal: a.client.profile?.goal ?? null,
                progress: Number(progress.toFixed(2)),
            };
        });
    }


    async addClient(trainerUserId: string, email: string): Promise<ClientResponseDto> {
        const trainer = await this.getTrainer(trainerUserId);

        const clientUser = await this.prisma.user.findUnique({
            where: { email },
        });

        if (!clientUser) {
            throw new NotFoundException(`User with email ${email} not found`);
        }

        // Check if user is already assigned to a trainer
        const existingAssignment = await this.prisma.trainerClient.findUnique({
            where: { clientId: clientUser.id },
        });

        if (existingAssignment) {
            if (existingAssignment.trainerId === trainer.id) {
                throw new ConflictException('User is already your client');
            }
            throw new ConflictException('User is already assigned to another trainer');
        }

        const today = new Date().toISOString().split('T')[0];

        const assignment = await this.prisma.trainerClient.create({
            data: {
                trainerId: trainer.id,
                clientId: clientUser.id,
            },
            include: {
                client: {
                    include: {
                        profile: true,
                    },
                },
            },
        });

        // Calculate progress for the newly added client for today
        const clientMeals = await this.prisma.meal.findMany({
            where: {
                userId: clientUser.id,
                date: today,
            },
        });

        const consumedCalories = clientMeals.reduce((acc, meal) => {
            const foodItems = (meal.data as any) || [];
            return acc + foodItems.reduce((sum: number, item: any) => sum + (item.calories || 0), 0);
        }, 0);

        const targetCalories = assignment.client.profile?.targetCalories || 0;
        const progress = targetCalories > 0 ? Math.min(1, consumedCalories / targetCalories) : 0;

        return {
            id: assignment.client.id,
            name: assignment.client.name,
            email: assignment.client.email,
            avatarUrl: assignment.client.avatarUrl,
            assignedAt: assignment.assignedAt,
            goal: assignment.client.profile?.goal ?? null,
            progress: Number(progress.toFixed(2)),
        };
    }


    async removeClient(trainerUserId: string, clientId: string): Promise<void> {
        const trainer = await this.getTrainer(trainerUserId);

        const assignment = await this.prisma.trainerClient.findUnique({
            where: { clientId },
        });

        if (!assignment || assignment.trainerId !== trainer.id) {
            throw new ForbiddenException('User is not your client');
        }

        await this.prisma.trainerClient.delete({
            where: { clientId },
        });
    }

    async getClientDetails(trainerUserId: string, clientId: string) {
        const trainer = await this.getTrainer(trainerUserId);

        // 1. Verify client is assigned to this trainer
        const assignment = await this.prisma.trainerClient.findUnique({
            where: { clientId },
            include: {
                client: {
                    include: {
                        profile: true,
                    },
                },
            },
        });

        if (!assignment || assignment.trainerId !== trainer.id) {
            throw new ForbiddenException('Client is not assigned to you');
        }

        const { client } = assignment;
        const profile = client.profile;

        // 2. Fetch today's calorie progress
        const today = new Date().toISOString().split('T')[0];
        const clientMeals = await this.prisma.meal.findMany({
            where: { userId: clientId, date: today },
        });

        const consumedCalories = clientMeals.reduce((acc, meal) => {
            const foodItems = (meal.data as any) || [];
            return acc + foodItems.reduce((sum: number, item: any) => sum + (item.calories || 0), 0);
        }, 0);

        const targetCalories = profile?.targetCalories || 0;
        const progress = targetCalories > 0 ? Math.min(1, consumedCalories / targetCalories) : 0;

        // 3. Aggregate data into detailed DTO format
        return {
            id: client.id,
            name: client.name,
            email: client.email,
            avatarUrl: client.avatarUrl,
            assignedAt: assignment.assignedAt,
            goal: profile?.goal ?? null,
            progress: Number(progress.toFixed(2)),
            // Flattened profile data
            age: profile?.age ?? undefined,
            gender: profile?.gender ?? undefined,
            heightCm: profile?.heightCm ?? undefined,
            weightKg: profile?.weightKg ?? undefined,
            bodyType: profile?.bodyType ?? undefined,
            fitnessLevel: profile?.fitnessLevel ?? undefined,
            bmi: profile?.bmi ?? undefined,
            targetCalories: profile?.targetCalories ?? undefined,
            targetProtein: profile?.targetProtein ?? undefined,
            targetCarbs: profile?.targetCarbs ?? undefined,
            targetFat: profile?.targetFat ?? undefined,
            targetFiber: profile?.targetFiber ?? undefined,
            injuries: profile?.injuries ?? undefined,
            dietaryPreference: profile?.dietaryPreference ?? undefined,
            workoutFrequencyPerWeek: profile?.workoutFrequencyPerWeek ?? undefined,
            preferredWorkoutType: profile?.preferredWorkoutType ?? undefined,
        };
    }



    async adminAddTrainer(dto: { email: string; specialization?: string; experience?: number }) {
        const user = await this.prisma.user.findUnique({
            where: { email: dto.email },
        });

        if (!user) {
            throw new NotFoundException(`User with email ${dto.email} not found`);
        }

        const existingTrainer = await this.prisma.trainer.findUnique({
            where: { userId: user.id },
        });

        if (existingTrainer) {
            throw new ConflictException('User is already a trainer');
        }

        return this.prisma.trainer.create({
            data: {
                userId: user.id,
                specialization: dto.specialization,
                experience: dto.experience,
            },
            include: {
                user: true,
            },
        });
    }

    private async getTrainer(userId: string) {

        const trainer = await this.prisma.trainer.findUnique({
            where: { userId },
        });

        if (!trainer) {
            throw new ForbiddenException('Trainer record not found for this user');
        }

        return trainer;
    }
}
