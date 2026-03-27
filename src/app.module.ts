import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { PrismaModule } from './prisma/prisma.module';
import { UserController } from './user/user.controller';
import { UserService } from './user/user.service';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { FirebaseModule } from './firebase/firebase.module';
import { WorkoutModule } from './workout/workout.module';
import { DietModule } from './diet/diet.module';
import { SleepModule } from './sleep/sleep.module';
import { TrainerModule } from './trainer/trainer.module';

@Module({
  imports: [
    FirebaseModule,
    PrismaModule,
    UserModule,
    AuthModule,
    WorkoutModule,
    DietModule,
    SleepModule,
    TrainerModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule { }

