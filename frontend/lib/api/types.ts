export type EventCategory =
  | "HACKATHON"
  | "CULTURAL"
  | "SPORTS"
  | "WORKSHOP"
  | "TECH_FEST"
  | "OTHER";

export type EventStatus =
  | "DRAFT"
  | "PUBLISHED"
  | "ONGOING"
  | "COMPLETED"
  | "CANCELLED";

export interface Event {
  id: string;
  college_id: string;
  created_by?: string;
  slug: string;
  title: string;
  description?: string;
  category: EventCategory;
  poster_url?: string;
  venue?: string;
  city?: string;
  status: EventStatus;
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

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  pages: number;
}

export interface ApiSuccess<T> {
  success: true;
  data: T;
  meta?: PaginationMeta;
}

export interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError;