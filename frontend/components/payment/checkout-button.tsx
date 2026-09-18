"use client";

import * as React from "react";
import { Loader2, CreditCard, Lock } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCreatePayment, useVerifyPayment } from "@/lib/hooks/use-payments";
import { createHmac } from "crypto";
import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

interface Props {
  registrationId: string;
  amountPaise: number;
  onSuccess?: () => void;
}

// ---- Razorpay global types ----
declare global {
  interface Window {
    Razorpay?: new (options: RazorpayOptions) => RazorpayInstance;
  }
}

interface RazorpayOptions {
  key: string;
  amount: number;
  currency: string;
  name: string;
  description: string;
  order_id: string;
  handler: (response: RazorpayResponse) => void;
  prefill?: { name?: string; email?: string; contact?: string };
  theme?: { color?: string };
  modal?: { ondismiss?: () => void };
}

interface RazorpayInstance {
  open: () => void;
  on: (event: string, handler: () => void) => void;
}

interface RazorpayResponse {
  razorpay_payment_id: string;
  razorpay_order_id: string;
  razorpay_signature: string;
}

export function CheckoutButton({ registrationId, amountPaise, onSuccess }: Props) {
  const createMutation = useCreatePayment();
  const verifyMutation = useVerifyPayment();
  const router = useRouter();
  const qc = useQueryClient();
  const [isProcessing, setIsProcessing] = React.useState(false);

  const handlePay = async () => {
    setIsProcessing(true);
    try {
      // 1. Create payment on backend
      let payment;
      try {
        payment = await createMutation.mutateAsync(registrationId);
      } catch (err: any) {
        // Handle "already paid" — user already completed payment
        if (err?.code === "CONFLICT" && err?.message?.toLowerCase().includes("already paid")) {
          // Force refresh registration state
          await qc.invalidateQueries({ queryKey: ["my-registrations"] });
          await qc.refetchQueries({ queryKey: ["my-registrations"], type: "all" });
          // Redirect to tickets — the registration is confirmed
          router.push("/tickets");
          return;
        }
        throw err;
      }

      // 2. Mock provider — auto-verify
      if (payment.provider === "MOCK") {
        const mockPaymentId = "mock_pay_" + Math.random().toString(36).slice(2, 10);
        const signature = await mockSign(mockPaymentId);

        try {
          await verifyMutation.mutateAsync({
            payment_id: payment.payment_id,
            gateway_order_id: payment.gateway_order_id,
            gateway_payment_id: mockPaymentId,
            signature,
          });
        } catch (err: any) {
          // If verify says "already paid" — same handling
          if (err?.message?.toLowerCase().includes("already paid")) {
            await qc.invalidateQueries({ queryKey: ["my-registrations"] });
            await qc.refetchQueries({ queryKey: ["my-registrations"], type: "all" });
            router.push("/tickets");
            return;
          }
          throw err;
        }
        return;
      }

      // 3. Razorpay provider
      if (payment.provider === "RAZORPAY" && payment.public_key) {
        await loadRazorpayScript();
        if (!window.Razorpay) throw new Error("Razorpay SDK failed to load");

        const rzp = new window.Razorpay({
          key: payment.public_key,
          amount: payment.amount_paise,
          currency: payment.currency,
          name: "CampusX",
          description: "Event registration",
          order_id: payment.gateway_order_id,
          handler: async (response) => {
            try {
              await verifyMutation.mutateAsync({
                payment_id: payment.payment_id,
                gateway_order_id: response.razorpay_order_id,
                gateway_payment_id: response.razorpay_payment_id,
                signature: response.razorpay_signature,
              });
            } catch (err) {
              console.error("Verify failed:", err);
            }
          },
          theme: { color: "#7C3AED" },
          modal: {
            ondismiss: () => setIsProcessing(false),
          },
        });
        rzp.open();
        return;
      }

      throw new Error("Unknown payment provider");
    } catch (err) {
      console.error("Payment flow error:", err);
      setIsProcessing(false);
    }
  };

  return (
    <Button
      size="lg"
      className="w-full"
      onClick={handlePay}
      disabled={isProcessing}
    >
      {isProcessing ? (
        <>
          <Loader2 className="h-4 w-4 animate-spin" />
          Processing…
        </>
      ) : (
        <>
          <CreditCard className="h-4 w-4" />
          Pay ₹{(amountPaise / 100).toFixed(0)}
        </>
      )}
    </Button>
  );
}

// ---- Helpers ----

async function loadRazorpayScript(): Promise<void> {
  if (window.Razorpay) return;
  return new Promise((resolve, reject) => {
    const script = document.createElement("script");
    script.src = "https://checkout.razorpay.com/v1/checkout.js";
    script.onload = () => resolve();
    script.onerror = () => reject(new Error("Failed to load Razorpay SDK"));
    document.body.appendChild(script);
  });
}

// Mock signature — mirrors backend's MockProvider.SignPayment
async function mockSign(paymentId: string): Promise<string> {
  const secret = "mock_secret_dev_only";
  const encoder = new TextEncoder();
  const key = await crypto.subtle.importKey(
    "raw",
    encoder.encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"]
  );
  const signature = await crypto.subtle.sign("HMAC", key, encoder.encode(paymentId));
  return Array.from(new Uint8Array(signature))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}