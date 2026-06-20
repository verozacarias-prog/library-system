import { IsEmail, IsString, IsNotEmpty } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class LoginDto {
    @ApiProperty({ example: 'admin@library.com' })
    @IsEmail()
    email: string;

    @ApiProperty({ example: 'adminpass' })
    @IsString()
    @IsNotEmpty()
    password: string;
}