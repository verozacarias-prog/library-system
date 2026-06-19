import { Injectable, NotFoundException, ConflictException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { User } from './user.entity';

@Injectable()
export class UsersService {
    constructor(
        @InjectRepository(User)
        private readonly userRepository: Repository<User>,
    ) { }

    async create(data: { name: string; email: string; password: string; role?: string }): Promise<Omit<User, 'password'>> {
        const existing = await this.userRepository.findOne({ where: { email: data.email } });
        if (existing) throw new ConflictException('Email already registered');

        const hashed = await bcrypt.hash(data.password, 10);
        const user = this.userRepository.create({ ...data, password: hashed });
        const saved = await this.userRepository.save(user);
        const { password, ...result } = saved;
        return result;
    }

    async findAll(): Promise<Omit<User, 'password'>[]> {
        const users = await this.userRepository.find();
        return users.map(({ password, ...u }) => u);
    }

    private async findOneWithPassword(id: number): Promise<User> {
        const user = await this.userRepository.findOne({ where: { id } });
        if (!user) throw new NotFoundException(`User ${id} not found`);
        return user;
    }

    async findOne(id: number): Promise<Omit<User, 'password'>> {
        const { password, ...result } = await this.findOneWithPassword(id);
        return result;
    }

    async findByEmail(email: string): Promise<User | null> {
        return this.userRepository.findOne({ where: { email } });
    }

    async update(id: number, data: Partial<User>): Promise<Omit<User, 'password'>> {
        await this.findOneWithPassword(id);
        if (data.password) data.password = await bcrypt.hash(data.password, 10);
        await this.userRepository.update(id, data);
        const { password, ...result } = await this.findOneWithPassword(id);
        return result;
    }

    async remove(id: number): Promise<void> {
        await this.findOneWithPassword(id);
        await this.userRepository.delete(id);
    }
}