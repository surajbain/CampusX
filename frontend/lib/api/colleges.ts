import { apiFetch } from "./client";
import type { PaginationMeta } from "./types";

export interface College {
  id: string;
  name: string;
  slug: string;
  city: string;
  state?: string;
  logo_url?: string;
  website?: string;
  contact_email?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface ListResponse {
  data: College[];
  meta: PaginationMeta;
}

export interface ListCollegesParams {
  city?: string;
  page?: number;
  limit?: number;
  sort?: string;
}

export async function listColleges(
  params: ListCollegesParams = {}
): Promise<ListResponse> {
  return apiFetch<ListResponse>("/api/v1/colleges", {
    params: params as Record<string, string | number>,
  });
}

export async function getCollege(id: string): Promise<College> {
  return apiFetch<College>(`/api/v1/colleges/${id}`);
}

export async function getCollegeBySlug(slug: string): Promise<College> {
  return apiFetch<College>(`/api/v1/colleges/slug/${slug}`);
}

export async function listCities(): Promise<string[]> {
  const res = await apiFetch<{ cities: string[] }>("/api/v1/colleges/cities");
  return (res as unknown as { cities: string[] }).cities || [];
}

export async function getStats(): Promise<Record<string, number>> {
  return apiFetch<Record<string, number>>("/api/v1/colleges/stats");
}