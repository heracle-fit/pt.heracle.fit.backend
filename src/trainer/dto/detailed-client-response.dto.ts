import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { ClientResponseDto } from './client-response.dto';

export class DetailedClientResponseDto extends ClientResponseDto {
    @ApiPropertyOptional()
    age?: number;

    @ApiPropertyOptional()
    gender?: string;

    @ApiPropertyOptional()
    heightCm?: number;

    @ApiPropertyOptional()
    weightKg?: number;

    @ApiPropertyOptional()
    bodyType?: string;

    @ApiPropertyOptional()
    fitnessLevel?: string;

    @ApiPropertyOptional()
    bmi?: number;

    @ApiPropertyOptional()
    targetCalories?: number;

    @ApiPropertyOptional()
    targetProtein?: number;

    @ApiPropertyOptional()
    targetCarbs?: number;

    @ApiPropertyOptional()
    targetFat?: number;

    @ApiPropertyOptional()
    targetFiber?: number;

    @ApiPropertyOptional()
    injuries?: string;

    @ApiPropertyOptional()
    dietaryPreference?: string;

    @ApiPropertyOptional()
    workoutFrequencyPerWeek?: number;

    @ApiPropertyOptional()
    preferredWorkoutType?: string;
}
