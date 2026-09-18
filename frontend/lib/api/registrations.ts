import { apiFetch } from "./client";
import type { PaginationMeta } from "./types";

export type RegistrationType = "ATTEND" | "PARTICIPATE";
export type RegistrationStatus =
  | "PENDING"
  | "CONFIRMED"
  | "CANCELLED"
  | "WAITLISTED";

export interface Registration {
  id: string;
  event_id: string;
  user_id: string;
  college_id: string;
  type: RegistrationType;
  status: RegistrationStatus;
  notes?: string;

  event_title?: string;
  event_slug?: string;
  event_poster_url?: string;
  event_category?: string;
  event_starts_at?: string;
  event_ends_at?: string;
  event_venue?: string;
  event_city?: string;
  event_price_paise: number;
  event_currency?: string;
  event_college_name?: string;

  user_name?: string;
  user_email?: string;

  created_at: string;
  updated_at: string;
}

interface ListResponse {
  data: Registration[];
  meta: PaginationMeta;
}

// apiFetch returns {data, meta} — unwrap
function unwrap<T>(res: unknown): T {
  return (res as { data: T }).data;
}

// ---- Register for event ----
export async function registerForEvent(
  eventId: string,
  input: { type: RegistrationType; notes?: string }
): Promise<Registration> {
  const res = await apiFetch<Registration>(
    `/api/v1/events/${eventId}/register`,
    {
      method: "POST",
      body: JSON.stringify(input),
    }
  );
  return unwrap<Registration>(res);
}

// ---- List user's registrations ----
export async function listMyRegistrations(params: {
  page?: number;
  limit?: number;
} = {}): Promise<ListResponse> {
  return apiFetch<ListResponse>("/api/v1/registrations/my", {
    params: params as Record<string, string | number>,
  });
}

// ---- Get one ----
export async function getRegistration(id: string): Promise<Registration> {
  const res = await apiFetch<Registration>(`/api/v1/registrations/${id}`);
  return unwrap<Registration>(res);
}

// ---- Cancel ----
export async function cancelRegistration(id: string): Promise<void> {
  await apiFetch(`/api/v1/registrations/${id}`, { method: "DELETE" });
}

// ---- List event registrations (admin/organizer) ----
export async function listEventRegistrations(
  eventId: string,
  params: { page?: number; limit?: number } = {}
): Promise<ListResponse> {
  return apiFetch<ListResponse>(`/api/v1/events/${eventId}/registrations`, {
    params: params as Record<string, string | number>,
  });
}