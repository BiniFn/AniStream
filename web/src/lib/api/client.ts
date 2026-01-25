import createClient from 'openapi-fetch';
import { PUBLIC_API_URL } from '$env/static/public';
import type { paths } from './openapi';

export const apiClient = createClient<paths>({
	baseUrl: PUBLIC_API_URL || 'http://localhost:8080',
	credentials: 'include',
});
