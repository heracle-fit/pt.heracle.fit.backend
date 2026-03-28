import { Injectable, NotFoundException, ForbiddenException } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { UpdateSplitDto } from './dto/update-split.dto';
import { SplitResponseDto } from './dto/split-response.dto';

@Injectable()
export class SplitService {
    constructor(private readonly prisma: PrismaService) {}

    async getSplit(userId: string): Promise<SplitResponseDto> {
        const split = await this.prisma.workoutSplit.findUnique({
            where: { userId },
        });

        if (!split) {
            throw new NotFoundException(`Workout split not found for user ${userId}`);
        }

        return split;
    }

    async trainerGetSplit(trainerUserId: string, clientId: string): Promise<SplitResponseDto> {
        await this.verifyTrainerClient(trainerUserId, clientId);
        return this.getSplit(clientId);
    }

    async upsertSplit(
        trainerUserId: string,
        clientId: string,
        dto: UpdateSplitDto,
    ): Promise<SplitResponseDto> {
        const trainer = await this.verifyTrainerClient(trainerUserId, clientId);

        return this.prisma.workoutSplit.upsert({
            where: { userId: clientId },
            create: {
                userId: clientId,
                trainerId: trainer.id,
                splitData: dto.splitData,
            },
            update: {
                trainerId: trainer.id,
                splitData: dto.splitData,
                updatedAt: new Date(),
            },
        });
    }

    private async verifyTrainerClient(trainerUserId: string, clientId: string) {
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

        return trainer;
    }
}
