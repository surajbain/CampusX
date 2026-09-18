"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Loader2, CheckCircle2, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/lib/stores/auth-store";
import {
  useRegisterForEvent,
  useMyRegistrations,
} from "@/lib/hooks/use-registrations";
import { CheckoutButton } from "@/components/payment/checkout-button";
import { cn } from "@/lib/utils";

interface Props {
  eventId: string;
  pricePaise: number;
  allowTeams: boolean;
  size?: "md" | "lg";
  className?: string;
}

export function RegisterButton({
  eventId,
  pricePaise,
  allowTeams,
  size = "lg",
  className,
}: Props) {
  const router = useRouter();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const user = useAuthStore((s) => s.user);
  const registerMutation = useRegisterForEvent();
  const [showTypeDialog, setShowTypeDialog] = React.useState(false);

  const { data: myRegs } = useMyRegistrations({ limit: 100 });
  const existingReg = ((myRegs as any)?.data ?? []).find(
    (r: any) =>
      r.event_id === eventId &&
      (r.status === "CONFIRMED" || r.status === "PENDING")
  );

  const isPaid = pricePaise > 0;

  const handleClick = () => {
    if (!isAuthenticated) {
      router.push(`/login?redirect=/events/${eventId}`);
      return;
    }

    if (user?.role !== "STUDENT") {
      alert("Only students can register for events.");
      return;
    }

    // For paid events with type choice, ask
    if (allowTeams) {
      setShowTypeDialog(true);
    } else {
      registerMutation.mutate({ eventId, input: { type: "ATTEND" } });
    }
  };

  const handleTypeSelect = (type: "ATTEND" | "PARTICIPATE") => {
    setShowTypeDialog(false);
    registerMutation.mutate({ eventId, input: { type } });
  };

  // ---------- Already registered ----------
  if (existingReg) {
    // If confirmed → show "Registered"
    if (existingReg.status === "CONFIRMED") {
      return (
        <Button
          size={size}
          variant="secondary"
          className={cn("w-full", className)}
          disabled
        >
          <CheckCircle2 className="h-4 w-4 text-[var(--color-success)]" />
          Registered
        </Button>
      );
    }

    // If pending → show checkout button (for paid events)
    if (existingReg.status === "PENDING") {
      if (isPaid) {
        return (
          <div className={cn("space-y-2", className)}>
            <div className="text-xs text-center text-[var(--fg-muted)]">
              Registration pending — complete payment to confirm
            </div>
            <CheckoutButton
              registrationId={existingReg.id}
              amountPaise={pricePaise}
            />
          </div>
        );
      }
      return (
        <Button
          size={size}
          variant="outline"
          className={cn("w-full", className)}
          disabled
        >
          <CheckCircle2 className="h-4 w-4" />
          Pending confirmation
        </Button>
      );
    }
  }

  // ---------- Not registered yet ----------
  return (
    <>
      <Button
        size={size}
        className={cn("w-full", className)}
        onClick={handleClick}
        disabled={registerMutation.isPending}
      >
        {registerMutation.isPending ? (
          <>
            <Loader2 className="h-4 w-4 animate-spin" />
            Registering…
          </>
        ) : (
          "Register now"
        )}
      </Button>

      {/* Type selection dialog */}
      {showTypeDialog && (
        <div
          className="fixed inset-0 z-50 grid place-items-center bg-black/50 p-4"
          onClick={() => setShowTypeDialog(false)}
        >
          <div
            className="w-full max-w-md rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] p-6 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-display text-xl font-bold tracking-tight">
              How do you want to join?
            </h3>
            <p className="mt-1 text-sm text-[var(--fg-muted)]">
              Choose your registration type
            </p>

            <div className="mt-5 space-y-3">
              <button
                onClick={() => handleTypeSelect("ATTEND")}
                className="w-full rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] p-4 text-left hover:border-[var(--color-brand)] transition-colors"
              >
                <div className="font-medium">Attend only</div>
                <div className="mt-1 text-xs text-[var(--fg-muted)]">
                  Watch, learn, network — no participation
                </div>
              </button>

              <button
                onClick={() => handleTypeSelect("PARTICIPATE")}
                className="w-full rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] p-4 text-left hover:border-[var(--color-brand)] transition-colors"
              >
                <div className="font-medium flex items-center gap-2">
                  Participate
                  <Users className="h-3.5 w-3.5 text-[var(--color-brand)]" />
                </div>
                <div className="mt-1 text-xs text-[var(--fg-muted)]">
                  Compete — requires team formation
                </div>
              </button>
            </div>

            <Button
              variant="ghost"
              size="sm"
              className="mt-4 w-full"
              onClick={() => setShowTypeDialog(false)}
            >
              Cancel
            </Button>
          </div>
        </div>
      )}
    </>
  );
}