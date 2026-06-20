import { Controller, Post, Body } from '@nestjs/common';
import { AuthService } from './auth.service';
import { LoginDto } from './dto/login.dto';
import { LoginResponseDto } from './dto/login-response.dto';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';

@ApiTags('auth')
@Controller('auth')
export class AuthController {
    constructor(private readonly authService: AuthService) { }

    @Post('login')
    @ApiOperation({ summary: 'Log in and receive a JWT' })
    @ApiResponse({ status: 201, description: 'Login successful — include the token as Authorization: Bearer <token> on protected routes', type: LoginResponseDto })
    @ApiResponse({ status: 401, description: 'Invalid email or password' })
    login(@Body() data: LoginDto) {
        return this.authService.login(data.email, data.password);
    }
}
