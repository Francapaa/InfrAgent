export interface UserProfile {
  id: string;
  email: string;
  name: string;
  company_name?: string;
  webhook_url?: string;
  metodo?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CompleteProfileRequest {
  company_name: string;
  webhook_url: string;
}

export interface CompleteProfileResponse {
  client_id: string;
  api_key: string;
  webhook_secret: string;
}

export interface ApiError {
  error: string;
}

export interface ApiResponse<T> {
  data: T | null;
  error: string | null;
}

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
