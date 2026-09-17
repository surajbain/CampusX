"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import * as adminApi from "@/lib/api/admin";
import { ApiClientError } from "@/lib/api/client";

// ---- Colleges ----

export function useCreateCollege() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: adminApi.createCollege,
    onSuccess: () => {
      toast.success("College created");
      qc.invalidateQueries({ queryKey: ["colleges"] });
    },
    onError: (err) =>
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to create college"
      ),
  });
}

export function useUpdateCollege() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: adminApi.UpdateCollegeInput }) =>
      adminApi.updateCollege(id, input),
    onSuccess: (_, vars) => {
      toast.success("College updated");
      qc.invalidateQueries({ queryKey: ["colleges"] });
      qc.invalidateQueries({ queryKey: ["college", vars.id] });
    },
    onError: (err) =>
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to update college"
      ),
  });
}

export function useDeleteCollege() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: adminApi.deleteCollege,
    onSuccess: () => {
      toast.success("College deleted");
      qc.invalidateQueries({ queryKey: ["colleges"] });
    },
    onError: (err) =>
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to delete college"
      ),
  });
}

export function useCollege(id: string) {
  return useQuery({
    queryKey: ["college", id],
    queryFn: () => adminApi.getCollege(id),
    enabled: !!id,
  });
}

// ---- Users ----

export function useCollegeUsers(
  collegeId: string,
  params: { page?: number; limit?: number } = {}
) {
  return useQuery({
    queryKey: ["college-users", collegeId, params],
    queryFn: () => adminApi.listCollegeUsers(collegeId, params),
    enabled: !!collegeId,
  });
}

export function useCreateCollegeUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ collegeId, input }: { collegeId: string; input: adminApi.CreateUserInput }) =>
      adminApi.createCollegeUser(collegeId, input),
    onSuccess: (_, vars) => {
      toast.success("User created");
      qc.invalidateQueries({ queryKey: ["college-users", vars.collegeId] });
    },
    onError: (err) =>
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to create user"
      ),
  });
}

export function useUpdateCollegeUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      collegeId,
      userId,
      input,
    }: {
      collegeId: string;
      userId: string;
      input: adminApi.UpdateUserInput;
    }) => adminApi.updateCollegeUser(collegeId, userId, input),
    onSuccess: (_, vars) => {
      toast.success("User updated");
      qc.invalidateQueries({ queryKey: ["college-users", vars.collegeId] });
    },
    onError: (err) =>
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to update user"
      ),
  });
}

export function useDeleteCollegeUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ collegeId, userId }: { collegeId: string; userId: string }) =>
      adminApi.deleteCollegeUser(collegeId, userId),
    onSuccess: (_, vars) => {
      toast.success("User deleted");
      qc.invalidateQueries({ queryKey: ["college-users", vars.collegeId] });
    },
    onError: (err) =>
      toast.error(
        err instanceof ApiClientError ? err.message : "Failed to delete user"
      ),
  });
}

// ---- Stats ----

export function useAdminStats() {
  return useQuery({
    queryKey: ["admin-stats"],
    queryFn: adminApi.getAdminStats,
    staleTime: 30 * 1000,
  });
}