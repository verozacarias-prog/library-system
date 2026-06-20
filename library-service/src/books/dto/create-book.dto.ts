import { IsString, IsNotEmpty, IsInt, Min, IsISBN } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';
import { IsNotFutureYear } from './validators';

export class CreateBookDto {

    @ApiProperty({ example: 'Clean Code' })
    @IsString()
    @IsNotEmpty()
    title: string;

    @ApiProperty({ example: 'Robert C. Martin' })
    @IsString()
    @IsNotEmpty()
    author: string;

    @ApiProperty({ example: '9780132350884' })
    @IsISBN()
    isbn: string;

    @ApiProperty({ example: 2008 })
    @IsInt()
    @Min(1450)
    @IsNotFutureYear()
    year: number;

    @ApiProperty({ example: 'tech' })
    @IsString()
    @IsNotEmpty()
    genre: string;

    @ApiProperty({ example: 3 })
    @IsInt()
    @Min(1)
    available_copies: number;
}
