"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import * as regApi from "@/lib/api/registrations";
import { ApiClientError } from "@/lib/api/client";
import { useAuthStore } from "@/lib/stores/auth-store";

// ---- Register for event ----
export function useRegisterForEvent() {
  const qc = useQueryClient();
  const setUser = useAuthStore((s) => s.setUser);

  return useMutation({
    mutationFn: ({
      eventId,
      input,
    }: {
      eventId: string;
      input: { type: regApi.RegistrationType; notes?: string };
    }) => regApi.registerForEvent(eventId, input),
    onSuccess: (reg) => {
      const message =
        reg.status === "CONFIRMED"
          ? "🎉 You're registered!"
          : "Registration submitted. Payment pending.";
      toast.success(message);
      qc.invalidateQueries({ queryKey: ["my-registrations"] });
      qc.invalidateQueries({ queryKey: ["event", reg.event_id] });
      _ = setUser; // keep reference
    },
    onError: (err) => {
      const msg =
        err instanceof ApiClientError ? err.message : "Registration failed";
      toast.error(msg);
    },
  });
}

// ---- List my registrations ----
export function useMyRegistrations(params: { page?: number; limit?: number } = {}) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: ["my-registrations", params],
    queryFn: () => regApi.listMyRegistrations(params),
    enabled: isAuthenticated,
    staleTime: 0,                    // always stale
    refetchOnMount: "always",        // always refetch on mount
    refetchOnWindowFocus: true,      // refetch on tab focus
    refetchInterval: 5000,           // poll every 5s while mounted (dev only)
  });
}

// ---- Get one ----
export function useRegistration(id: string) {
  return useQuery({
    queryKey: ["registration", id],
    queryFn: () => regApi.getRegistration(id),
    enabled: !!id,
  });
}

// ---- Cancel ----
export function useCancelRegistration() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: regApi.cancelRegistration,
    onSuccess: () => {
      toast.success("Registration cancelled");
      qc.invalidateQueries({ queryKey: ["my-registrations"] });
    },
    onError: (err) => {
      const msg =
        err instanceof ApiClientError ? err.message : "Cancellation failed";
      toast.error(msg);
    },
  });
}

// ---- List event registrations (admin) ----
export function useEventRegistrations(
  eventId: string,
  params: { page?: number; limit?: number } = {}
) {
  return useQuery({
    queryKey: ["event-registrations", eventId, params],
    queryFn: () => regApi.listEventRegistrations(eventId, params),
    enabled: !!eventId,
  });
}

// silence lint
let _: unknown = undefined;