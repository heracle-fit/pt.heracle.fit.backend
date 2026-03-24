import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { TodayWorkoutResponseDto, WorkoutSessionDto } from './dto/today-workout.dto';
import { SaveWorkoutPreferencesDto, WorkoutPreferencesResponseDto } from './dto/workout-preferences.dto';
import { ExerciseListResponseDto } from './dto/exercise-list.dto';
import { CreateSessionRequestDto } from './dto/create-session-request.dto';
import { UpdateSessionRequestDto } from './dto/update-session-request.dto';
import { SessionResponseDto } from './dto/session-response.dto';
import { NotFoundException } from '@nestjs/common';

@Injectable()
export class WorkoutService {
    constructor(private readonly prisma: PrismaService) { }


    async getTodayWorkout(userId: string): Promise<TodayWorkoutResponseDto | null> {
        // Check required preference fields — return null if onboarding is incomplete
        const prefs = await this.prisma.userProfile.findUnique({
            where: { userId },
            select: {
                fitnessLevel: true,
                workoutFrequencyPerWeek: true,
                preferredWorkoutType: true,
                availableDays: true,
                preferredWorkoutTime: true,
            },
        });

        const ready =
            prefs &&
            prefs.fitnessLevel &&
            prefs.workoutFrequencyPerWeek &&
            prefs.preferredWorkoutType &&
            prefs.availableDays?.length > 0 &&
            prefs.preferredWorkoutTime;

        if (!ready) return null;

        return {
            title: 'Suggested Muscle',
            highlight: 'Bicep & Back',
            subtext: 'Optional for hypertrophy',
            duration: 45,
            intensity: 'hard',
            session: this.getSessions(),
        };
    }

    getSessions(): WorkoutSessionDto[] {
        return [
            {
                id: 'session-001',
                title: 'Chest & Triceps',
                content: 'Push day focused on upper body strength',
                category: 'Strength',
                exercisesCount: 3,
                position: 1,
                exercises: [
                    {
                        id: 'ex-001',
                        name: 'Bench Press',
                        desc: 'Chest, Triceps, Shoulders',
                        image: 'https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg',
                        sets: [
                            { kg: 60, reps: 10 },
                            { kg: 65, reps: 8 },
                            { kg: 70, reps: 6 },
                        ],
                    },
                    {
                        id: 'ex-002',
                        name: 'Tricep Dips',
                        desc: 'Triceps, Chest',
                        image: 'https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg',
                        sets: [
                            { kg: 0, reps: 12 },
                            { kg: 0, reps: 10 },
                            { kg: 0, reps: 10 },
                        ],
                    },
                    {
                        id: 'ex-003',
                        name: 'Incline Dumbbell Press',
                        desc: 'Upper Chest, Shoulders',
                        image: 'https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg',
                        sets: [
                            { kg: 24, reps: 10 },
                            { kg: 26, reps: 8 },
                            { kg: 28, reps: 6 },
                        ],
                    },
                ],
            },
            {
                id: 'session-002',
                title: 'Back & Biceps',
                content: 'Pull day targeting back width and bicep peak',
                category: 'Strength',
                exercisesCount: 3,
                position: 2,
                exercises: [
                    {
                        id: 'ex-004',
                        name: 'Deadlift',
                        desc: 'Lower Back, Glutes, Hamstrings',
                        image: 'https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg',
                        sets: [
                            { kg: 80, reps: 8 },
                            { kg: 90, reps: 6 },
                            { kg: 100, reps: 4 },
                        ],
                    },
                    {
                        id: 'ex-005',
                        name: 'Pull-Ups',
                        desc: 'Lats, Biceps, Rear Delts',
                        image: 'https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg',
                        sets: [
                            { kg: 0, reps: 10 },
                            { kg: 0, reps: 8 },
                            { kg: 0, reps: 8 },
                        ],
                    },
                    {
                        id: 'ex-006',
                        name: 'Barbell Curl',
                        desc: 'Biceps, Forearms',
                        image: 'https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg',
                        sets: [
                            { kg: 30, reps: 12 },
                            { kg: 35, reps: 10 },
                            { kg: 35, reps: 8 },
                        ],
                    },
                ],
            },
        ];
    }

    async getExercises(): Promise<ExerciseListResponseDto[]> {
        const exercises = await this.prisma.exercise.findMany({
            orderBy: { name: 'asc' },
        });

        return exercises.map((ex) => ({
            id: ex.id,
            name: ex.name,
            secondaryInfo: ex.secondaryInfo,
            exerciseType: ex.exerciseType,
            image: `https://r2.heracle.fit/exercises/${ex.id}.jpg`,
        }));
    }

