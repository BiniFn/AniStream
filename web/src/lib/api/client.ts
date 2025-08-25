import createClient from 'openapi-fetch';
import type { paths } from './openapi';
import { PUBLIC_API_URL } from '$env/static/public';

export const apiClient = createClient<paths>({
	baseUrl: PUBLIC_API_URL || 'http://localhost:8080',
	credentials: 'include',
});
