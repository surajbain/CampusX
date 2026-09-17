"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Building2,
  CalendarDays,
  Users,
  Settings,
  GraduationCap,
  ChevronLeft,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/lib/stores/auth-store";

const items = [
  { href: "/admin", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/colleges", label: "Colleges", icon: Building2 },
  { href: "/admin/events", label: "Events", icon: CalendarDays },
  { href: "/admin/users", label: "Users", icon: Users },
  { href: "/admin/settings", label: "Settings", icon: Settings },
];

export function AdminSidebar() {
  const pathname = usePathname();
  const user = useAuthStore((s) => s.user);
  const isSuperAdmin = user?.role === "SUPER_ADMIN";

  // Only super admin can see "Colleges" section
  const visibleItems = items.filter((item) => {
    if (item.href === "/admin/colleges") return isSuperAdmin;
    return true;
  });

  return (
    <aside className="hidden lg:flex lg:flex-col w-64 shrink-0 border-r border-[var(--border)] bg-[var(--bg-subtle)] min-h-[calc(100vh-4rem)]">
      <div className="p-4 border-b border-[var(--border)]">
        <Link
          href="/"
          className="flex items-center gap-2 text-sm text-[var(--fg-muted)] hover:text-[var(--fg)] transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
          Back to site
        </Link>
      </div>

      <nav className="flex-1 p-3 space-y-1">
        {visibleItems.map(({ href, label, icon: Icon }) => {
          const active =
            href === "/admin" ? pathname === "/admin" : pathname.startsWith(href);
          return (
            <Link
              key={href}
              href={href}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors",
                active
                  ? "bg-[var(--color-brand)]/10 text-[var(--color-brand)]"
                  : "text-[var(--fg-muted)] hover:bg-[var(--bg-muted)] hover:text-[var(--fg)]"
              )}
            >
              <Icon className="h-4 w-4" />
              {label}
            </Link>
          );
        })}
      </nav>

      {user && (
        <div className="p-4 border-t border-[var(--border)]">
          <div className="text-xs text-[var(--fg-subtle)]">Signed in as</div>
          <div className="text-sm font-medium truncate">{user.full_name}</div>
          <div className="text-xs text-[var(--fg-muted)]">{user.role}</div>
        </div>
      )}
    </aside>
  );
}

// Mobile bottom nav for admin
export function AdminBottomNav() {
  const pathname = usePathname();
  const user = useAuthStore((s) => s.user);
  const isSuperAdmin = user?.role === "SUPER_ADMIN";

  const visibleItems = items.slice(0, 4).filter((item) => {
    if (item.href === "/admin/colleges") return isSuperAdmin;
    return true;
  });

  return (
    <nav className="lg:hidden fixed bottom-0 left-0 right-0 z-50 border-t border-[var(--border)] glass pb-[env(safe-area-inset-bottom)]">
      <div
        className={cn(
          "grid h-16",
          visibleItems.length === 4
            ? "grid-cols-4"
            : visibleItems.length === 3
            ? "grid-cols-3"
            : "grid-cols-2"
        )}
      >
        {visibleItems.map(({ href, label, icon: Icon }) => {
          const active =
            href === "/admin" ? pathname === "/admin" : pathname.startsWith(href);
          return (
            <Link
              key={href}
              href={href}
              className={cn(
                "flex flex-col items-center justify-center gap-1 text-xs transition-colors",
                active
                  ? "text-[var(--color-brand)]"
                  : "text-[var(--fg-muted)]"
              )}
            >
              <Icon className="h-5 w-5" />
              <span className="font-medium">{label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}