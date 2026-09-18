import { apiFetch } from "./client";

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

// apiFetch returns {data, meta} wrapper — unwrap .data
function unwrap<T>(res: unknown): T {
  return (res as { data: T }).data;
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  const res = await apiFetch<AuthResponse>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return unwrap<AuthResponse>(res);
}

export async function register(input: RegisterInput): Promise<AuthResponse> {
  const res = await apiFetch<AuthResponse>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return unwrap<AuthResponse>(res);
}

export async function refresh(refreshToken: string): Promise<AuthResponse> {
  const res = await apiFetch<AuthResponse>("/api/v1/auth/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  return unwrap<AuthResponse>(res);
}

export async function logout(accessToken: string, refreshToken: string): Promise<void> {
  await apiFetch("/api/v1/auth/logout", {
    method: "POST",
    accessToken,
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
}

export async function me(accessToken: string): Promise<User> {
  const res = await apiFetch<User>("/api/v1/auth/me", { accessToken });
  return unwrap<User>(res);
}