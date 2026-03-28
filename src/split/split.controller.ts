import { Controller, Get, Put, Body, Param, Req, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiBearerAuth, ApiOkResponse } from '@nestjs/swagger';
import { SplitService } from './split.service';
import { UpdateSplitDto } from './dto/update-split.dto';
import { SplitResponseDto } from './dto/split-response.dto';
import { Trainer } from '../common/decorators/trainer.decorator';

@ApiTags('Workout Split')
@ApiBearerAuth('JWT')
@Controller('split')
export class SplitController {
    constructor(private readonly splitService: SplitService) {}

    @Get()
    @ApiOperation({ summary: 'Get current user workout split' })
    @ApiOkResponse({ type: SplitResponseDto })
    async getMySplit(@Req() req: any): Promise<SplitResponseDto> {
        return this.splitService.getSplit(req.user.id);
    }

    @Get('trainer/:clientId')
    @Trainer()
    @ApiOperation({ summary: 'Get a client workout split (Trainer Only)' })
    @ApiOkResponse({ type: SplitResponseDto })
    async trainerGetClientSplit(
        @Req() req: any,
        @Param('clientId') clientId: string,
    ): Promise<SplitResponseDto> {
        return this.splitService.trainerGetSplit(req.user.id, clientId);
    }

    @Put('trainer/:clientId')
    @Trainer()
    @ApiOperation({ summary: 'Add or update a client workout split (Trainer Only)' })
    @ApiOkResponse({ type: SplitResponseDto })
    async upsertClientSplit(
        @Req() req: any,
        @Param('clientId') clientId: string,
        @Body() dto: UpdateSplitDto,
    ): Promise<SplitResponseDto> {
        return this.splitService.upsertSplit(req.user.id, clientId, dto);
    }
}
