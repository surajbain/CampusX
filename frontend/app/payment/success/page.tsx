"use client";

import * as React from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { CheckCircle2, Ticket, ArrowRight } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { usePayment } from "@/lib/hooks/use-payments";
import { formatPrice } from "@/lib/utils";

export default function PaymentSuccessPage() {
  const params = useSearchParams();
  const paymentId = params.get("id") || "";
  const { data: payment, isLoading } = usePayment(paymentId);

  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md">
        <Card className="p-8 text-center">
          {/* Success icon with animation */}
          <div className="mx-auto grid h-20 w-20 place-items-center rounded-full bg-[var(--color-success)]/15 animate-[bounce_1s_ease-in-out]">
            <CheckCircle2 className="h-10 w-10 text-[var(--color-success)]" />
          </div>

          <h1 className="mt-6 text-display text-2xl font-bold tracking-tight">
            Payment successful! 🎉
          </h1>
          <p className="mt-2 text-sm text-[var(--fg-muted)]">
            Your registration is confirmed. Ticket will be issued shortly.
          </p>

          {isLoading && (
            <div className="mt-6 text-xs text-[var(--fg-subtle)]">Loading…</div>
          )}

          {payment && (
            <div className="mt-6 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] p-4 text-left">
              <div className="text-xs text-[var(--fg-subtle)]">Amount paid</div>
              <div className="text-display text-2xl font-bold mt-0.5 text-[var(--color-success)]">
                {formatPrice(payment.amount_paise, payment.currency)}
              </div>
              <div className="mt-3 text-xs text-[var(--fg-subtle)]">
                Payment ID
              </div>
              <div className="text-xs font-mono mt-0.5 truncate">
                {payment.id}
              </div>
              {payment.gateway_payment_id && (
                <>
                  <div className="mt-3 text-xs text-[var(--fg-subtle)]">
                    Gateway Payment ID
                  </div>
                  <div className="text-xs font-mono mt-0.5 truncate">
                    {payment.gateway_payment_id}
                  </div>
                </>
              )}
            </div>
          )}

          <div className="mt-6 space-y-2">
            <Link href="/tickets" className="block">
              <Button size="lg" className="w-full">
                <Ticket className="h-4 w-4" />
                View my tickets
              </Button>
            </Link>
            <Link href="/events" className="block">
              <Button variant="ghost" size="lg" className="w-full">
                Browse more events
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </Card>

        <p className="mt-6 text-center text-xs text-[var(--fg-subtle)]">
          A confirmation email will be sent to your registered email address.
        </p>
      </div>
    </div>
  );
}