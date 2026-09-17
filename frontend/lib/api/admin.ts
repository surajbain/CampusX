import { apiFetch } from "./client";
import type { PaginationMeta } from "./types";
import type { User } from "./auth";
import type { College } from "./colleges";

// ---- College CRUD ----

export interface CreateCollegeInput {
  name: string;
  slug: string;
  city: string;
  state?: string;
  logo_url?: string;
  website?: string;
  contact_email?: string;
  contact_phone?: string;
}

export interface UpdateCollegeInput {
  name?: string;
  city?: string;
  state?: string;
  logo_url?: string;
  website?: string;
  contact_email?: string;
  contact_phone?: string;
  is_active?: boolean;
}

export async function createCollege(
  input: CreateCollegeInput
): Promise<College> {
  return apiFetch<College>("/api/v1/colleges", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function updateCollege(
  id: string,
  input: UpdateCollegeInput
): Promise<College> {
  return apiFetch<College>(`/api/v1/colleges/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
}

export async function deleteCollege(id: string): Promise<void> {
  await apiFetch(`/api/v1/colleges/${id}`, { method: "DELETE" });
}

export async function getCollege(id: string): Promise<College> {
  return apiFetch<College>(`/api/v1/colleges/${id}`);
}

// ---- User Management ----

export interface AdminUser {
  id: string;
  email: string;
  full_name: string;
  phone?: string;
  role: string;
  college_id?: string;
  is_active: boolean;
  email_verified: boolean;
  last_login_at?: string;
  created_at: string;
}

export interface CreateUserInput {
  email: string;
  password: string;
  full_name: string;
  phone?: string;
  role: "ORGANIZER" | "VOLUNTEER" | "COLLEGE_ADMIN";
}

export interface UpdateUserInput {
  full_name?: string;
  phone?: string;
  role?: string;
  is_active?: boolean;
}

interface ListUsersResponse {
  data: AdminUser[];
  meta: PaginationMeta;
}

export async function listCollegeUsers(
  collegeId: string,
  params: { page?: number; limit?: number } = {}
): Promise<ListUsersResponse> {
  return apiFetch<ListUsersResponse>(`/api/v1/colleges/${collegeId}/users`, {
    params: params as Record<string, string | number>,
  });
}

export async function createCollegeUser(
  collegeId: string,
  input: CreateUserInput
): Promise<AdminUser> {
  return apiFetch<AdminUser>(`/api/v1/colleges/${collegeId}/users`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function updateCollegeUser(
  collegeId: string,
  userId: string,
  input: UpdateUserInput
): Promise<AdminUser> {
  return apiFetch<AdminUser>(
    `/api/v1/colleges/${collegeId}/users/${userId}`,
    {
      method: "PATCH",
      body: JSON.stringify(input),
    }
  );
}

export async function deleteCollegeUser(
  collegeId: string,
  userId: string
): Promise<void> {
  await apiFetch(`/api/v1/colleges/${collegeId}/users/${userId}`, {
    method: "DELETE",
  });
}

// ---- Stats ----

export async function getAdminStats(): Promise<Record<string, number>> {
  return apiFetch<Record<string, number>>("/api/v1/colleges/stats");
}

export type { User };