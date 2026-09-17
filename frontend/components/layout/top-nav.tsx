"use client";

import * as React from "react";
import Link from "next/link";
import { Search, User, GraduationCap, Ticket, Compass } from "lucide-react";
import { ThemeToggle } from "./theme-toggle";
import { UserMenu } from "./user-menu";
import { useAuthStore } from "@/lib/stores/auth-store";

export function TopNav() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return (
    <header className="sticky top-0 z-50 w-full border-b border-[var(--border)] glass">
      <div className="cx-container flex h-16 items-center gap-4">
        <Link href="/" className="flex items-center gap-2 group shrink-0">
          <div className="grid h-9 w-9 place-items-center rounded-xl bg-brand-gradient shadow-md transition-transform group-hover:scale-105">
            <GraduationCap className="h-5 w-5 text-white" />
          </div>
          <span className="text-display text-xl font-bold tracking-tight hidden sm:block">
            Campus<span className="text-brand-gradient">X</span>
          </span>
        </Link>

        {/* Nav links (desktop) */}
        <nav className="hidden md:flex items-center gap-1 ml-4">
          <Link
            href="/events"
            className="flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm font-medium text-[var(--fg-muted)] hover:text-[var(--fg)] hover:bg-[var(--bg-muted)] transition-colors"
          >
            <Compass className="h-4 w-4" />
            Events
          </Link>
          <Link
            href="/colleges"
            className="rounded-lg px-3 py-2 text-sm font-medium text-[var(--fg-muted)] hover:text-[var(--fg)] hover:bg-[var(--bg-muted)] transition-colors"
          >
            Colleges
          </Link>
        </nav>

        {/* Search (visual) */}
        <div className="flex-1 max-w-md mx-auto hidden md:block">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[var(--fg-subtle)]" />
            <input
              type="text"
              placeholder="Search events, colleges, cities…"
              className="w-full h-10 rounded-full border border-[var(--border)] bg-[var(--bg-subtle)] pl-10 pr-4 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:bg-[var(--bg)] transition-colors"
            />
          </div>
        </div>

        <div className="ml-auto flex items-center gap-1">
          <ThemeToggle />

          {isAuthenticated ? (
            <>
              <Link
                href="/tickets"
                className="hidden sm:flex h-10 w-10 items-center justify-center rounded-full hover:bg-[var(--bg-muted)] transition-colors"
                aria-label="My tickets"
              >
                <Ticket className="h-5 w-5" />
              </Link>
              <UserMenu />
            </>
          ) : (
            <Link
              href="/login"
              className="flex h-10 items-center gap-2 rounded-full border border-[var(--border)] px-4 text-sm font-medium hover:bg-[var(--bg-muted)] transition-colors"
            >
              <User className="h-4 w-4" />
              <span className="hidden sm:inline">Sign in</span>
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}