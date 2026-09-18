"use client";

import * as React from "react";
import Link from "next/link";
import {
  Calendar,
  MapPin,
  X,
  CheckCircle2,
  Clock,
  Ban,
  ExternalLink,
} from "lucide-react";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { formatPrice, formatDateRange, cn } from "@/lib/utils";
import { categoryMeta } from "@/lib/constants";
import { useCancelRegistration } from "@/lib/hooks/use-registrations";
import type { Registration } from "@/lib/api/registrations";

export function RegistrationCard({ registration }: { registration: Registration }) {
  const cancelMutation = useCancelRegistration();
  const cat = categoryMeta(registration.event_category);
  const dates = formatDateRange(
    registration.event_starts_at,
    registration.event_ends_at
  );
  const price = formatPrice(registration.event_price_paise);

  const statusConfig: Record<
    string,
    { label: string; variant: "success" | "warning" | "danger" | "default"; icon: any }
  > = {
    CONFIRMED: { label: "Confirmed", variant: "success", icon: CheckCircle2 },
    PENDING: { label: "Payment pending", variant: "warning", icon: Clock },
    CANCELLED: { label: "Cancelled", variant: "danger", icon: Ban },
    WAITLISTED: { label: "Waitlisted", variant: "default", icon: Clock },
  };

  const status = statusConfig[registration.status] ?? statusConfig.PENDING;
  const StatusIcon = status.icon;
  const canCancel =
    registration.status === "CONFIRMED" || registration.status === "PENDING";

  return (
    <Card className="overflow-hidden">
      <div className="flex flex-col sm:flex-row">
        {/* Poster */}
        <div
          className="relative w-full sm:w-48 shrink-0 aspect-[16/9] sm:aspect-auto overflow-hidden"
          style={{
            background: `linear-gradient(135deg, ${cat.color}30, ${cat.color}10)`,
          }}
        >
          {registration.event_poster_url ? (
            <img
              src={registration.event_poster_url}
              alt={registration.event_title}
              className="h-full w-full object-cover"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center text-4xl">
              {cat.emoji}
            </div>
          )}
        </div>

        {/* Body */}
        <div className="flex-1 p-5">
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0 flex-1">
              <div className="flex flex-wrap items-center gap-2 mb-2">
                <Badge variant={status.variant}>
                  <StatusIcon className="h-3 w-3" />
                  {status.label}
                </Badge>
                <Badge
                  variant={
                    registration.type === "PARTICIPATE" ? "brand" : "outline"
                  }
                >
                  {registration.type === "PARTICIPATE" ? "Participating" : "Attending"}
                </Badge>
              </div>

              <Link
                href={`/events/${registration.event_id}`}
                className="block group"
              >
                <h3 className="text-display text-base font-semibold leading-snug line-clamp-2 group-hover:text-[var(--color-brand)] transition-colors">
                  {registration.event_title}
                </h3>
              </Link>

              <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1.5 text-xs text-[var(--fg-muted)]">
                {dates && (
                  <span className="inline-flex items-center gap-1">
                    <Calendar className="h-3.5 w-3.5" />
                    {dates}
                  </span>
                )}
                {registration.event_venue && (
                  <span className="inline-flex items-center gap-1">
                    <MapPin className="h-3.5 w-3.5" />
                    {registration.event_venue}
                  </span>
                )}
              </div>

              <div className="mt-2 text-xs text-[var(--fg-subtle)]">
                {registration.event_college_name}
              </div>
            </div>

            <div className="text-right shrink-0">
              <div
                className={cn(
                  "text-display text-lg font-bold",
                  registration.event_price_paise === 0 &&
                    "text-[var(--color-success)]"
                )}
              >
                {price}
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="mt-4 pt-4 border-t border-[var(--border)] flex gap-2">
            <Link href={`/events/${registration.event_id}`} className="flex-1">
              <Button variant="secondary" size="sm" className="w-full">
                <ExternalLink className="h-3.5 w-3.5" />
                View event
              </Button>
            </Link>

            {canCancel && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  if (confirm("Cancel your registration for this event?")) {
                    cancelMutation.mutate(registration.id);
                  }
                }}
                disabled={cancelMutation.isPending}
              >
                <X className="h-3.5 w-3.5" />
                Cancel
              </Button>
            )}
          </div>
        </div>
      </div>
    </Card>
  );
}