import { ApiError, ApiResponse, HttpMethod } from './types';

const BASE_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8080';

function getAuthToken(): string | null {
  if (typeof document === 'undefined') return null;
  
  const cookies = document.cookie.split(';');
  const authCookie = cookies.find(cookie => cookie.trim().startsWith('auth_token='));
  
  if (!authCookie) return null;
  
  return authCookie.split('=')[1]?.trim() || null;
}

interface FetchOptions extends Omit<RequestInit, 'method'> {
  method?: HttpMethod;
}

export async function fetchWithAuth<T>(
  endpoint: string,
  options: FetchOptions = {}
): Promise<ApiResponse<T>> {
  const token = getAuthToken();
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const fetchOptions: RequestInit = {
    ...options,
    headers,
    method: options.method || 'GET',
  };

  try {
    const response = await fetch(`${BASE_URL}${endpoint}`, fetchOptions);
    
    if (!response.ok) {
      const errorData: ApiError = await response.json().catch(() => ({ error: 'Unknown error' }));
      return {
        data: null,
        error: errorData.error || `HTTP Error: ${response.status}`,
      };
    }

    if (response.status === 204) {
      return { data: null as T, error: null };
    }

    const data: T = await response.json();
    return { data, error: null };
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Network error';
    return {
      data: null,
      error: errorMessage,
    };
  }
}

export function get<T>(endpoint: string, options?: Omit<FetchOptions, 'method'>): Promise<ApiResponse<T>> {
  return fetchWithAuth<T>(endpoint, { ...options, method: 'GET' });
}

export function post<T>(
  endpoint: string,
  body: Record<string, unknown>,
  options?: Omit<FetchOptions, 'method' | 'body'>
): Promise<ApiResponse<T>> {
  return fetchWithAuth<T>(endpoint, { ...options, method: 'POST', body: JSON.stringify(body) });
}
