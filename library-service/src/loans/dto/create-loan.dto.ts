import { IsInt, IsPositive } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class CreateLoanDto {
    @ApiProperty({ example: 1 })
    @IsInt()
    @IsPositive()
    book_id: number;
}
