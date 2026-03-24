import { ApiProperty } from '@nestjs/swagger';

export class CreateSessionRequestDto {
    @ApiProperty({ example: 'Morning Push Session', description: 'Name of the workout session' })
    name: string;

    @ApiProperty({ example: ['Strength', 'Chest'], type: [String], description: 'Categories associated with the session' })
    category: string[];

    @ApiProperty({
        description: 'JSON array of exercise targets',
        example: [
            {
                "exercise id": 1,
                "set1": { "targetRep": 10, "targetKg": 10 }
            }
        ]
    })
    sessionData: any;
}
