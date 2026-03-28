import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IsEmail, IsNotEmpty, IsOptional, IsString, IsInt, Min } from 'class-validator';

export class AddTrainerDto {
    @ApiProperty({ example: 'trainer@example.com' })
    @IsEmail()
    @IsNotEmpty()
    email: string;

    @ApiPropertyOptional({ example: 'Strength & Conditioning' })
    @IsOptional()
    @IsString()
    specialization?: string;

    @ApiPropertyOptional({ example: 5 })
    @IsOptional()
    @IsInt()
    @Min(0)
    experience?: number;
}
