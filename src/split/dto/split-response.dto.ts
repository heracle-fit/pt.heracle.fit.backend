import { ApiProperty } from '@nestjs/swagger';

export class SplitResponseDto {
    @ApiProperty()
    id: string;

    @ApiProperty()
    trainerId: string;

    @ApiProperty()
    userId: string;

    @ApiProperty()
    splitData: any;

    @ApiProperty()
    createdAt: Date;

    @ApiProperty()
    updatedAt: Date;
}
