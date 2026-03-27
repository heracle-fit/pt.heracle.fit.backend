import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class WorkoutLogResponseDto {
    @ApiProperty({ example: 1 })
    id: number;

    @ApiProperty({ example: 'uuid-...' })
    userId: string;

    @ApiPropertyOptional({ example: 1 })
    sessionId: number | null;

    @ApiProperty()
    logData: any;

    @ApiPropertyOptional()
    notes: string | null;

    @ApiProperty({ example: 'https://r2.heracle.fit/exercises' })
    exerciseImageBaseUrl: string;

    @ApiProperty()
    createdAt: Date;

    @ApiPropertyOptional()
    ptReview: string | null;

    @ApiProperty()
    updatedAt: Date;

}
