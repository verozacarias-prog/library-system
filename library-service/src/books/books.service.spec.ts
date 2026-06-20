import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { NotFoundException } from '@nestjs/common';
import { BooksService } from './books.service';
import { Book } from './book.entity';

const mockQueryBuilder = {
    update: jest.fn().mockReturnThis(),
    set: jest.fn().mockReturnThis(),
    where: jest.fn().mockReturnThis(),
    returning: jest.fn().mockReturnThis(),
    execute: jest.fn(),
};

const mockRepository = {
    create: jest.fn(),
    save: jest.fn(),
    findOne: jest.fn(),
    update: jest.fn(),
    delete: jest.fn(),
    createQueryBuilder: jest.fn(() => mockQueryBuilder),
};

describe('BooksService', () => {
    let service: BooksService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            providers: [
                BooksService,
                { provide: getRepositoryToken(Book), useValue: mockRepository },
            ],
        }).compile();

        service = module.get<BooksService>(BooksService);
        jest.clearAllMocks();
    });

    it('should create a book', async () => {
        const data = { title: 'Clean Code', author: 'Martin', isbn: '123', year: 2008, genre: 'tech', available_copies: 3 };
        mockRepository.create.mockReturnValue(data);
        mockRepository.save.mockResolvedValue({ id: 1, ...data });

        const result = await service.create(data);

        expect(result).toEqual({ id: 1, ...data });
        expect(mockRepository.save).toHaveBeenCalled();
    });

    it('should throw NotFoundException when book not found', async () => {
        mockRepository.findOne.mockResolvedValue(null);
        await expect(service.findOne(999)).rejects.toThrow(NotFoundException);
    });

    it('should update available copies', async () => {
        mockQueryBuilder.execute.mockResolvedValue({
            affected: 1,
            raw: [{ id: 1, title: 'Clean Code', available_copies: 4 }],
        });

        const result = await service.updateCopies(1, 1);

        expect(result).toEqual({ id: 1, title: 'Clean Code', available_copies: 4 });
        expect(mockQueryBuilder.where).toHaveBeenCalledWith(
            'id = :id AND available_copies + :delta >= 0',
            { id: 1, delta: 1 },
        );
    });
});