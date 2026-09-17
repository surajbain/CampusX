"use client";

import { motion } from "framer-motion";
import { HomeSearchBar } from "./search-bar";

export function Hero() {
  return (
    <section className="relative overflow-hidden">
      {/* Gradient background */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute inset-0 bg-gradient-to-b from-[var(--color-brand-from)]/10 via-[var(--color-brand-to)]/5 to-transparent" />
        <div className="absolute top-0 left-1/2 -translate-x-1/2 h-[500px] w-[800px] rounded-full bg-[var(--color-brand-from)]/20 blur-[120px] opacity-60" />
        <div className="absolute top-20 right-1/4 h-[300px] w-[500px] rounded-full bg-[var(--color-brand-to)]/20 blur-[100px] opacity-50" />
      </div>

      <div className="cx-container py-20 md:py-28">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-3xl mx-auto text-center"
        >
          {/* Live badge */}
          <div className="inline-flex items-center gap-2 rounded-full border border-[var(--border)] bg-[var(--bg-elevated)]/60 backdrop-blur px-4 py-1.5 text-xs font-medium text-[var(--fg-muted)] mb-6">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-[var(--color-success)] opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-[var(--color-success)]" />
            </span>
            Live across 2 colleges
          </div>

          {/* Headline */}
          <h1 className="text-display text-4xl md:text-6xl lg:text-7xl font-bold tracking-tight leading-[1.05]">
            Discover what&apos;s happening{" "}
            <span className="text-brand-gradient">around you</span>
          </h1>

          {/* Subtitle */}
          <p className="mt-6 text-base md:text-lg text-[var(--fg-muted)] max-w-2xl mx-auto">
            Hackathons, cultural fests, sports meets, and workshops — all from
            colleges near you. Free to discover. Free to register.
          </p>

          {/* Search bar */}
          <HomeSearchBar />

          {/* Trust signals */}
          <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
            <div className="flex items-center gap-2 text-xs text-[var(--fg-muted)]">
              <span className="text-base">🎓</span> 2 colleges
            </div>
            <span className="text-[var(--fg-subtle)]">·</span>
            <div className="flex items-center gap-2 text-xs text-[var(--fg-muted)]">
              <span className="text-base">🎉</span> 4 live events
            </div>
            <span className="text-[var(--fg-subtle)]">·</span>
            <div className="flex items-center gap-2 text-xs text-[var(--fg-muted)]">
              <span className="text-base">💸</span> ₹0 platform fee
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}