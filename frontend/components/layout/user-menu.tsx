"use client";

import * as React from "react";
import Link from "next/link";
import { LogOut, Ticket, User as UserIcon, LayoutDashboard } from "lucide-react";
import { useAuthStore } from "@/lib/stores/auth-store";
import { useLogout } from "@/lib/hooks/use-auth";
import { initials, cn } from "@/lib/utils";

export function UserMenu() {
  const user = useAuthStore((s) => s.user);
  const [open, setOpen] = React.useState(false);
  const logout = useLogout();
  const ref = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    function onDocClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", onDocClick);
    return () => document.removeEventListener("mousedown", onDocClick);
  }, []);

  if (!user) return null;

  const isAdmin =
    user.role === "COLLEGE_ADMIN" ||
    user.role === "SUPER_ADMIN" ||
    user.role === "ORGANIZER";

  return (
    <div className="relative" ref={ref}>
      <button
        onClick={() => setOpen((v) => !v)}
        className="flex h-10 items-center gap-2 rounded-full border border-[var(--border)] pl-1 pr-3 hover:bg-[var(--bg-muted)] transition-colors"
      >
        <div className="grid h-8 w-8 place-items-center rounded-full bg-brand-gradient text-white text-xs font-bold">
          {initials(user.full_name)}
        </div>
        <span className="hidden sm:block text-sm font-medium truncate max-w-[100px]">
          {user.full_name.split(" ")[0]}
        </span>
      </button>

      {open && (
        <div className="absolute right-0 mt-2 w-60 rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-xl overflow-hidden z-50">
          {/* Header */}
          <div className="p-3 border-b border-[var(--border)]">
            <div className="text-sm font-medium truncate">{user.full_name}</div>
            <div className="text-xs text-[var(--fg-muted)] truncate">{user.email}</div>
            <div className="mt-1.5 inline-flex items-center rounded-full bg-[var(--bg-muted)] px-2 py-0.5 text-[10px] font-medium text-[var(--fg-muted)]">
              {user.role}
            </div>
          </div>

          {/* Items */}
          <div className="p-1">
            <MenuItem href="/tickets" icon={Ticket} onClick={() => setOpen(false)}>
              My Tickets
            </MenuItem>
            <MenuItem href="/profile" icon={UserIcon} onClick={() => setOpen(false)}>
              Profile
            </MenuItem>
            {isAdmin && (
              <MenuItem
                href="/admin"
                icon={LayoutDashboard}
                onClick={() => setOpen(false)}
              >
                Admin Dashboard
              </MenuItem>
            )}
          </div>

          {/* Logout */}
          <div className="p-1 border-t border-[var(--border)]">
            <button
              onClick={() => {
                setOpen(false);
                logout.mutate();
              }}
              disabled={logout.isPending}
              className="w-full flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-[var(--color-danger)] hover:bg-[var(--color-danger)]/10 transition-colors disabled:opacity-50"
            >
              <LogOut className="h-4 w-4" />
              {logout.isPending ? "Logging out…" : "Log out"}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

function MenuItem({
  href,
  icon: Icon,
  children,
  onClick,
}: {
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  children: React.ReactNode;
  onClick?: () => void;
}) {
  return (
    <Link
      href={href}
      onClick={onClick}
      className={cn(
        "flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors",
        "hover:bg-[var(--bg-muted)]"
      )}
    >
      <Icon className="h-4 w-4 text-[var(--fg-muted)]" />
      {children}
    </Link>
  );
}