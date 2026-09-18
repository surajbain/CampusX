import { apiFetch } from "./client";
import type { PaginationMeta } from "./types";

export interface Ticket {
  id: string;
  registration_id: string;
  user_id: string;
  event_id: string;
  ticket_code: string;
  qr_payload: string;
  qr_signature: string;
  valid_from?: string;
  valid_until?: string;
  status: "ACTIVE" | "USED" | "REVOKED" | "EXPIRED";
  issued_at: string;
  used_at?: string;

  // Event details
  event_title?: string;
  event_slug?: string;
  event_poster_url?: string;
  event_category?: string;
  event_starts_at?: string;
  event_ends_at?: string;
  event_venue?: string;
  event_city?: string;
  event_college_name?: string;
    event_whatsapp_link?: string; 

  // User details
  user_name?: string;
  user_email?: string;
}

interface ListResponse {
  data: Ticket[];
  meta: PaginationMeta;
}

function unwrap<T>(res: unknown): T {
  return (res as { data: T }).data;
}

export async function listMyTickets(params: {
  page?: number;
  limit?: number;
} = {}): Promise<ListResponse> {
  return apiFetch<ListResponse>("/api/v1/tickets/my", {
    params: params as Record<string, string | number>,
  });
}

export async function getTicket(id: string): Promise<Ticket> {
  const res = await apiFetch<Ticket>(`/api/v1/tickets/${id}`);
  return unwrap<Ticket>(res);
}
export interface Ticket {
  // ... existing fields ...
  event_whatsapp_link?: string;
}