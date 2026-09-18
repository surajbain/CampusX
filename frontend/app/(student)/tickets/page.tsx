"use client";

import * as React from "react";
import Link from "next/link";
import { TicketX, Compass } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { TicketCard } from "@/components/ticket/ticket-card";
import { useMyTickets } from "@/lib/hooks/use-tickets";
import { cn } from "@/lib/utils";


const FILTERS = [
  { value: "ALL", label: "All" },
  { value: "ACTIVE", label: "Active" },
  { value: "USED", label: "Used" },
  { value: "EXPIRED", label: "Expired" },
] as const;

export default function TicketsPage() {
  const [filter, setFilter] =
    React.useState<(typeof FILTERS)[number]["value"]>("ALL");
  const { data, isLoading } = useMyTickets({ limit: 100 });

  const tickets = ((data as any)?.data ?? []) as any[];
  const filtered =
    filter === "ALL" ? tickets : tickets.filter((t) => t.status === filter);

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-display text-3xl font-bold tracking-tight">
          My tickets
        </h1>
        <p className="mt-1.5 text-sm text-[var(--fg-muted)]">
          {tickets.length} total · {filtered.length} shown
        </p>
      </div>

      {tickets.length > 0 && (
        <div className="flex gap-2 mb-6 overflow-x-auto no-scrollbar">
          {FILTERS.map((f) => (
            <button
              key={f.value}
              onClick={() => setFilter(f.value)}
              className={cn(
                "inline-flex items-center rounded-full border px-3.5 py-1.5 text-sm font-medium whitespace-nowrap transition-all",
                filter === f.value
                  ? "border-[var(--color-brand)] bg-[var(--color-brand)]/10 text-[var(--color-brand)]"
                  : "border-[var(--border)] bg-[var(--bg-elevated)] text-[var(--fg-muted)] hover:border-[var(--border-strong)] hover:text-[var(--fg)]"
              )}
            >
              {f.label}
            </button>
          ))}
        </div>
      )}

      {isLoading && (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-32 rounded-2xl" />
          ))}
        </div>
      )}

      {!isLoading && filtered.length === 0 && (
        <EmptyState
          icon={filter === "ALL" ? Compass : TicketX}
          title={
            filter === "ALL"
              ? "No tickets yet"
              : `No ${filter.toLowerCase()} tickets`
          }
          description={
            filter === "ALL"
              ? "Register for events to get your tickets."
              : "Try changing the filter."
          }
          action={
            filter === "ALL" ? (
              <Link href="/events">
                <Button>
                  <Compass className="h-4 w-4" />
                  Browse events
                </Button>
              </Link>
            ) : undefined
          }
        />
      )}

      {!isLoading && filtered.length > 0 && (
        <div className="space-y-3">
          {filtered.map((ticket) => (
            <TicketCard key={ticket.id} ticket={ticket} />
          ))}
        </div>
      )}
    </div>
  );
}