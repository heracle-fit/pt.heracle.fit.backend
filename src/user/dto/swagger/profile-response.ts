import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class ProfileResponseDto {
    @ApiProperty({ example: 'a1b2c3d4-e5f6-7890-abcd-ef1234567890' })
    id: string;

    @ApiProperty({ example: 'sanjaysagar' })
    username: string;

    @ApiProperty({ example: 'SANJAY SAGAR' })
    name: string;

    @ApiProperty({ example: 'sanjaysagar.main@gmail.com' })
    email: string;

    @ApiPropertyOptional({ example: 'https://lh3.googleusercontent.com/...' })
    avatarUrl: string | null;

    @ApiPropertyOptional({ example: 'Athlete' })
    bio: string | null;

    @ApiProperty({ example: '2026-02-25T11:30:00.000Z' })
    createdAt: Date;

    @ApiPropertyOptional({ example: '2026-02-25T12:00:00.000Z' })
    updatedAt: Date | null;
}
