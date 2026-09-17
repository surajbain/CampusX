"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Loader2, ShieldAlert } from "lucide-react";
import { useAuthStore } from "@/lib/stores/auth-store";
import { EmptyState } from "@/components/ui/empty-state";
import { Button } from "@/components/ui/button";

const ALLOWED_ROLES = ["COLLEGE_ADMIN", "SUPER_ADMIN", "ORGANIZER"];

export function AdminGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const [hydrated, setHydrated] = React.useState(false);

  React.useEffect(() => {
    setHydrated(true);
  }, []);

  React.useEffect(() => {
    if (hydrated && !isAuthenticated) {
      router.replace("/login?redirect=/admin");
    }
  }, [hydrated, isAuthenticated, router]);

  if (!hydrated) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-[var(--fg-muted)]" />
      </div>
    );
  }

  if (!isAuthenticated) return null;

  const isAllowed = user && ALLOWED_ROLES.includes(user.role);

  if (!isAllowed) {
    return (
      <EmptyState
        icon={ShieldAlert}
        title="Access denied"
        description="You don't have permission to access the admin dashboard. Only college admins, organizers, and super admins can access this area."
        action={
          <Button onClick={() => router.push("/")}>Go home</Button>
        }
      />
    );
  }

  return <>{children}</>;
}