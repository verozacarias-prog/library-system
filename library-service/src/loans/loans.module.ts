import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { LoansController } from './loans.controller';

@Module({
    imports: [HttpModule.register({ timeout: 5000 })],
    controllers: [LoansController],
})
export class LoansModule { }
