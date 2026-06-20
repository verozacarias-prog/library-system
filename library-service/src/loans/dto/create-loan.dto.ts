import { IsInt, IsPositive } from 'class-validator';

export class CreateLoanDto {
    @IsInt()
    @IsPositive()
    user_id: number;

    @IsInt()
    @IsPositive()
    book_id: number;
}
