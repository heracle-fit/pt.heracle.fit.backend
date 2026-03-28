import { ApiProperty } from '@nestjs/swagger';

export class AndroidTokenDto {
	@ApiProperty({ example: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Ij...', description: 'Google ID token from Android client' })
	idToken: string;

	@ApiProperty({ required: false, example: 'ya29.a0AfB_byC...', description: 'Google Access token from Android client' })
	accessToken?: string;
}
