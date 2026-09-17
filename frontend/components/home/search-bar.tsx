"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Search } from "lucide-react";

export function HomeSearchBar() {
  const router = useRouter();
  const [q, setQ] = React.useState("");

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    const params = q.trim() ? `?q=${encodeURIComponent(q.trim())}` : "";
    router.push(`/events${params}`);
  };

  return (
    <form onSubmit={submit} className="mt-8 max-w-xl mx-auto">
      <div className="relative">
        <Search className="absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-[var(--fg-subtle)]" />
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Search hackathons, fests, workshops…"
          className="w-full h-14 rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)]/80 backdrop-blur pl-12 pr-32 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-4 focus:ring-[var(--color-brand)]/15 shadow-lg transition-all"
        />
        <button
          type="submit"
          className="absolute right-2 top-1/2 -translate-y-1/2 h-10 rounded-xl bg-brand-gradient px-5 text-sm font-medium text-white shadow-md hover:brightness-110 transition-all"
        >
          Search
        </button>
      </div>
    </form>
  );
}