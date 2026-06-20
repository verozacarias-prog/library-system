import { Entity, PrimaryGeneratedColumn, Column } from 'typeorm';
import { ApiProperty } from '@nestjs/swagger';

@Entity('books')
export class Book {
    @ApiProperty({ example: 1 })
    @PrimaryGeneratedColumn()
    id: number;

    @ApiProperty({ example: 'Clean Code' })
    @Column()
    title: string;

    @ApiProperty({ example: 'Robert C. Martin' })
    @Column()
    author: string;

    @ApiProperty({ example: '9780132350884' })
    @Column({ unique: true })
    isbn: string;

    @ApiProperty({ example: 2008 })
    @Column()
    year: number;

    @ApiProperty({ example: 'tech' })
    @Column()
    genre: string;

    @ApiProperty({ example: 3 })
    @Column({ name: 'available_copies' })
    available_copies: number;
}