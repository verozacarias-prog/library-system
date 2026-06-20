import { Controller, Get, Post, Patch, Delete, Param, Body, Request, ForbiddenException, HttpCode, HttpStatus } from '@nestjs/common';
import { ApiTags, ApiBearerAuth, ApiOperation, ApiResponse, ApiParam } from '@nestjs/swagger';
import { UsersService } from './users.service';
import { UseGuards } from '@nestjs/common';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { RolesGuard } from '../auth/roles.guard';
import { Roles } from '../auth/roles.decorator';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { UserResponseDto } from './dto/user-response.dto';

@ApiTags('users')
@ApiBearerAuth()
@Controller('users')
export class UsersController {
    constructor(private readonly usersService: UsersService) { }

    @Post()
    @ApiOperation({ summary: 'Register a new user' })
    @ApiResponse({ status: 201, description: 'User created', type: UserResponseDto })
    @ApiResponse({ status: 400, description: 'Validation error' })
    @ApiResponse({ status: 409, description: 'Email already in use' })
    create(@Body() data: CreateUserDto) {
        return this.usersService.create(data);
    }

    @Get()
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    @ApiOperation({ summary: 'List all users (admin only)' })
    @ApiResponse({ status: 200, description: 'List of users', type: [UserResponseDto] })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 403, description: 'Requires admin role' })
    findAll() {
        return this.usersService.findAll();
    }

    @Get(':id')
    @UseGuards(JwtAuthGuard)
    @ApiOperation({ summary: 'Get a user by ID' })
    @ApiParam({ name: 'id', description: 'User ID', example: 1 })
    @ApiResponse({ status: 200, description: 'User found', type: UserResponseDto })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 404, description: 'User not found' })
    findOne(@Param('id') id: string) {
        return this.usersService.findOne(Number(id));
    }

    @Patch(':id')
    @UseGuards(JwtAuthGuard)
    @ApiOperation({ summary: 'Update a user (role field requires admin)' })
    @ApiParam({ name: 'id', description: 'User ID', example: 1 })
    @ApiResponse({ status: 200, description: 'User updated', type: UserResponseDto })
    @ApiResponse({ status: 400, description: 'Validation error' })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 403, description: 'Only admins can change the role field' })
    @ApiResponse({ status: 404, description: 'User not found' })
    update(@Param('id') id: string, @Body() data: UpdateUserDto, @Request() req) {
        if (data.role !== undefined && req.user.role !== 'admin') {
            throw new ForbiddenException('Only admins can change roles');
        }
        return this.usersService.update(Number(id), data);
    }

    @Delete(':id')
    @UseGuards(JwtAuthGuard, RolesGuard)
    @Roles('admin')
    @HttpCode(HttpStatus.NO_CONTENT)
    @ApiOperation({ summary: 'Delete a user (admin only)' })
    @ApiParam({ name: 'id', description: 'User ID', example: 1 })
    @ApiResponse({ status: 204, description: 'User deleted' })
    @ApiResponse({ status: 401, description: 'Missing or invalid token' })
    @ApiResponse({ status: 403, description: 'Requires admin role' })
    @ApiResponse({ status: 404, description: 'User not found' })
    remove(@Param('id') id: string) {
        return this.usersService.remove(Number(id));
    }
}
