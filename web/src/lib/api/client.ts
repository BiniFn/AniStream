import createClient from 'openapi-fetch';
import { env } from '$env/dynamic/public';
import type { paths } from './openapi';

export const apiClient = createClient<paths>({
	baseUrl: env.PUBLIC_API_URL || 'http://localhost:8080',
	credentials: 'include',
});
