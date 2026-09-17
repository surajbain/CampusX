"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { CalendarX2 } from "lucide-react";
import { listEvents } from "@/lib/api/events";
import { CategoryChips } from "@/components/event/category-chip";
import { EventCard } from "@/components/event/event-card";
import { EventGridSkeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import type { EventCategory } from "@/lib/api/types";

type Filter = EventCategory | "ALL";

export function EventsSection() {
  const [category, setCategory] = React.useState<Filter>("ALL");

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["events", category],
    queryFn: () =>
      listEvents({
        category: category === "ALL" ? undefined : category,
        limit: 12,
      }),
  });

  const events = (data as unknown as { data?: unknown[] })?.data ?? [];
  const meta = (data as unknown as { meta?: { total: number } })?.meta;

  return (
    <section id="events" className="cx-container pb-24">
      <div className="flex items-end justify-between mb-6">
        <div>
          <h2 className="text-display text-2xl md:text-3xl font-bold tracking-tight">
            Upcoming events
          </h2>
          <p className="mt-1 text-sm text-[var(--fg-muted)]">
            {meta?.total ? `${meta.total} events live now` : "Discover events across colleges"}
          </p>
        </div>
      </div>

      <CategoryChips
        value={category}
        onChange={setCategory}
        className="mb-8"
      />

      {isLoading && <EventGridSkeleton count={6} />}

      {isError && (
        <EmptyState
          icon={CalendarX2}
          title="Couldn't load events"
          description={
            error instanceof Error
              ? error.message
              : "Something went wrong. Please try again."
          }
        />
      )}

      {!isLoading && !isError && events.length === 0 && (
        <EmptyState
          icon={CalendarX2}
          title="No events found"
          description="Try a different category or check back later."
        />
      )}

      {!isLoading && !isError && events.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {(events as any[]).map((event) => (
            <EventCard key={event.id} event={event} />
          ))}
        </div>
      )}
    </section>
  );
}