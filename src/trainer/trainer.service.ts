import { Injectable, NotFoundException, ConflictException, ForbiddenException } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { ClientResponseDto } from './dto/client-response.dto';

@Injectable()
export class TrainerService {
    constructor(private readonly prisma: PrismaService) {}

    async getClients(trainerUserId: string): Promise<ClientResponseDto[]> {
        const trainer = await this.getTrainer(trainerUserId);

        const assignments = await this.prisma.trainerClient.findMany({
            where: { trainerId: trainer.id },
            include: {
                client: true,
            },
        });

        return assignments.map(a => ({
            id: a.client.id,
            name: a.client.name,
            email: a.client.email,
            avatarUrl: a.client.avatarUrl,
            assignedAt: a.assignedAt,
        }));
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

        const assignment = await this.prisma.trainerClient.create({
            data: {
                trainerId: trainer.id,
                clientId: clientUser.id,
            },
            include: {
                client: true,
            },
        });

        return {
            id: assignment.client.id,
            name: assignment.client.name,
            email: assignment.client.email,
            avatarUrl: assignment.client.avatarUrl,
            assignedAt: assignment.assignedAt,
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
