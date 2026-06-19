import { ConflictException, Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Book } from './book.entity';

const PG_UNIQUE_VIOLATION = '23505';

@Injectable()
export class BooksService {
    constructor(
        @InjectRepository(Book)
        private readonly bookRepository: Repository<Book>,
    ) { }

    async create(data: Partial<Book>): Promise<Book> {
        try {
            const book = this.bookRepository.create(data);
            return this.bookRepository.save(book);
        } catch (err: any) {
            if (err?.code === PG_UNIQUE_VIOLATION) {
                throw new ConflictException(`Book with ISBN ${data.isbn} already exists`);
            }
            throw err;
        }

    }

    async findAll(filters: { author?: string; genre?: string; available?: boolean }, page = 1, limit = 10): Promise<{ data: Book[]; total: number }> {
        const query = this.bookRepository.createQueryBuilder('book');

        if (filters.author) query.andWhere('book.author ILIKE :author', { author: `%${filters.author}%` });
        if (filters.genre) query.andWhere('book.genre = :genre', { genre: filters.genre });
        if (filters.available) query.andWhere('book.available_copies > 0');

        query.skip((page - 1) * limit).take(limit);

        const [data, total] = await query.getManyAndCount();
        return { data, total };
    }

    async findOne(id: number): Promise<Book> {
        const book = await this.bookRepository.findOne({ where: { id } });
        if (!book) throw new NotFoundException(`Book ${id} not found`);
        return book;
    }

    async update(id: number, data: Partial<Book>): Promise<Book> {
        await this.findOne(id);
        await this.bookRepository.update(id, data);
        return this.findOne(id);
    }

    async remove(id: number): Promise<void> {
        await this.findOne(id);
        await this.bookRepository.delete(id);
    }

    async updateCopies(id: number, delta: number): Promise<Book> {
        const book = await this.findOne(id);
        const next = book.available_copies + delta;
        if (next < 0) throw new ConflictException("not enough copies available");
        book.available_copies = next;
        return this.bookRepository.save(book);
    }
}
