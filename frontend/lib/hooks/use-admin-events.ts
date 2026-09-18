"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import * as adminApi from "@/lib/api/admin";
import { ApiClientError } from "@/lib/api/client";

// ---------- List admin's events ----------
export function useAdminEvents(params: {
  page?: number;
  limit?: number;
  status?: string;
} = {}) {
  return useQuery({
    queryKey: ["admin-events", params],
    queryFn: () => adminApi.listAdminEvents(params),
    staleTime: 30 * 1000,
  });
}

// ---------- Single event ----------
export function useAdminEvent(id: string) {
  return useQuery({
    queryKey: ["admin-event", id],
    queryFn: () => adminApi.getAdminEvent(id),
    enabled: !!id,
  });
}

// ---------- Create ----------
export function useCreateEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: adminApi.createEvent,
    onSuccess: (event) => {
      toast.success(`Event created: ${event.title}`);
      qc.invalidateQueries({ queryKey: ["admin-events"] });
      qc.invalidateQueries({ queryKey: ["events"] });
    },
    onError: (err) => {
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to create event"
      );
    },
  });
}

// ---------- Update ----------
export function useUpdateEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: Partial<adminApi.CreateEventInput> }) =>
      adminApi.updateEvent(id, input),
    onSuccess: (_, vars) => {
      toast.success("Event updated");
      qc.invalidateQueries({ queryKey: ["admin-events"] });
      qc.invalidateQueries({ queryKey: ["admin-event", vars.id] });
      qc.invalidateQueries({ queryKey: ["events"] });
      qc.invalidateQueries({ queryKey: ["event", vars.id] });
    },
    onError: (err) => {
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to update event"
      );
    },
  });
}

// ---------- Delete ----------
export function useDeleteEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: adminApi.deleteEvent,
    onSuccess: () => {
      toast.success("Event deleted");
      qc.invalidateQueries({ queryKey: ["admin-events"] });
      qc.invalidateQueries({ queryKey: ["events"] });
    },
    onError: (err) => {
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to delete event"
      );
    },
  });
}

// ---------- Publish ----------
export function usePublishEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: adminApi.publishEvent,
    onSuccess: () => {
      toast.success("Event published! 🎉");
      qc.invalidateQueries({ queryKey: ["admin-events"] });
      qc.invalidateQueries({ queryKey: ["events"] });
    },
    onError: (err) => {
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to publish"
      );
    },
  });
}

// ---------- Cancel ----------
export function useCancelEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: adminApi.cancelEvent,
    onSuccess: () => {
      toast.success("Event cancelled");
      qc.invalidateQueries({ queryKey: ["admin-events"] });
      qc.invalidateQueries({ queryKey: ["events"] });
    },
    onError: (err) => {
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to cancel"
      );
    },
  });
}