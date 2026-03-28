import { Controller, Get, Post, Delete, Body, Param, Req } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiBearerAuth, ApiOkResponse, ApiBody } from '@nestjs/swagger';
import { TrainerService } from './trainer.service';
import { Trainer } from '../common/decorators/trainer.decorator';
import { Admin } from '../common/decorators/admin.decorator';
import { AddClientDto } from './dto/add-client.dto';
import { AddTrainerDto } from './dto/add-trainer.dto';
import { ClientResponseDto } from './dto/client-response.dto';


@ApiTags('Trainer')
@ApiBearerAuth('JWT')
@Controller('trainer')
export class TrainerController {
    constructor(private readonly trainerService: TrainerService) {}

    @Get('clients')
    @Trainer()
    @ApiOperation({ summary: 'Get all clients for the trainer' })
    @ApiOkResponse({ type: [ClientResponseDto] })
    async getClients(@Req() req: any): Promise<ClientResponseDto[]> {
        return this.trainerService.getClients(req.user.id);
    }

    @Post('clients/add')
    @Trainer()
    @ApiOperation({ summary: 'Add a client by email' })
    @ApiBody({ type: AddClientDto })
    @ApiOkResponse({ type: ClientResponseDto })
    async addClient(@Req() req: any, @Body() body: AddClientDto): Promise<ClientResponseDto> {
        return this.trainerService.addClient(req.user.id, body.email);
    }

    @Delete('clients/remove/:clientId')
    @Trainer()
    @ApiOperation({ summary: 'Remove a client' })
    @ApiOkResponse({ description: 'Client removed successfully' })
    async removeClient(@Req() req: any, @Param('clientId') clientId: string): Promise<void> {
        return this.trainerService.removeClient(req.user.id, clientId);
    }

    @Post('admin/add')
    @Admin()
    @ApiOperation({ summary: 'Add a new trainer (Admin Only)' })
    @ApiBody({ type: AddTrainerDto })
    @ApiOkResponse({ description: 'Trainer added successfully' })
    async adminAddTrainer(@Body() body: AddTrainerDto) {
        return this.trainerService.adminAddTrainer(body);
    }
}

