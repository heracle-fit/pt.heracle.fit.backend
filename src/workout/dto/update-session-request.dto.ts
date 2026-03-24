import { ApiPropertyOptional } from '@nestjs/swagger';

export class UpdateSessionRequestDto {
    @ApiPropertyOptional({ example: 'Updated Session Name' })
    name?: string;

    @ApiPropertyOptional({ example: ['Strength'], type: [String] })
    category?: string[];

    @ApiPropertyOptional({ description: 'JSON array of exercise targets' })
    sessionData?: any;
}
