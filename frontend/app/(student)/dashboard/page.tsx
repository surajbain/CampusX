"use client";

import Link from "next/link";
import { Ticket, CheckCircle2, Clock, Compass, ArrowRight } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { Badge } from "@/components/ui/badge";
import { useAuthStore } from "@/lib/stores/auth-store";
import { useMyRegistrations } from "@/lib/hooks/use-registrations";
import { RegistrationCard } from "@/components/student/registration-card";

export default function StudentDashboard() {
  const user = useAuthStore((s) => s.user);
  const { data, isLoading } = useMyRegistrations({ limit: 100 });

  const registrations = ((data as any)?.data ?? []) as any[];

  const confirmed = registrations.filter((r) => r.status === "CONFIRMED").length;
  const pending = registrations.filter((r) => r.status === "PENDING").length;
  const upcoming = registrations
    .filter((r) => {
      if (!r.event_starts_at) return false;
      return new Date(r.event_starts_at) > new Date();
    })
    .slice(0, 3);

  return (
    <div>
      {/* Welcome */}
      <div className="mb-8">
        <h1 className="text-display text-3xl font-bold tracking-tight">
          Hey {user?.full_name?.split(" ")[0] || "there"} 👋
        </h1>
        <p className="mt-1.5 text-sm text-[var(--fg-muted)]">
          Here&apos;s what&apos;s happening with your events
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
        <StatCard
          label="Total registrations"
          value={registrations.length}
          icon={Ticket}
          loading={isLoading}
          color="brand"
        />
        <StatCard
          label="Confirmed"
          value={confirmed}
          icon={CheckCircle2}
          loading={isLoading}
          color="success"
        />
        <StatCard
          label="Pending payment"
          value={pending}
          icon={Clock}
          loading={isLoading}
          color="warning"
        />
      </div>

      {/* Upcoming */}
      <div className="mb-8">
        <div className="flex items-end justify-between mb-4">
          <div>
            <h2 className="text-display text-xl font-semibold tracking-tight">
              Your upcoming events
            </h2>
            <p className="mt-1 text-xs text-[var(--fg-muted)]">
              Events you&apos;re registered for
            </p>
          </div>
          <Link href="/tickets">
            <Button variant="ghost" size="sm">
              View all
              <ArrowRight className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>

        {isLoading && (
          <div className="space-y-3">
            {[1, 2].map((i) => (
              <Skeleton key={i} className="h-32 rounded-2xl" />
            ))}
          </div>
        )}

        {!isLoading && upcoming.length === 0 && (
          <EmptyState
            icon={Compass}
            title="No upcoming events"
            description="Discover events across colleges and register to attend."
            action={
              <Link href="/events">
                <Button>
                  <Compass className="h-4 w-4" />
                  Browse events
                </Button>
              </Link>
            }
          />
        )}

        {!isLoading && upcoming.length > 0 && (
          <div className="space-y-3">
            {upcoming.map((reg) => (
              <RegistrationCard key={reg.id} registration={reg} />
            ))}
          </div>
        )}
      </div>

      {/* CTA */}
      <Card className="p-6 bg-gradient-to-br from-[var(--color-brand-from)]/10 to-[var(--color-brand-to)]/5 border-[var(--color-brand)]/20">
        <div className="flex items-start gap-4">
          <div className="grid h-12 w-12 shrink-0 place-items-center rounded-xl bg-brand-gradient text-white">
            <Compass className="h-6 w-6" />
          </div>
          <div className="flex-1">
            <h3 className="text-display font-semibold">
              Discover more events
            </h3>
            <p className="mt-1 text-sm text-[var(--fg-muted)]">
              Hackathons, cultural fests, sports meets — all at your fingertips.
            </p>
            <Link href="/events" className="inline-block mt-3">
              <Button size="sm">
                Browse events
                <ArrowRight className="h-3.5 w-3.5" />
              </Button>
            </Link>
          </div>
        </div>
      </Card>
    </div>
  );
}

function StatCard({
  label,
  value,
  icon: Icon,
  loading,
  color,
}: {
  label: string;
  value: number;
  icon: React.ComponentType<{ className?: string }>;
  loading?: boolean;
  color: "brand" | "success" | "warning";
}) {
  const colors = {
    brand: "bg-[var(--color-brand)]/10 text-[var(--color-brand)]",
    success: "bg-[var(--color-success)]/10 text-[var(--color-success)]",
    warning: "bg-[var(--color-warning)]/10 text-[var(--color-warning)]",
  };

  return (
    <Card className="p-5">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-xs text-[var(--fg-subtle)]">{label}</div>
          {loading ? (
            <Skeleton className="h-8 w-12 mt-2" />
          ) : (
            <div className="text-display text-3xl font-bold mt-1">{value}</div>
          )}
        </div>
        <div className={`grid h-11 w-11 place-items-center rounded-xl ${colors[color]}`}>
          <Icon className="h-5 w-5" />
        </div>
      </div>
    </Card>
  );
}

let _: unknown = undefined;