    async getWorkoutPreferences(userId: string): Promise<WorkoutPreferencesResponseDto | null> {
        const profile = await this.prisma.userProfile.findUnique({
            where: { userId },
            select: {
                id: true,
                fitnessLevel: true,
                workoutFrequencyPerWeek: true,
                preferredWorkoutType: true,
                injuries: true,
                availableDays: true,
                preferredWorkoutTime: true,
                sessionDurationMins: true,
                updatedAt: true,
            },
        });
        return profile;
    }

    async saveWorkoutPreferences(
        userId: string,
        dto: SaveWorkoutPreferencesDto,
    ): Promise<WorkoutPreferencesResponseDto> {
        const profile = await this.prisma.userProfile.upsert({
            where: { userId },
            create: {
                userId,
                fitnessLevel: dto.fitnessLevel,
                workoutFrequencyPerWeek: dto.workoutFrequencyPerWeek,
                preferredWorkoutType: dto.preferredWorkoutType,
                injuries: dto.injuries,
                availableDays: dto.availableDays ?? [],
                preferredWorkoutTime: dto.preferredWorkoutTime,
                sessionDurationMins: dto.sessionDurationMins,
            },
            update: {
                ...(dto.fitnessLevel !== undefined && { fitnessLevel: dto.fitnessLevel }),
                ...(dto.workoutFrequencyPerWeek !== undefined && { workoutFrequencyPerWeek: dto.workoutFrequencyPerWeek }),
                ...(dto.preferredWorkoutType !== undefined && { preferredWorkoutType: dto.preferredWorkoutType }),
                ...(dto.injuries !== undefined && { injuries: dto.injuries }),
                ...(dto.availableDays !== undefined && { availableDays: dto.availableDays }),
                ...(dto.preferredWorkoutTime !== undefined && { preferredWorkoutTime: dto.preferredWorkoutTime }),
                ...(dto.sessionDurationMins !== undefined && { sessionDurationMins: dto.sessionDurationMins }),
                updatedAt: new Date(),
            },
            select: {
                id: true,
                fitnessLevel: true,
                workoutFrequencyPerWeek: true,
                preferredWorkoutType: true,
                injuries: true,
                availableDays: true,
                preferredWorkoutTime: true,
                sessionDurationMins: true,
                updatedAt: true,
            },
        });
        return profile;
    }

    // --- Session CRUD ---

    async createSession(userId: string, dto: CreateSessionRequestDto): Promise<SessionResponseDto> {
        return this.prisma.session.create({
            data: {
                userId,
                name: dto.name,
                category: dto.category,
                sessionData: dto.sessionData,
            },
            select: {
                id: true,
                name: true,
                category: true,
                sessionData: true,
                createdAt: true,
                updatedAt: true,
            },
        }) as any;
    }

    async getSession(userId: string, id: number): Promise<SessionResponseDto> {
        const session = await this.prisma.session.findFirst({
            where: { id, userId },
            select: {
                id: true,
                name: true,
                category: true,
                sessionData: true,
                createdAt: true,
                updatedAt: true,
            },
        });

        if (!session) {
            throw new NotFoundException(`Session with ID ${id} not found`);
        }

        return session as any;
    }

    async getUserSessions(userId: string): Promise<SessionResponseDto[]> {
        return this.prisma.session.findMany({
            where: { userId },
            orderBy: { createdAt: 'desc' },
            select: {
                id: true,
                name: true,
                category: true,
                sessionData: true,
                createdAt: true,
                updatedAt: true,
            },
        }) as any;
    }

    async updateSession(userId: string, id: number, dto: UpdateSessionRequestDto): Promise<SessionResponseDto> {
        // Ensure ownership
        await this.getSession(userId, id);

        return this.prisma.session.update({
            where: { id },
            data: {
                ...(dto.name !== undefined && { name: dto.name }),
                ...(dto.category !== undefined && { category: dto.category }),
                ...(dto.sessionData !== undefined && { sessionData: dto.sessionData }),
                updatedAt: new Date(),
            },
            select: {
                id: true,
                name: true,
                category: true,
                sessionData: true,
                createdAt: true,
                updatedAt: true,
            },
        }) as any;
    }

    async deleteSession(userId: string, id: number): Promise<void> {
        // Ensure ownership
        await this.getSession(userId, id);

        await this.prisma.session.delete({
            where: { id },
        });
    }
}

