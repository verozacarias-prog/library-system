import { IsString, IsNotEmpty, IsInt, Min, IsISBN } from 'class-validator';
import { IsNotFutureYear } from './validators';

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
    @Min(1450)
    @IsNotFutureYear()
    year: number;

    @IsString()
    @IsNotEmpty()
    genre: string;

    @IsInt()
    @Min(1)
    available_copies: number;
}
