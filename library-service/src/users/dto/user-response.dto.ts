import { ApiProperty } from '@nestjs/swagger';

export class UserResponseDto {
    @ApiProperty({ example: 1 })
    id: number;

    @ApiProperty({ example: 'Jane Doe' })
    name: string;

    @ApiProperty({ example: 'jane@example.com' })
    email: string;

    @ApiProperty({ enum: ['user', 'admin'], example: 'user' })
    role: string;
}
