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

    @ApiPropertyOptional({ example: 'weight_loss' })
    goal: string | null;

    @ApiProperty({ example: 0.75, description: "Today's progress toward nutritional targets (0.0 to 1.0)" })
    progress: number;
}

