import { Controller, Get, Post, Patch, Delete, Param, Body } from '@nestjs/common';
import { UsersService } from './users.service';

@Controller('users')
export class UsersController {
    constructor(private readonly usersService: UsersService) { }

    @Post()
    create(@Body() data: { name: string; email: string; password: string; role?: string }) {
        return this.usersService.create(data);
    }

    @Get()
    findAll() {
        return this.usersService.findAll();
    }

    @Get(':id')
    findOne(@Param('id') id: string) {
        return this.usersService.findOne(Number(id));
    }

    @Patch(':id')
    update(@Param('id') id: string, @Body() data: any) {
        return this.usersService.update(Number(id), data);
    }

    @Delete(':id')
    remove(@Param('id') id: string) {
        return this.usersService.remove(Number(id));
    }
}
