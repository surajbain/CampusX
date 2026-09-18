"use client";

import * as React from "react";
import { User as UserIcon, Mail, Phone, Building2, Calendar } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/lib/stores/auth-store";
import { initials } from "@/lib/utils";

export default function ProfilePage() {
  const user = useAuthStore((s) => s.user);

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <div className="text-sm text-[var(--fg-muted)]">Loading profile…</div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl">
      <div className="mb-8">
        <h1 className="text-display text-3xl font-bold tracking-tight">
          Profile
        </h1>
        <p className="mt-1.5 text-sm text-[var(--fg-muted)]">
          Your account information
        </p>
      </div>

      {/* Avatar + name */}
      <Card className="p-6 mb-6">
        <div className="flex items-center gap-5">
          <div className="grid h-20 w-20 place-items-center rounded-full bg-brand-gradient text-white text-2xl font-bold">
            {initials(user.full_name)}
          </div>
          <div className="flex-1">
            <h2 className="text-display text-2xl font-bold tracking-tight">
              {user.full_name}
            </h2>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <Badge variant="brand">{user.role}</Badge>
              {user.email_verified && (
                <Badge variant="success">Email verified</Badge>
              )}
            </div>
          </div>
          <Button variant="secondary" disabled>
            Edit profile (coming soon)
          </Button>
        </div>
      </Card>

      {/* Details */}
      <Card className="p-6">
        <h3 className="text-display text-lg font-semibold mb-5">
          Account details
        </h3>

        <div className="space-y-4">
          <DetailRow
            icon={UserIcon}
            label="Full name"
            value={user.full_name}
          />
          <DetailRow
            icon={Mail}
            label="Email"
            value={user.email}
          />
          {user.phone && (
            <DetailRow
              icon={Phone}
              label="Phone"
              value={user.phone}
            />
          )}
          {user.college_id && (
            <DetailRow
              icon={Building2}
              label="College ID"
              value={user.college_id}
              mono
            />
          )}
          <DetailRow
            icon={Calendar}
            label="Member since"
            value={new Date(user.created_at).toLocaleDateString("en-IN", {
              day: "numeric",
              month: "long",
              year: "numeric",
            })}
          />
        </div>
      </Card>
    </div>
  );
}

function DetailRow({
  icon: Icon,
  label,
  value,
  mono,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
  mono?: boolean;
}) {
  return (
    <div className="flex items-start gap-3 py-3 border-b border-[var(--border)] last:border-0">
      <div className="grid h-9 w-9 shrink-0 place-items-center rounded-lg bg-[var(--bg-muted)]">
        <Icon className="h-4 w-4 text-[var(--fg-muted)]" />
      </div>
      <div className="flex-1 min-w-0">
        <div className="text-xs text-[var(--fg-subtle)]">{label}</div>
        <div
          className={`mt-0.5 text-sm font-medium truncate ${
            mono ? "font-mono" : ""
          }`}
        >
          {value}
        </div>
      </div>
    </div>
  );
}