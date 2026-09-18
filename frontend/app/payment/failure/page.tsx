"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { XCircle, ArrowLeft, RefreshCw } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function PaymentFailurePage() {
  const params = useSearchParams();
  const reason = params.get("reason") || "Payment could not be completed";
  const regId = params.get("registration");

  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md">
        <Card className="p-8 text-center">
          <div className="mx-auto grid h-20 w-20 place-items-center rounded-full bg-[var(--color-danger)]/15">
            <XCircle className="h-10 w-10 text-[var(--color-danger)]" />
          </div>

          <h1 className="mt-6 text-display text-2xl font-bold tracking-tight">
            Payment failed
          </h1>
          <p className="mt-2 text-sm text-[var(--fg-muted)]">
            {reason}
          </p>

          <div className="mt-6 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] p-4 text-left">
            <div className="text-xs text-[var(--fg-subtle)]">
              What you can do:
            </div>
            <ul className="mt-2 text-xs space-y-1.5 text-[var(--fg-muted)]">
              <li>• Check your payment method and try again</li>
              <li>• Use a different payment method</li>
              <li>• Contact your bank if the issue persists</li>
              <li>• Your registration is still saved — just retry payment</li>
            </ul>
          </div>

          <div className="mt-6 space-y-2">
            {regId && (
              <Link href={`/tickets`} className="block">
                <Button size="lg" className="w-full">
                  <RefreshCw className="h-4 w-4" />
                  Retry payment
                </Button>
              </Link>
            )}
            <Link href="/events" className="block">
              <Button variant="ghost" size="lg" className="w-full">
                <ArrowLeft className="h-4 w-4" />
                Back to events
              </Button>
            </Link>
          </div>
        </Card>

        <p className="mt-6 text-center text-xs text-[var(--fg-subtle)]">
          If amount was deducted, it will be refunded within 5–7 business days.
        </p>
      </div>
    </div>
  );
}