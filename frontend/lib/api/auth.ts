import { apiFetch, ApiClientError } from "./client";

export interface User {
  id: string;
  email: string;
  full_name: string;
  phone?: string;
  role: "STUDENT" | "COLLEGE_ADMIN" | "ORGANIZER" | "VOLUNTEER" | "SUPER_ADMIN" | "ADVERTISER";
  college_id?: string;
  email_verified: boolean;
  created_at: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  email: string;
  password: string;
  full_name: string;
  phone?: string;
  college_id: string;
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  return apiFetch<AuthResponse>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function register(input: RegisterInput): Promise<AuthResponse> {
  return apiFetch<AuthResponse>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function refresh(refreshToken: string): Promise<AuthResponse> {
  return apiFetch<AuthResponse>("/api/v1/auth/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
}

export async function logout(accessToken: string, refreshToken: string): Promise<void> {
  await apiFetch("/api/v1/auth/logout", {
    method: "POST",
    accessToken,
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
}

export async function me(accessToken: string): Promise<User> {
  return apiFetch<User>("/api/v1/auth/me", { accessToken });
}

export { ApiClientError };