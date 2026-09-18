"use client";

import { useQuery } from "@tanstack/react-query";
import * as ticketApi from "@/lib/api/tickets";
import { useAuthStore } from "@/lib/stores/auth-store";

export function useMyTickets(params: { page?: number; limit?: number } = {}) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: ["my-tickets", params],
    queryFn: () => ticketApi.listMyTickets(params),
    enabled: isAuthenticated,
    staleTime: 30 * 1000,
    refetchOnMount: "always",
  });
}

export function useTicket(id: string) {
  return useQuery({
    queryKey: ["ticket", id],
    queryFn: () => ticketApi.getTicket(id),
    enabled: !!id,
    staleTime: 60 * 1000,
  });
}