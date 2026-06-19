import { IsString, IsNotEmpty, IsInt, Min, IsISBN, IsOptional } from 'class-validator';

export class CreateBookDto {
    @IsString()
    @IsNotEmpty()
    title: string;

    @IsString()
    @IsNotEmpty()
    author: string;

    @IsISBN()
    isbn: string;

    @IsInt()
    @Min(0)
    year: number;

    @IsString()
    @IsNotEmpty()
    genre: string;

    @IsInt()
    @Min(0)
    available_copies: number;
}
