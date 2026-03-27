import { Module } from '@nestjs/common';
import { TrainerService } from './trainer.service';
import { TrainerController } from './trainer.controller';
import { PrismaModule } from '../prisma/prisma.module';

@Module({
    imports: [PrismaModule],
    controllers: [TrainerController],
    providers: [TrainerService],
    exports: [TrainerService],
})
export class TrainerModule {}
