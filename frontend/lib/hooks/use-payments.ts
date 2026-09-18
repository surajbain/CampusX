"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import * as payApi from "@/lib/api/payments";
import { ApiClientError } from "@/lib/api/client";

// ---- Create payment ----
export function useCreatePayment() {
  return useMutation({
    mutationFn: payApi.createPayment,
    onError: (err) => {
      const msg =
        err instanceof ApiClientError ? err.message : "Payment creation failed";
      toast.error(msg);
    },
  });
}

// ---- Verify payment ----
export function useVerifyPayment() {
  const qc = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: payApi.verifyPayment,
    onSuccess: (payment) => {
      qc.invalidateQueries({ queryKey: ["my-registrations"] });
      qc.invalidateQueries({ queryKey: ["registration"] });
      toast.success("🎉 Payment successful!");
      router.push(`/payment/success?id=${payment.id}`);
    },
    onError: (err) => {
      const msg =
        err instanceof ApiClientError ? err.message : "Payment verification failed";
      toast.error(msg);
    },
  });
}

// ---- Get payment ----
export function usePayment(id: string) {
  return useQuery({
    queryKey: ["payment", id],
    queryFn: () => payApi.getPayment(id),
    enabled: !!id,
  });
}