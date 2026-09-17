"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import * as authApi from "@/lib/api/auth";
import { useAuthStore } from "@/lib/stores/auth-store";
import { ApiClientError } from "@/lib/api/client";

export function useLogin() {
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.user, data.access_token, data.refresh_token);
      toast.success(`Welcome back, ${data.user.full_name.split(" ")[0]}!`);
      const redirect =
        typeof window !== "undefined"
          ? new URLSearchParams(window.location.search).get("redirect")
          : null;
      router.push(redirect || "/");
      router.refresh();
    },
    onError: (err) => {
      const msg = err instanceof ApiClientError ? err.message : "Login failed";
      toast.error(msg);
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setAuth(data.user, data.access_token, data.refresh_token);
      toast.success("Account created! Welcome to CampusX 🎉");
      router.push("/");
      router.refresh();
    },
    onError: (err) => {
      const msg =
        err instanceof ApiClientError ? err.message : "Registration failed";
      toast.error(msg);
    },
  });
}

export function useLogout() {
  const router = useRouter();
  const clear = useAuthStore((s) => s.clear);
  const accessToken = useAuthStore((s) => s.accessToken);
  const refreshToken = useAuthStore((s) => s.refreshToken);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      if (accessToken && refreshToken) {
        try {
          await authApi.logout(accessToken, refreshToken);
        } catch {
          // ignore — clear locally anyway
        }
      }
    },
    onSuccess: () => {
      clear();
      queryClient.clear();
      toast.success("Logged out");
      router.push("/");
      router.refresh();
    },
  });
}

export function useCurrentUser() {
  const accessToken = useAuthStore((s) => s.accessToken);
  const setUser = useAuthStore((s) => s.setUser);

  return useQuery({
    queryKey: ["current-user", accessToken],
    queryFn: async () => {
      if (!accessToken) return null;
      const user = await authApi.me(accessToken);
      setUser(user);
      return user;
    },
    enabled: !!accessToken,
    staleTime: 5 * 60 * 1000,
  });
}