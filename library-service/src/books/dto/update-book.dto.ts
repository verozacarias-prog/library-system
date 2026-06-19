import { IsString, IsNotEmpty, IsInt, Min, IsISBN, IsOptional } from 'class-validator';

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
    @Min(0)
    year?: number;

    @IsOptional()
    @IsString()
    @IsNotEmpty()
    genre?: string;

    @IsOptional()
    @IsInt()
    @Min(0)
    availableCopies?: number;
}
