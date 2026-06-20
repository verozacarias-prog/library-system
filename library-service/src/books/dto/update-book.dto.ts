import { IsString, IsNotEmpty, IsInt, Min, IsISBN, IsOptional } from 'class-validator';
import { ApiPropertyOptional } from '@nestjs/swagger';
import { IsNotFutureYear } from './validators';

export class UpdateBookDto {
    @ApiPropertyOptional({ example: 'Clean Code' })
    @IsOptional()
    @IsString()
    @IsNotEmpty()
    title?: string;

    @ApiPropertyOptional({ example: 'Robert C. Martin' })
    @IsOptional()
    @IsString()
    @IsNotEmpty()
    author?: string;

    @ApiPropertyOptional({ example: '9780132350884' })
    @IsOptional()
    @IsISBN()
    isbn?: string;

    @ApiPropertyOptional({ example: 2008, minimum: 1450 })
    @IsOptional()
    @IsInt()
    @Min(1450)
    @IsNotFutureYear()
    year?: number;

    @ApiPropertyOptional({ example: 'tech' })
    @IsOptional()
    @IsString()
    @IsNotEmpty()
    genre?: string;

    @ApiPropertyOptional({ example: 3, minimum: 0 })
    @IsOptional()
    @IsInt()
    @Min(0)
    available_copies?: number;
}
