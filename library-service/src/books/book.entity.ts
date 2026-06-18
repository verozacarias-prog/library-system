import { Entity, PrimaryGeneratedColumn, Column } from 'typeorm';

@Entity('books')
export class Book {
    @PrimaryGeneratedColumn()
    id: number;

    @Column()
    title: string;

    @Column()
    author: string;

    @Column({ unique: true })
    isbn: string;

    @Column()
    year: number;

    @Column()
    genre: string;

    @Column({ name: 'available_copies', default: 1 })
    availableCopies: number;
}