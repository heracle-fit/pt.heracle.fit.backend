import { Body, Controller, Delete, Get, Param, Patch, Post, Req } from '@nestjs/common';
import { ApiBearerAuth, ApiBody, ApiOkResponse, ApiOperation, ApiTags } from '@nestjs/swagger';
import { Trainer } from '../common/decorators/trainer.decorator';

import { WorkoutService } from './workout.service';
import { TodayWorkoutResponseDto, WorkoutSessionDto } from './dto/today-workout.dto';
import { SaveWorkoutPreferencesDto, WorkoutPreferencesResponseDto } from './dto/workout-preferences.dto';
import { ExerciseListResponseDto } from './dto/exercise-list.dto';
import { CreateSessionRequestDto } from './dto/create-session-request.dto';
import { UpdateSessionRequestDto } from './dto/update-session-request.dto';
import { SessionResponseDto } from './dto/session-response.dto';
import { CreateWorkoutLogRequestDto } from './dto/create-workout-log-request.dto';
import { UpdateWorkoutLogRequestDto } from './dto/update-workout-log-request.dto';
import { WorkoutLogResponseDto } from './dto/workout-log-response.dto';

@ApiTags('Workout')
@ApiBearerAuth('JWT')
@Controller('workout')
export class WorkoutController {
    constructor(private readonly workoutService: WorkoutService) { }

    @Get('today')
    @ApiOperation({
        summary: "Get today's suggested workout",
        description: 'Returns the recommended workout plan for today. Returns null if workout preferences are not set.',
    })
    @ApiOkResponse({ type: TodayWorkoutResponseDto })
    async getTodayWorkout(@Req() req: any): Promise<TodayWorkoutResponseDto | null> {
        return this.workoutService.getTodayWorkout(req.user.id);
    }

    @Get('sessions')
    @ApiOperation({
        summary: 'Get workout sessions',
        description: 'Returns all available workout sessions with their exercises',
    })
    @ApiOkResponse({ type: [WorkoutSessionDto] })
    getSessions(): WorkoutSessionDto[] {
        return this.workoutService.getSessions();
    }

    @Get('exercises')
    @ApiOperation({
        summary: 'Get all exercises',
        description: 'Returns a list of all available exercises from the database',
    })
    @ApiOkResponse({ type: [ExerciseListResponseDto] })
    async getExercises(): Promise<ExerciseListResponseDto[]> {
        return this.workoutService.getExercises();
    }

    @Get('preferences')
    @ApiOperation({
        summary: 'Get workout preferences',
        description: 'Returns the stored workout preferences for the authenticated user',
    })
    @ApiOkResponse({ type: WorkoutPreferencesResponseDto })
    async getWorkoutPreferences(@Req() req: any) {
        return this.workoutService.getWorkoutPreferences(req.user.id);
    }

    @Post('preferences')
    @ApiOperation({
        summary: 'Save workout preferences',
        description: 'Creates or updates workout preferences. All fields are optional — send only what you want to change.',
    })
    @ApiBody({ type: SaveWorkoutPreferencesDto })
    @ApiOkResponse({ type: WorkoutPreferencesResponseDto })
    async saveWorkoutPreferences(@Req() req: any, @Body() body: SaveWorkoutPreferencesDto) {
        return this.workoutService.saveWorkoutPreferences(req.user.id, body);
    }

    // --- Session CRUD ---

    @Post('session')
    @ApiOperation({ summary: 'Create a new workout session' })
    @ApiBody({ type: CreateSessionRequestDto })
    @ApiOkResponse({ type: SessionResponseDto })
    async createSession(@Req() req: any, @Body() body: CreateSessionRequestDto): Promise<SessionResponseDto> {
        return this.workoutService.createSession(req.user.id, body);
    }

    @Get('session/:id')
    @ApiOperation({ summary: 'Get a specific workout session' })
    @ApiOkResponse({ type: SessionResponseDto })
    async getSession(@Req() req: any, @Param('id') id: string): Promise<SessionResponseDto> {
        return this.workoutService.getSession(req.user.id, +id);
    }

    @Get('my-sessions')
    @ApiOperation({ summary: 'Get all workout sessions for the authenticated user' })
    @ApiOkResponse({ type: [SessionResponseDto] })
    async getUserSessions(@Req() req: any): Promise<SessionResponseDto[]> {
        return this.workoutService.getUserSessions(req.user.id);
    }

    @Patch('session/:id')
    @ApiOperation({ summary: 'Update a workout session' })
    @ApiBody({ type: UpdateSessionRequestDto })
    @ApiOkResponse({ type: SessionResponseDto })
    async updateSession(
        @Req() req: any,
        @Param('id') id: string,
        @Body() body: UpdateSessionRequestDto,
    ): Promise<SessionResponseDto> {
        return this.workoutService.updateSession(req.user.id, +id, body);
    }

    @Delete('session/:id')
    @ApiOperation({ summary: 'Delete a workout session' })
    @ApiOkResponse({ description: 'Session deleted successfully' })
    async deleteSession(@Req() req: any, @Param('id') id: string): Promise<void> {
        return this.workoutService.deleteSession(req.user.id, +id);
    }

