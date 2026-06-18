import { Test, TestingModule } from '@nestjs/testing';
import { JwtService } from '@nestjs/jwt';
import { UnauthorizedException } from '@nestjs/common';
import { AuthService } from './auth.service';
import { UsersService } from '../users/users.service';
import * as bcrypt from 'bcrypt';

const mockUsersService = {
    findByEmail: jest.fn(),
};

const mockJwtService = {
    sign: jest.fn().mockReturnValue('mock_token'),
};

describe('AuthService', () => {
    let service: AuthService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            providers: [
                AuthService,
                { provide: UsersService, useValue: mockUsersService },
                { provide: JwtService, useValue: mockJwtService },
            ],
        }).compile();

        service = module.get<AuthService>(AuthService);
        jest.clearAllMocks();
    });

    it('should return access_token on valid credentials', async () => {
        const hashed = await bcrypt.hash('password123', 10);
        mockUsersService.findByEmail.mockResolvedValue({
            id: 1, email: 'vero@test.com', password: hashed, role: 'admin',
        });

        const result = await service.login('vero@test.com', 'password123');

        expect(result).toHaveProperty('access_token');
        expect(mockJwtService.sign).toHaveBeenCalled();
    });

    it('should throw UnauthorizedException when user not found', async () => {
        mockUsersService.findByEmail.mockResolvedValue(null);
        await expect(service.login('noexiste@test.com', 'pass')).rejects.toThrow(UnauthorizedException);
    });

    it('should throw UnauthorizedException when password is wrong', async () => {
        const hashed = await bcrypt.hash('correctpass', 10);
        mockUsersService.findByEmail.mockResolvedValue({
            id: 1, email: 'vero@test.com', password: hashed, role: 'user',
        });
        await expect(service.login('vero@test.com', 'wrongpass')).rejects.toThrow(UnauthorizedException);
    });
});