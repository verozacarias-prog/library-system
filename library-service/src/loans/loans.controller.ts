import {
    Controller, Post, Patch, Get, Param, Body, UseGuards,
    HttpException, HttpStatus,
} from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { CreateLoanDto } from './dto/create-loan.dto';
import { LoanResponseDto } from './dto/loan-response.dto';
import { ApiTags, ApiBearerAuth, ApiOperation, ApiResponse, ApiParam } from '@nestjs/swagger';

@ApiTags('loans')
@ApiBearerAuth()
@Controller('loans')
@UseGuards(JwtAuthGuard)
export class LoansController {
    constructor(private readonly http: HttpService) { }

    private get baseUrl(): string {
        return process.env.LOANS_SERVICE_URL ?? 'http://localhost:8081';
    }

    @Post()
    @ApiOperation({ summary: 'Create a loan' })
    @ApiResponse({ status: 201, description: 'Loan created', type: LoanResponseDto })
    @ApiResponse({ status: 400, description: 'Validation error' })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 409, description: 'No copies available or user already has this book on loan' })
    async create(@Body() body: CreateLoanDto) {
        try {
            const { data } = await firstValueFrom(
                this.http.post(`${this.baseUrl}/loans`, body),
            );
            return data;
        } catch (err: any) {
            const status = err?.response?.status ?? HttpStatus.INTERNAL_SERVER_ERROR;
            const message = err?.response?.data ?? 'loans service error';
            throw new HttpException(message, status);
        }
    }

    @Patch(':id')
    @ApiOperation({ summary: 'Return a book (marks loan as returned)' })
    @ApiParam({ name: 'id', description: 'Loan ID', example: 1 })
    @ApiResponse({ status: 200, description: 'Loan updated to returned', type: LoanResponseDto })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 404, description: 'Loan not found' })
    @ApiResponse({ status: 409, description: 'Loan already returned' })
    async returnLoan(@Param('id') id: string) {
        try {
            const { data } = await firstValueFrom(
                this.http.patch(`${this.baseUrl}/loans/${id}`),
            );
            return data;
        } catch (err: any) {
            const status = err?.response?.status ?? HttpStatus.INTERNAL_SERVER_ERROR;
            const message = err?.response?.data ?? 'loans service error';
            throw new HttpException(message, status);
        }
    }

    @Get('users/:userId')
    @ApiOperation({ summary: 'Get active loans for a user' })
    @ApiParam({ name: 'userId', description: 'User ID', example: 1 })
    @ApiResponse({ status: 200, description: 'Active loans', type: [LoanResponseDto] })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    async getActiveLoans(@Param('userId') userId: string) {
        try {
            const { data } = await firstValueFrom(
                this.http.get(`${this.baseUrl}/loans/users/${userId}`),
            );
            return data;
        } catch (err: any) {
            const status = err?.response?.status ?? HttpStatus.INTERNAL_SERVER_ERROR;
            const message = err?.response?.data ?? 'loans service error';
            throw new HttpException(message, status);
        }
    }

    @Get('users/:userId/history')
    @ApiOperation({ summary: 'Get full loan history for a user' })
    @ApiParam({ name: 'userId', description: 'User ID', example: 1 })
    @ApiResponse({ status: 200, description: 'All loans (active and returned)', type: [LoanResponseDto] })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    async getLoanHistory(@Param('userId') userId: string) {
        try {
            const { data } = await firstValueFrom(
                this.http.get(`${this.baseUrl}/loans/users/${userId}/history`),
            );
            return data;
        } catch (err: any) {
            const status = err?.response?.status ?? HttpStatus.INTERNAL_SERVER_ERROR;
            const message = err?.response?.data ?? 'loans service error';
            throw new HttpException(message, status);
        }
    }
}
