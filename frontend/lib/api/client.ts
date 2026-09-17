"use client";

import { useAuthStore } from "@/lib/stores/auth-store";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export class ApiClientError extends Error {
  code: string;
  status: number;
  details?: Record<string, unknown>;

  constructor(
    message: string,
    code: string,
    status: number,
    details?: Record<string, unknown>
  ) {
    super(message);
    this.code = code;
    this.status = status;
    this.details = details;
  }
}

interface FetchOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined | null>;
  accessToken?: string;
  skipAuth?: boolean;
  skipRefresh?: boolean;
}

let refreshPromise: Promise<string | null> | null = null;

async function attemptRefresh(): Promise<string | null> {
  const { refreshToken, setTokens, clear } = useAuthStore.getState();
  if (!refreshToken) return null;

  try {
    const res = await fetch(`${API_URL}/api/v1/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    if (!res.ok) {
      clear();
      return null;
    }
    const body = await res.json();
    const newAccess = body?.data?.access_token as string | undefined;
    const newRefresh = body?.data?.refresh_token as string | undefined;
    if (!newAccess || !newRefresh) {
      clear();
      return null;
    }
    setTokens(newAccess, newRefresh);
    return newAccess;
  } catch {
    clear();
    return null;
  }
}

export async function apiFetch<T>(
  path: string,
  options: FetchOptions = {}
): Promise<T> {
  const {
    params,
    accessToken,
    skipAuth,
    skipRefresh,
    headers,
    ...rest
  } = options;

  const url = new URL(path, API_URL);
  if (params) {
    Object.entries(params).forEach(([k, v]) => {
      if (v !== undefined && v !== null && v !== "") {
        url.searchParams.set(k, String(v));
      }
    });
  }

  // Resolve token: explicit > store
  const storeToken = skipAuth ? undefined : useAuthStore.getState().accessToken;
  const token = accessToken ?? storeToken;

  const doFetch = async (bearer?: string) =>
    fetch(url.toString(), {
      ...rest,
      headers: {
        "Content-Type": "application/json",
        ...(bearer ? { Authorization: `Bearer ${bearer}` } : {}),
        ...headers,
      },
    });

  let res = await doFetch(token);

  // Auto refresh on 401
  if (res.status === 401 && !skipAuth && !skipRefresh && !accessToken) {
    if (!refreshPromise) {
      refreshPromise = attemptRefresh().finally(() => {
        refreshPromise = null;
      });
    }
    const newToken = await refreshPromise;
    if (newToken) {
      res = await doFetch(newToken);
    }
  }

  // Parse JSON safely
  let body: unknown;
  try {
    body = await res.json();
  } catch {
    throw new ApiClientError(
      "Invalid response from server",
      "INVALID_RESPONSE",
      res.status
    );
  }

  if (!res.ok) {
    const err = body as {
      error?: { code: string; message: string; details?: Record<string, unknown> };
    };
    throw new ApiClientError(
      err.error?.message || "Request failed",
      err.error?.code || "UNKNOWN",
      res.status,
      err.error?.details
    );
  }

  const ok = body as { success: true; data: T; meta?: unknown };
  return { data: ok.data, meta: ok.meta } as T & { meta?: unknown } as T;
}

export const API_BASE = API_URL;