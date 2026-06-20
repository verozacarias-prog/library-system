import { ApiProperty } from '@nestjs/swagger';

export class LoanResponseDto {
    @ApiProperty({ example: 1 })
    id: number;

    @ApiProperty({ example: 1 })
    user_id: number;

    @ApiProperty({ example: 1 })
    book_id: number;

    @ApiProperty({ example: '2024-01-15T10:30:00Z' })
    loaned_at: string;

    @ApiProperty({ example: null, nullable: true })
    returned_at: string | null;

    @ApiProperty({ enum: ['active', 'returned'], example: 'active' })
    status: string;
}
