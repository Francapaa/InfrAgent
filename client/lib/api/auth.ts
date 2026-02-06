import { get, post } from './client';
import {
  UserProfile,
  CompleteProfileRequest,
  CompleteProfileResponse,
  ApiResponse,
} from './types';

export async function getCurrentUser(): Promise<ApiResponse<UserProfile>> {
  return get<UserProfile>('/auth/me');
}

export async function completeProfile(
  data: CompleteProfileRequest
): Promise<ApiResponse<CompleteProfileResponse>> {
  return post<CompleteProfileResponse>('/auth/complete-registration', {
    company_name: data.company_name,
    webhook_url: data.webhook_url,
  });
}
