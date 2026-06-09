import { apiRequest } from "@/shared/services/api-client";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  ChangePasswordRequest,
  UserInfo,
} from "@/features/auth/types/auth.types";

export const authService = {
  register(data: RegisterRequest): Promise<AuthResponse> {
    return apiRequest<AuthResponse>("/auth/register", {
      method: "POST",
      body: data,
      skipAuth: true,
    });
  },

  login(data: LoginRequest): Promise<AuthResponse> {
    return apiRequest<AuthResponse>("/auth/login", {
      method: "POST",
      body: data,
      skipAuth: true,
    });
  },

  me(): Promise<UserInfo> {
    return apiRequest<UserInfo>("/auth/me");
  },

  logout(refreshToken: string): Promise<void> {
    return apiRequest<void>("/auth/logout", {
      method: "POST",
      body: { refreshToken },
      skipAuth: true,
    });
  },

  changePassword(data: ChangePasswordRequest): Promise<void> {
    return apiRequest<void>("/auth/change-password", {
      method: "POST",
      body: data,
    });
  },

  refresh(refreshToken: string): Promise<AuthResponse> {
    return apiRequest<AuthResponse>("/auth/refresh", {
      method: "POST",
      body: { refreshToken },
      skipAuth: true,
    });
  },
};
