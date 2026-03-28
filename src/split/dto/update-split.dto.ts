import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty } from 'class-validator';

export class UpdateSplitDto {
    @ApiProperty({
        example: [
            { day: 'monday', sessionId: 1 },
            { day: 'wednesday', sessionId: 2 },
            { day: 'friday', sessionId: 1 }
        ],
        description: 'Array of day-to-session mappings'
    })
    @IsNotEmpty()
    splitData: any;
}
