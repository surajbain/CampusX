"use client";

import Link from "next/link";
import { Building2, CalendarDays, Users, ArrowRight } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { useAdminStats } from "@/lib/hooks/use-admin";
import { useAuthStore } from "@/lib/stores/auth-store";

export default function AdminDashboard() {
  const { data: stats, isLoading } = useAdminStats();
  const user = useAuthStore((s) => s.user);

  const isSuperAdmin = user?.role === "SUPER_ADMIN";

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-display text-3xl font-bold tracking-tight">
          Welcome back, {user?.full_name?.split(" ")[0]} 👋
        </h1>
        <p className="mt-1.5 text-sm text-[var(--fg-muted)]">
          Here&apos;s what&apos;s happening across CampusX
        </p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-5 mb-8">
        <StatCard
          label="Active colleges"
          value={stats?.colleges ?? 0}
          icon={Building2}
          loading={isLoading}
        />
        <StatCard
          label="Published events"
          value={stats?.published_events ?? 0}
          icon={CalendarDays}
          loading={isLoading}
        />
        <StatCard
          label="Total users"
          value={0} // Will fill when API exists
          icon={Users}
          loading={isLoading}
        />
      </div>

      {/* Quick actions */}
      <div>
        <h2 className="text-display text-lg font-semibold mb-4">
          Quick actions
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {isSuperAdmin && (
            <QuickActionCard
              href="/admin/colleges"
              icon={Building2}
              title="Manage colleges"
              description="Create, edit, or delete colleges on the platform"
            />
          )}
          <QuickActionCard
            href="/admin/events"
            icon={CalendarDays}
            title="Manage events"
            description="Create and manage events for your college"
          />
        </div>
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
  icon: Icon,
  loading,
}: {
  label: string;
  value: number;
  icon: React.ComponentType<{ className?: string }>;
  loading?: boolean;
}) {
  return (
    <Card className="p-5">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-xs text-[var(--fg-subtle)]">{label}</div>
          {loading ? (
            <Skeleton className="h-8 w-16 mt-2" />
          ) : (
            <div className="text-display text-3xl font-bold mt-1">{value}</div>
          )}
        </div>
        <div className="grid h-11 w-11 place-items-center rounded-xl bg-[var(--color-brand)]/10">
          <Icon className="h-5 w-5 text-[var(--color-brand)]" />
        </div>
      </div>
    </Card>
  );
}

function QuickActionCard({
  href,
  icon: Icon,
  title,
  description,
}: {
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  title: string;
  description: string;
}) {
  return (
    <Link href={href}>
      <Card interactive className="p-5 h-full">
        <div className="flex items-start gap-4">
          <div className="grid h-11 w-11 shrink-0 place-items-center rounded-xl bg-[var(--color-brand)]/10">
            <Icon className="h-5 w-5 text-[var(--color-brand)]" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <h3 className="text-display font-semibold">{title}</h3>
              <ArrowRight className="h-3.5 w-3.5 text-[var(--fg-subtle)]" />
            </div>
            <p className="mt-1 text-sm text-[var(--fg-muted)]">{description}</p>
          </div>
        </div>
      </Card>
    </Link>
  );
}