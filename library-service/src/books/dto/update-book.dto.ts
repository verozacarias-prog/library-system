import { IsString, IsNotEmpty, IsInt, Min, IsISBN, IsOptional } from 'class-validator';
import { IsNotFutureYear } from './validators';

export class UpdateBookDto {
    @IsOptional()
    @IsString()
    @IsNotEmpty()
    title?: string;

    @IsOptional()
    @IsString()
    @IsNotEmpty()
    author?: string;

    @IsOptional()
    @IsISBN()
    isbn?: string;

    @IsOptional()
    @IsInt()
    @Min(1450)
    @IsNotFutureYear()
    year?: number;

    @IsOptional()
    @IsString()
    @IsNotEmpty()
    genre?: string;

    @IsOptional()
    @IsInt()
    @Min(0)
    available_copies?: number;
}
