import { apiFetch } from "./client";
import type { Event, PaginationMeta } from "./types";

interface ListResponse {
  data: Event[];
  meta: PaginationMeta;
}

export interface ListEventsParams {
  category?: string;
  city?: string;
  status?: string;
  q?: string;
  featured?: boolean;
  college_id?: string;
  page?: number;
  limit?: number;
  sort?: string;
}

/**
 * Fetch public events (published only).
 */
export async function listEvents(
  params: ListEventsParams = {}
): Promise<ListResponse> {
  return apiFetch<ListResponse>("/api/v1/events", {
    params: params as Record<string, string | number | boolean>,
  });
}

/**
 * Fetch featured events.
 */
export async function listFeaturedEvents(
  limit = 8
): Promise<ListResponse> {
  return apiFetch<ListResponse>("/api/v1/events/featured", {
    params: { limit },
  });
}

/**
 * Fetch single event by ID.
 */
export async function getEvent(id: string): Promise<Event> {
  return apiFetch<Event>(`/api/v1/events/${id}`);
}

/**
 * Fetch categories list.
 */
export async function listCategories(): Promise<string[]> {
  const res = await apiFetch<{ categories: string[] }>(
    "/api/v1/events/categories"
  );
  return (res as unknown as { categories: string[] }).categories || [];
}