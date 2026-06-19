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

    @IsOptional()
    @IsInt()
    @Min(0)
    availableCopies?: number;
}
