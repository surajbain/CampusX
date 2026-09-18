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

// ============================================
// EVENT MANAGEMENT (Admin)
// ============================================

export interface CreateEventInput {
  title: string;
  slug?: string;
  description?: string;
  category: "HACKATHON" | "CULTURAL" | "SPORTS" | "WORKSHOP" | "TECH_FEST" | "OTHER";
  poster_url?: string;
  venue?: string;
  city?: string;
  starts_at?: string;
  ends_at?: string;
  registration_opens?: string;
  registration_closes?: string;
  price_paise?: number;
  capacity?: number;
  allow_teams?: boolean;
  team_size_min?: number;
  team_size_max?: number;
  prize_pool_paise?: number;
  contact_email?: string;
  whatsapp_link?: string;
  rules?: string;
  is_featured?: boolean;
}

export interface AdminEvent {
  id: string;
  college_id: string;
  created_by?: string;
  slug: string;
  title: string;
  description?: string;
  category: string;
  poster_url?: string;
  venue?: string;
  city?: string;
  status: "DRAFT" | "PUBLISHED" | "ONGOING" | "COMPLETED" | "CANCELLED";
  starts_at?: string;
  ends_at?: string;
  registration_opens?: string;
  registration_closes?: string;
  price_paise: number;
  currency: string;
  capacity?: number;
  allow_teams: boolean;
  team_size_min?: number;
  team_size_max?: number;
  prize_pool_paise: number;
  contact_email?: string;
  whatsapp_link?: string;
  rules?: string;
  is_featured: boolean;
  college_name?: string;
  college_slug?: string;
  college_city?: string;
  created_at: string;
  updated_at: string;
}

interface ListEventsAdminResponse {
  data: AdminEvent[];
  meta: PaginationMeta;
}

// List events the current admin can manage
export async function listAdminEvents(params: {
  page?: number;
  limit?: number;
  status?: string;
} = {}): Promise<ListEventsAdminResponse> {
  return apiFetch<ListEventsAdminResponse>("/api/v1/events/my", {
    params: params as Record<string, string | number>,
  });
}

// Create event
export async function createEvent(input: CreateEventInput): Promise<AdminEvent> {
  const res = await apiFetch<AdminEvent>("/api/v1/events", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return (res as unknown as { data: AdminEvent }).data;
}

// Update event
export async function updateEvent(
  id: string,
  input: Partial<CreateEventInput>
): Promise<AdminEvent> {
  const res = await apiFetch<AdminEvent>(`/api/v1/events/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return (res as unknown as { data: AdminEvent }).data;
}

// Delete event
export async function deleteEvent(id: string): Promise<void> {
  await apiFetch(`/api/v1/events/${id}`, { method: "DELETE" });
}

// Publish event
export async function publishEvent(id: string): Promise<AdminEvent> {
  const res = await apiFetch<AdminEvent>(`/api/v1/events/${id}/publish`, {
    method: "POST",
  });
  return (res as unknown as { data: AdminEvent }).data;
}

// Cancel event
export async function cancelEvent(id: string): Promise<AdminEvent> {
  const res = await apiFetch<AdminEvent>(`/api/v1/events/${id}/cancel`, {
    method: "POST",
  });
  return (res as unknown as { data: AdminEvent }).data;
}

// Get single event
export async function getAdminEvent(id: string): Promise<AdminEvent> {
  const res = await apiFetch<AdminEvent>(`/api/v1/events/${id}`);
  return (res as unknown as { data: AdminEvent }).data;
}