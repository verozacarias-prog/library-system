import { Controller, Get, Post, Patch, Delete, Param, Body, Query, Request, UseGuards, BadRequestException, ForbiddenException, HttpCode, HttpStatus } from '@nestjs/common';
import { BooksService } from './books.service';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { RolesGuard } from '../auth/roles.guard';
import { Roles } from '../auth/roles.decorator';
import { CreateBookDto } from './dto/create-book.dto';
import { UpdateBookDto } from './dto/update-book.dto';


@Controller('books')
export class BooksController {
    constructor(private readonly booksService: BooksService) { }

    @Post()
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    create(@Body() data: CreateBookDto) {
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
        const parsedPage = Math.max(1, parseInt(page, 10) || 1);
        const parsedLimit = Math.min(100, Math.max(1, parseInt(limit, 10) || 10));
        return this.booksService.findAll(
            { author, genre, available: available === 'true' },
            parsedPage,
            parsedLimit,
        );
    }

    @Get(':id')
    findOne(@Param('id') id: string) {
        return this.booksService.findOne(Number(id));
    }

    @Patch(':id')
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    update(@Param('id') id: string, @Body() data: UpdateBookDto) {
        return this.booksService.update(Number(id), data);
    }

    @Delete(':id')
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    @HttpCode(HttpStatus.NO_CONTENT)
    remove(@Param('id') id: string) {
        return this.booksService.remove(Number(id));
    }

    @Patch(':id/copies')
    @UseGuards(JwtAuthGuard)
    updateCopies(@Param('id') id: string, @Body('delta') delta: number, @Request() req) {
        if (req.user.role !== 'service') {
            throw new ForbiddenException('This endpoint is reserved for internal service use');
        }
        const d = Number(delta);
        if (!Number.isInteger(d) || d === 0) {
            throw new BadRequestException('delta must be a non-zero integer');
        }
        return this.booksService.updateCopies(Number(id), d);
    }
}
