import { SetMetadata } from '@nestjs/common';

export const IS_TRAINER_KEY = 'isTrainer';
export const Trainer = () => SetMetadata(IS_TRAINER_KEY, true);
