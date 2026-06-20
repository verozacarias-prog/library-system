import { Controller, Get, Post, Patch, Delete, Param, Body, Query, Request, UseGuards, BadRequestException, ForbiddenException, HttpCode, HttpStatus } from '@nestjs/common';
import { BooksService } from './books.service';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { RolesGuard } from '../auth/roles.guard';
import { Roles } from '../auth/roles.decorator';
import { CreateBookDto } from './dto/create-book.dto';
import { UpdateBookDto } from './dto/update-book.dto';
import { Book } from './book.entity';
import { ApiTags, ApiBearerAuth, ApiQuery, ApiOperation, ApiResponse, ApiParam, ApiExcludeEndpoint } from '@nestjs/swagger';

@ApiTags('books')
@ApiBearerAuth()
@Controller('books')
export class BooksController {
    constructor(private readonly booksService: BooksService) { }

    @Post()
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    @ApiOperation({ summary: 'Create a book (admin only)' })
    @ApiResponse({ status: 201, description: 'Book created', type: Book })
    @ApiResponse({ status: 400, description: 'Validation error' })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 403, description: 'Requires admin role' })
    create(@Body() data: CreateBookDto) {
        return this.booksService.create(data);
    }

    @ApiOperation({ summary: 'List books with optional filters and pagination' })
    @ApiQuery({ name: 'author', required: false, description: 'Filter by author name' })
    @ApiQuery({ name: 'genre', required: false, description: 'Filter by genre' })
    @ApiQuery({ name: 'available', required: false, type: Boolean, description: 'Filter by availability' })
    @ApiQuery({ name: 'page', required: false, type: Number, description: 'Page number (default: 1)' })
    @ApiQuery({ name: 'limit', required: false, type: Number, description: 'Results per page, max 100 (default: 10)' })
    @ApiResponse({ status: 200, description: 'Paginated list of books', schema: { example: { data: [], total: 0 } } })
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
    @ApiOperation({ summary: 'Get a book by ID' })
    @ApiParam({ name: 'id', description: 'Book ID', example: 1 })
    @ApiResponse({ status: 200, description: 'Book found', type: Book })
    @ApiResponse({ status: 404, description: 'Book not found' })
    findOne(@Param('id') id: string) {
        return this.booksService.findOne(Number(id));
    }

    @Patch(':id')
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    @ApiOperation({ summary: 'Update a book (admin only)' })
    @ApiParam({ name: 'id', description: 'Book ID', example: 1 })
    @ApiResponse({ status: 200, description: 'Book updated', type: Book })
    @ApiResponse({ status: 400, description: 'Validation error' })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 403, description: 'Requires admin role' })
    @ApiResponse({ status: 404, description: 'Book not found' })
    update(@Param('id') id: string, @Body() data: UpdateBookDto) {
        return this.booksService.update(Number(id), data);
    }

    @Delete(':id')
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    @HttpCode(HttpStatus.NO_CONTENT)
    @ApiOperation({ summary: 'Delete a book (admin only)' })
    @ApiParam({ name: 'id', description: 'Book ID', example: 1 })
    @ApiResponse({ status: 204, description: 'Book deleted' })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 403, description: 'Requires admin role' })
    @ApiResponse({ status: 404, description: 'Book not found' })
    remove(@Param('id') id: string) {
        return this.booksService.remove(Number(id));
    }

    @Patch(':id/copies')
    @UseGuards(JwtAuthGuard)
    @ApiExcludeEndpoint()
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
