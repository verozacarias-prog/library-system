import {
    Controller, Post, Patch, Get, Param, Body, UseGuards,
    HttpException, HttpStatus,
} from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';

@Controller('loans')
@UseGuards(JwtAuthGuard)
export class LoansController {
    constructor(private readonly http: HttpService) {}

    private get baseUrl(): string {
        return process.env.LOANS_SERVICE_URL ?? 'http://localhost:8081';
    }

    @Post()
    async create(@Body() body: { user_id: number; book_id: number }) {
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
