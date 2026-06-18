import { Controller, Get, Post, Patch, Delete, Param, Body, Query, UseGuards } from '@nestjs/common';
import { BooksService } from './books.service';
import { Book } from './book.entity';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';


@Controller('books')
export class BooksController {
    constructor(private readonly booksService: BooksService) { }

    @Post()
    @UseGuards(JwtAuthGuard)
    create(@Body() data: Partial<Book>) {
        return this.booksService.create(data);
    }

    @Get()
    findAll(
        @Query('author') author?: string,
        @Query('genre') genre?: string,
        @Query('available') available?: string,
        @Query('page') page = '1',
        @Query('limit') limit = '10',
    ) {
        return this.booksService.findAll(
            { author, genre, available: available === 'true' },
            Number(page),
            Number(limit),
        );
    }

    @Get(':id')
    findOne(@Param('id') id: string) {
        return this.booksService.findOne(Number(id));
    }

    @Patch(':id')
    @UseGuards(JwtAuthGuard)
    update(@Param('id') id: string, @Body() data: Partial<Book>) {
        return this.booksService.update(Number(id), data);
    }

    @Delete(':id')
    @UseGuards(JwtAuthGuard)
    remove(@Param('id') id: string) {
        return this.booksService.remove(Number(id));
    }

    @Patch(':id/copies')
    @UseGuards(JwtAuthGuard)
    updateCopies(@Param('id') id: string, @Body('delta') delta: number) {
        return this.booksService.updateCopies(Number(id), delta);
    }
}