    @Patch('trainer/session/:clientId/:sessionId')
    @Trainer()
    @ApiOperation({ summary: 'Edit a client session data (Trainer Only)' })
    @ApiBody({ type: UpdateSessionRequestDto })
    @ApiOkResponse({ type: SessionResponseDto })
    async trainerUpdateSession(
        @Req() req: any,
        @Param('clientId') clientId: string,
        @Param('sessionId') sessionId: string,
        @Body() body: UpdateSessionRequestDto,
    ): Promise<SessionResponseDto> {
        return this.workoutService.trainerUpdateSession(req.user.id, clientId, +sessionId, body);
    }

    @Post('trainer/session/:clientId')
    @Trainer()
    @ApiOperation({ summary: 'Create a session for a client (Trainer Only)' })
    @ApiBody({ type: CreateSessionRequestDto })
    @ApiOkResponse({ type: SessionResponseDto })
    async trainerCreateSession(
        @Req() req: any,
        @Param('clientId') clientId: string,
        @Body() body: CreateSessionRequestDto,
    ): Promise<SessionResponseDto> {
        return this.workoutService.trainerCreateSession(req.user.id, clientId, body);
    }

    @Get('trainer/sessions/:clientId')
    @Trainer()
    @ApiOperation({ summary: "Get all sessions for a client (Trainer Only)" })
    @ApiOkResponse({ type: [SessionResponseDto] })
    async trainerGetSessions(
        @Req() req: any,
        @Param('clientId') clientId: string,
    ): Promise<SessionResponseDto[]> {
        return this.workoutService.trainerGetSessions(req.user.id, clientId);
    }

    @Delete('trainer/session/:clientId/:sessionId')
    @Trainer()
    @ApiOperation({ summary: 'Delete a client session (Trainer Only)' })
    @ApiOkResponse({ description: 'Session deleted successfully' })
    async trainerDeleteSession(
        @Req() req: any,
        @Param('clientId') clientId: string,
        @Param('sessionId') sessionId: string,
    ): Promise<void> {
        return this.workoutService.trainerDeleteSession(req.user.id, clientId, +sessionId);
    }



    // --- WorkoutLog CRUD ---

    @Post('log')
    @ApiOperation({ summary: 'Create a new workout log' })
    @ApiBody({ type: CreateWorkoutLogRequestDto })
    @ApiOkResponse({ type: WorkoutLogResponseDto })
    async createWorkoutLog(@Req() req: any, @Body() body: CreateWorkoutLogRequestDto): Promise<WorkoutLogResponseDto> {
        return this.workoutService.createWorkoutLog(req.user.id, body);
    }

    @Get('log/:id')
    @ApiOperation({ summary: 'Get a specific workout log' })
    @ApiOkResponse({ type: WorkoutLogResponseDto })
    async getWorkoutLog(@Req() req: any, @Param('id') id: string): Promise<WorkoutLogResponseDto> {
        return this.workoutService.getWorkoutLog(req.user.id, +id);
    }

    @Get('logs')
    @ApiOperation({ summary: 'Get all workout logs for the authenticated user' })
    @ApiOkResponse({ type: [WorkoutLogResponseDto] })
    async getWorkoutLogs(@Req() req: any): Promise<WorkoutLogResponseDto[]> {
        return this.workoutService.getWorkoutLogs(req.user.id);
    }

    @Patch('log/:id')
    @ApiOperation({ summary: 'Update a workout log' })
    @ApiBody({ type: UpdateWorkoutLogRequestDto })
    @ApiOkResponse({ type: WorkoutLogResponseDto })
    async updateWorkoutLog(
        @Req() req: any,
        @Param('id') id: string,
        @Body() body: UpdateWorkoutLogRequestDto,
    ): Promise<WorkoutLogResponseDto> {
        return this.workoutService.updateWorkoutLog(req.user.id, +id, body);
    }

    @Delete('log/:id')
    @ApiOperation({ summary: 'Delete a workout log' })
    @ApiOkResponse({ description: 'Workout log deleted successfully' })
    async deleteWorkoutLog(@Req() req: any, @Param('id') id: string): Promise<void> {
        return this.workoutService.deleteWorkoutLog(req.user.id, +id);
    }

    @Get('trainer/logs/:clientId')
    @Trainer()
    @ApiOperation({ summary: 'Get all workout logs for a client (Trainer Only)' })
    @ApiOkResponse({ type: [WorkoutLogResponseDto] })
    async trainerGetWorkoutLogs(
        @Req() req: any,
        @Param('clientId') clientId: string,
    ): Promise<WorkoutLogResponseDto[]> {
        return this.workoutService.trainerGetWorkoutLogs(req.user.id, clientId);
    }


    @Patch('trainer/log-review/:logId')
    @Trainer()
    @ApiOperation({ summary: 'Add a PT review to a client workout log (Trainer Only)' })
    @ApiBody({
        schema: {
            type: 'object',
            properties: {
                review: { type: 'string' },
            },
            required: ['review'],
        },
    })
    @ApiOkResponse({ type: WorkoutLogResponseDto })
    async addWorkoutLogReview(
        @Req() req: any,
        @Param('logId') logId: string,
        @Body('review') review: string,
    ): Promise<WorkoutLogResponseDto> {
        return this.workoutService.trainerAddWorkoutLogReview(req.user.id, +logId, review);
    }
}



