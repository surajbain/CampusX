"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { Search, SlidersHorizontal, X, CalendarX2, Loader2 } from "lucide-react";
import { listEvents } from "@/lib/api/events";
import { listCities } from "@/lib/api/colleges";
import { CategoryChips } from "@/components/event/category-chip";
import { EventCard } from "@/components/event/event-card";
import { EventGridSkeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { EventCategory } from "@/lib/api/types";

type Filter = EventCategory | "ALL";

export function EventsBrowser() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [category, setCategory] = React.useState<Filter>(
    (searchParams.get("category") as Filter) || "ALL"
  );
  const [city, setCity] = React.useState(searchParams.get("city") || "");
  const [query, setQuery] = React.useState(searchParams.get("q") || "");
  const [debouncedQuery, setDebouncedQuery] = React.useState(query);
  const [page, setPage] = React.useState(1);
  const [showFilters, setShowFilters] = React.useState(false);

  // Sync URL
  React.useEffect(() => {
    const params = new URLSearchParams();
    if (category !== "ALL") params.set("category", category);
    if (city) params.set("city", city);
    if (debouncedQuery) params.set("q", debouncedQuery);
    if (page > 1) params.set("page", String(page));
    const qs = params.toString();
    router.replace(qs ? `/events?${qs}` : "/events", { scroll: false });
  }, [category, city, debouncedQuery, page, router]);

  // Debounce query
  React.useEffect(() => {
    const t = setTimeout(() => setDebouncedQuery(query), 350);
    return () => clearTimeout(t);
  }, [query]);

  // Reset page on filter change
  React.useEffect(() => {
    setPage(1);
  }, [category, city, debouncedQuery]);

  const { data, isLoading, isError, error, isFetching } = useQuery({
    queryKey: ["events", { category, city, q: debouncedQuery, page }],
    queryFn: () =>
      listEvents({
        category: category === "ALL" ? undefined : category,
        city: city || undefined,
        q: debouncedQuery || undefined,
        page,
        limit: 12,
      }),
  });

  const citiesQuery = useQuery({
    queryKey: ["cities"],
    queryFn: listCities,
    staleTime: 10 * 60 * 1000,
  });

  const events = ((data as any)?.data ?? []) as any[];
  const meta = (data as any)?.meta as
    | { total: number; page: number; pages: number }
    | undefined;

  const activeFilters = [category !== "ALL", city, debouncedQuery].filter(Boolean).length;

  const clearFilters = () => {
    setCategory("ALL");
    setCity("");
    setQuery("");
    setPage(1);
  };

  return (
    <div className="cx-container py-8 md:py-12">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-display text-3xl md:text-5xl font-bold tracking-tight">
          Discover events
        </h1>
        <p className="mt-2 text-sm text-[var(--fg-muted)]">
          {meta?.total
            ? `${meta.total} events across colleges`
            : "Browse events across colleges"}
        </p>
      </div>

      {/* Search + filter toggle */}
      <div className="flex gap-2 mb-4">
        <div className="relative flex-1">
          <Search className="absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[var(--fg-subtle)]" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search events by name, college, city…"
            className="w-full h-12 rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)] pl-10 pr-10 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
          />
          {query && (
            <button
              onClick={() => setQuery("")}
              className="absolute right-3 top-1/2 -translate-y-1/2 grid h-6 w-6 place-items-center rounded-full hover:bg-[var(--bg-muted)]"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          )}
        </div>

        <Button
          variant="secondary"
          onClick={() => setShowFilters((v) => !v)}
          className="h-12 shrink-0 gap-2"
        >
          <SlidersHorizontal className="h-4 w-4" />
          <span className="hidden sm:inline">Filters</span>
          {activeFilters > 0 && (
            <span className="grid h-5 min-w-5 place-items-center rounded-full bg-[var(--color-brand)] px-1.5 text-[10px] font-bold text-white">
              {activeFilters}
            </span>
          )}
        </Button>
      </div>

      {/* Categories */}
      <CategoryChips value={category} onChange={setCategory} className="mb-4" />

      {/* City filter */}
      {showFilters && (
        <div className="mb-6 p-4 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)]">
          <div className="flex flex-col sm:flex-row sm:items-center gap-3">
            <label className="text-sm font-medium">City</label>
            <select
              value={city}
              onChange={(e) => setCity(e.target.value)}
              className="h-10 rounded-lg border border-[var(--border)] bg-[var(--bg-elevated)] px-3 text-sm focus:outline-none focus:border-[var(--color-brand)] min-w-[200px]"
            >
              <option value="">All cities</option>
              {citiesQuery.data?.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>

            {activeFilters > 0 && (
              <Button variant="ghost" size="sm" onClick={clearFilters} className="sm:ml-auto">
                <X className="h-4 w-4" />
                Clear filters
              </Button>
            )}
          </div>
        </div>
      )}

      {/* Results */}
      <div className="mt-6">
        {isLoading && <EventGridSkeleton count={6} />}

        {isError && (
          <EmptyState
            icon={CalendarX2}
            title="Couldn't load events"
            description={
              error instanceof Error ? error.message : "Something went wrong."
            }
            action={<Button onClick={() => window.location.reload()}>Retry</Button>}
          />
        )}

        {!isLoading && !isError && events.length === 0 && (
          <EmptyState
            icon={CalendarX2}
            title="No events found"
            description="Try adjusting your filters or search query."
            action={
              activeFilters > 0 ? (
                <Button onClick={clearFilters}>Clear filters</Button>
              ) : undefined
            }
          />
        )}

        {!isLoading && !isError && events.length > 0 && (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
              {events.map((event) => (
                <EventCard key={event.id} event={event} />
              ))}
            </div>

            {/* Pagination */}
            {meta && meta.pages > 1 && (
              <div className="mt-12 flex items-center justify-center gap-2">
                <Button
                  variant="secondary"
                  disabled={page === 1}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  Previous
                </Button>
                <span className="px-4 text-sm text-[var(--fg-muted)]">
                  Page {page} of {meta.pages}
                </span>
                <Button
                  variant="secondary"
                  disabled={page >= meta.pages}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Next
                </Button>
              </div>
            )}

            {isFetching && (
              <div className="mt-6 flex items-center justify-center gap-2 text-sm text-[var(--fg-muted)]">
                <Loader2 className="h-4 w-4 animate-spin" />
                Updating…
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}