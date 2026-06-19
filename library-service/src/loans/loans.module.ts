import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { LoansController } from './loans.controller';

@Module({
    imports: [HttpModule],
    controllers: [LoansController],
})
export class LoansModule {}
