"use client";

import * as React from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { MapPin, Building2 } from "lucide-react";
import { TopNav } from "@/components/layout/top-nav";
import { BottomNav } from "@/components/layout/bottom-nav";
import { Footer } from "@/components/layout/footer";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { Button } from "@/components/ui/button";
import { listColleges, listCities } from "@/lib/api/colleges";
import { cn } from "@/lib/utils";

export default function CollegesPage() {
  const [city, setCity] = React.useState("");

  const { data, isLoading, isError } = useQuery({
    queryKey: ["colleges", { city }],
    queryFn: () => listColleges({ city: city || undefined, limit: 50 }),
  });

  const citiesQuery = useQuery({
    queryKey: ["cities"],
    queryFn: listCities,
    staleTime: 10 * 60 * 1000,
  });

  const colleges = ((data as any)?.data ?? []) as any[];

  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />

      <main className="flex-1 pb-20 md:pb-0">
        <div className="cx-container py-8 md:py-12">
          <div className="mb-8">
            <h1 className="text-display text-3xl md:text-5xl font-bold tracking-tight">
              Colleges on CampusX
            </h1>
            <p className="mt-2 text-sm text-[var(--fg-muted)]">
              Discover events from colleges near you
            </p>
          </div>

          {citiesQuery.data && citiesQuery.data.length > 0 && (
            <div className="flex gap-2 overflow-x-auto no-scrollbar mb-8 pb-1">
              <button
                onClick={() => setCity("")}
                className={cn(
                  "inline-flex items-center gap-1.5 rounded-full border px-3.5 py-1.5 text-sm font-medium whitespace-nowrap transition-all",
                  city === ""
                    ? "border-[var(--color-brand)] bg-[var(--color-brand)]/10 text-[var(--color-brand)]"
                    : "border-[var(--border)] bg-[var(--bg-elevated)] text-[var(--fg-muted)] hover:border-[var(--border-strong)] hover:text-[var(--fg)]"
                )}
              >
                All Cities
              </button>
              {citiesQuery.data.map((c) => (
                <button
                  key={c}
                  onClick={() => setCity(c)}
                  className={cn(
                    "inline-flex items-center gap-1.5 rounded-full border px-3.5 py-1.5 text-sm font-medium whitespace-nowrap transition-all",
                    city === c
                      ? "border-[var(--color-brand)] bg-[var(--color-brand)]/10 text-[var(--color-brand)]"
                      : "border-[var(--border)] bg-[var(--bg-elevated)] text-[var(--fg-muted)] hover:border-[var(--border-strong)] hover:text-[var(--fg)]"
                  )}
                >
                  <MapPin className="h-3.5 w-3.5" />
                  {c}
                </button>
              ))}
            </div>
          )}

          {isLoading && (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
              {Array.from({ length: 6 }).map((_, i) => (
                <Skeleton key={i} className="h-32 rounded-2xl" />
              ))}
            </div>
          )}

          {isError && (
            <EmptyState
              icon={Building2}
              title="Couldn't load colleges"
              action={
                <Button onClick={() => window.location.reload()}>Retry</Button>
              }
            />
          )}

          {!isLoading && !isError && colleges.length === 0 && (
            <EmptyState
              icon={Building2}
              title="No colleges found"
              description={
                city
                  ? `No colleges in ${city} yet.`
                  : "No colleges listed yet."
              }
            />
          )}

          {!isLoading && !isError && colleges.length > 0 && (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
              {colleges.map((c) => (
                <Link key={c.id} href={`/colleges/${c.slug}`}>
                  <Card interactive className="h-full p-5">
                    <div className="flex items-start gap-3">
                      <div className="grid h-12 w-12 shrink-0 place-items-center rounded-xl bg-brand-gradient text-white font-bold text-lg">
                        {c.name?.[0]?.toUpperCase() ?? "?"}
                      </div>
                      <div className="min-w-0 flex-1">
                        <h3 className="text-display font-semibold truncate">
                          {c.name}
                        </h3>
                        <div className="mt-1 flex items-center gap-1 text-xs text-[var(--fg-muted)]">
                          <MapPin className="h-3 w-3" />
                          {c.city}
                          {c.state && `, ${c.state}`}
                        </div>
                        <div className="mt-3 flex gap-2">
                          <Badge variant="success">Active</Badge>
                        </div>
                      </div>
                    </div>
                  </Card>
                </Link>
              ))}
            </div>
          )}
        </div>
      </main>

      <Footer />
      <BottomNav />
    </div>
  );
}