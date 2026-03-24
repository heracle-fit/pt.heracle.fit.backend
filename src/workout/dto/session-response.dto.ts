import { ApiProperty } from '@nestjs/swagger';

export class SessionResponseDto {
    @ApiProperty({ example: 1 })
    id: number;

    @ApiProperty({ example: 'Morning Push Session' })
    name: string;

    @ApiProperty({ example: ['Strength', 'Chest'], type: [String] })
    category: string[];

    @ApiProperty()
    sessionData: any;

    @ApiProperty()
    createdAt: Date;

    @ApiProperty()
    updatedAt: Date;
}
