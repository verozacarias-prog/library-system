import { IsEmail, IsString, IsNotEmpty, MinLength, IsOptional, IsIn } from 'class-validator';

export class UpdateUserDto {
    @IsOptional()
    @IsString()
    @IsNotEmpty()
    name?: string;

    @IsOptional()
    @IsEmail()
    email?: string;

    @IsOptional()
    @IsString()
    @MinLength(8)
    password?: string;

    @IsOptional()
    @IsIn(['user', 'admin'])
    role?: string;
}
