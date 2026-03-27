import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class ClientResponseDto {
    @ApiProperty({ example: 'uuid-...' })
    id: string;

    @ApiProperty({ example: 'John Doe' })
    name: string;

    @ApiProperty({ example: 'john@example.com' })
    email: string;

    @ApiPropertyOptional({ example: 'https://example.com/avatar.jpg' })
    avatarUrl: string | null;

    @ApiProperty()
    assignedAt: Date;
}
