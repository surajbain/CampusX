"use client";

import * as React from "react";
import Link from "next/link";
import {
  Calendar,
  MapPin,
  Trophy,
  Users,
  ArrowLeft,
  Share2,
  MessageCircle,
  Check,
  ExternalLink,
  Mail,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { formatPrice, formatDateRange, formatRelative, cn } from "@/lib/utils";
import { categoryMeta } from "@/lib/constants";
import type { Event } from "@/lib/api/types";
import { RegisterButton } from "./register-button";

const categoryVariantMap: Record<
  string,
  "hackathon" | "cultural" | "sports" | "workshop" | "techfest" | "other"
> = {
  HACKATHON: "hackathon",
  CULTURAL: "cultural",
  SPORTS: "sports",
  WORKSHOP: "workshop",
  TECH_FEST: "techfest",
  OTHER: "other",
};

export function EventDetailView({ event }: { event: Event }) {
  const cat = categoryMeta(event.category);
  const price = formatPrice(event.price_paise, event.currency);
  const dates = formatDateRange(event.starts_at, event.ends_at);
  const [copied, setCopied] = React.useState(false);
  const [tab, setTab] = React.useState<"about" | "schedule" | "rules" | "faq">(
    "about"
  );

  const shareUrl = typeof window !== "undefined" ? window.location.href : "";

  const copyLink = async () => {
    try {
      await navigator.clipboard.writeText(shareUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {}
  };

  const whatsappShare = () => {
    const msg = encodeURIComponent(`${event.title} — ${dates}\n${shareUrl}`);
    window.open(`https://wa.me/?text=${msg}`, "_blank");
  };

  return (
    <div className="pb-24 md:pb-0">
      {/* Hero */}
      <div className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10 bg-gradient-to-br from-[var(--color-brand-from)]/15 via-[var(--color-brand-to)]/10 to-transparent">
          <div
            className="absolute top-0 left-1/2 -translate-x-1/2 h-[400px] w-[700px] rounded-full blur-[120px] opacity-50"
            style={{ backgroundColor: `${cat.color}33` }}
          />
        </div>

        <div className="cx-container pt-6 pb-12">
          <Link
            href="/events"
            className="inline-flex items-center gap-2 text-sm text-[var(--fg-muted)] hover:text-[var(--fg)] mb-6 transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to events
          </Link>

          <div className="grid grid-cols-1 lg:grid-cols-[1fr_380px] gap-8">
            {/* Left */}
            <div>
              <div className="flex flex-wrap items-center gap-2 mb-4">
                <Badge variant={categoryVariantMap[event.category] || "other"}>
                  {cat.emoji} {cat.label}
                </Badge>
                {event.is_featured && (
                  <Badge variant="warning">🔥 Featured</Badge>
                )}
                <Badge
                  variant={event.status === "PUBLISHED" ? "success" : "default"}
                >
                  {event.status}
                </Badge>
              </div>

              <h1 className="text-display text-3xl md:text-5xl font-bold tracking-tight leading-tight">
                {event.title}
              </h1>

              <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-[var(--fg-muted)]">
                <Link
                  href={`/colleges/${event.college_slug}`}
                  className="font-medium text-[var(--fg)] hover:text-[var(--color-brand)] transition-colors"
                >
                  {event.college_name}
                </Link>
                {event.college_city && (
                  <>
                    <span>·</span>
                    <span className="inline-flex items-center gap-1">
                      <MapPin className="h-3.5 w-3.5" />
                      {event.college_city}
                    </span>
                  </>
                )}
              </div>

              {event.description && (
                <p className="mt-6 text-base text-[var(--fg-muted)] leading-relaxed max-w-3xl">
                  {event.description}
                </p>
              )}
            </div>

            {/* Right — Sidebar */}
            <div className="lg:sticky lg:top-20 lg:self-start">
              <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] overflow-hidden shadow-lg">
                <div className="relative aspect-[16/10] overflow-hidden bg-gradient-to-br from-[var(--color-brand-from)]/20 to-[var(--color-brand-to)]/20">
                  {event.poster_url ? (
                    <img
                      src={event.poster_url}
                      alt={event.title}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div
                      className="flex h-full w-full items-center justify-center text-7xl"
                      style={{
                        background: `linear-gradient(135deg, ${cat.color}30, ${cat.color}10)`,
                      }}
                    >
                      {cat.emoji}
                    </div>
                  )}
                </div>

                <div className="p-5 space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="text-xs text-[var(--fg-muted)]">Entry</div>
                      <div
                        className={cn(
                          "text-display text-2xl font-bold",
                          event.price_paise === 0 && "text-[var(--color-success)]"
                        )}
                      >
                        {price}
                      </div>
                    </div>
                    {event.capacity && (
                      <div className="text-right">
                        <div className="text-xs text-[var(--fg-muted)]">
                          Capacity
                        </div>
                        <div className="text-display text-lg font-semibold">
                          {event.capacity}
                        </div>
                      </div>
                    )}
                  </div>

                 <RegisterButton
                   eventId={event.id}
                   pricePaise={event.price_paise}
                   allowTeams={event.allow_teams}
                   size="lg"
                   className="w-full"
                  />

                  {event.whatsapp_link && (
                    <a
                      href={event.whatsapp_link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center justify-center gap-2 h-11 rounded-xl border border-[var(--border)] text-sm font-medium hover:bg-[var(--bg-muted)] transition-colors"
                    >
                      <MessageCircle className="h-4 w-4" />
                      Join WhatsApp group
                    </a>
                  )}

                  <div className="flex gap-2 pt-2 border-t border-[var(--border)]">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={whatsappShare}
                      className="flex-1"
                    >
                      <MessageCircle className="h-4 w-4" />
                      WhatsApp
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={copyLink}
                      className="flex-1"
                    >
                      {copied ? (
                        <>
                          <Check className="h-4 w-4" />
                          Copied
                        </>
                      ) : (
                        <>
                          <Share2 className="h-4 w-4" />
                          Copy link
                        </>
                      )}
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Info strip */}
      <div className="border-y border-[var(--border)] bg-[var(--bg-subtle)]">
        <div className="cx-container py-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <InfoBlock icon={Calendar} label="Dates" value={dates || "TBA"} />
            <InfoBlock icon={MapPin} label="Venue" value={event.venue || "TBA"} />
            <InfoBlock
              icon={Trophy}
              label="Prize pool"
              value={
                event.prize_pool_paise > 0
                  ? formatPrice(event.prize_pool_paise)
                  : "—"
              }
            />
            <InfoBlock
              icon={Users}
              label={event.allow_teams ? "Team size" : "Mode"}
              value={
                event.allow_teams && event.team_size_max
                  ? `${event.team_size_min}–${event.team_size_max} members`
                  : "Individual"
              }
            />
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="cx-container py-12">
        <div className="flex gap-1 border-b border-[var(--border)] mb-8 overflow-x-auto no-scrollbar">
          {(["about", "schedule", "rules", "faq"] as const).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={cn(
                "px-4 py-3 text-sm font-medium whitespace-nowrap border-b-2 transition-colors -mb-px",
                tab === t
                  ? "border-[var(--color-brand)] text-[var(--fg)]"
                  : "border-transparent text-[var(--fg-muted)] hover:text-[var(--fg)]"
              )}
            >
              {t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>

        <div className="max-w-3xl">
          {tab === "about" && (
            <div>
              {event.description ? (
                <p className="text-[var(--fg-muted)] leading-relaxed whitespace-pre-wrap">
                  {event.description}
                </p>
              ) : (
                <p className="text-[var(--fg-subtle)]">
                  No description provided.
                </p>
              )}

              <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 gap-4">
                {event.registration_opens && (
                  <DetailRow
                    label="Registration opens"
                    value={new Date(event.registration_opens).toLocaleDateString(
                      "en-IN",
                      { day: "numeric", month: "long", year: "numeric" }
                    )}
                    hint={formatRelative(event.registration_opens)}
                  />
                )}
                {event.registration_closes && (
                  <DetailRow
                    label="Registration closes"
                    value={new Date(
                      event.registration_closes
                    ).toLocaleDateString("en-IN", {
                      day: "numeric",
                      month: "long",
                      year: "numeric",
                    })}
                    hint={formatRelative(event.registration_closes)}
                  />
                )}
                {event.contact_email && (
                  <DetailRow
                    label="Contact"
                    value={event.contact_email}
                    icon={Mail}
                  />
                )}
                {event.college_slug && (
                  <DetailRow
                    label="College"
                    value={event.college_name || ""}
                    icon={ExternalLink}
                    href={`/colleges/${event.college_slug}`}
                  />
                )}
              </div>
            </div>
          )}

          {tab === "schedule" && (
            <p className="text-[var(--fg-muted)]">
              Detailed schedule will be announced soon.
            </p>
          )}

          {tab === "rules" && (
            <div className="text-[var(--fg-muted)] whitespace-pre-wrap leading-relaxed">
              {event.rules || "Rules will be published soon."}
            </div>
          )}

          {tab === "faq" && (
            <p className="text-[var(--fg-muted)]">FAQ will be published soon.</p>
          )}
        </div>
      </div>
    </div>
  );
}

function InfoBlock({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-start gap-3">
      <div className="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-[var(--bg-elevated)] border border-[var(--border)]">
        <Icon className="h-4 w-4 text-[var(--fg-muted)]" />
      </div>
      <div className="min-w-0">
        <div className="text-xs text-[var(--fg-subtle)]">{label}</div>
        <div className="text-sm font-medium truncate">{value}</div>
      </div>
    </div>
  );
}

function DetailRow({
  label,
  value,
  hint,
  icon: Icon,
  href,
}: {
  label: string;
  value: string;
  hint?: string;
  icon?: React.ComponentType<{ className?: string }>;
  href?: string;
}) {
  const content = (
    <div className="rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)] p-4">
      <div className="text-xs text-[var(--fg-subtle)] mb-1">{label}</div>
      <div className="text-sm font-medium flex items-center gap-1.5">
        {Icon && <Icon className="h-3.5 w-3.5 text-[var(--fg-muted)]" />}
        {value}
      </div>
      {hint && <div className="text-xs text-[var(--fg-muted)] mt-1">{hint}</div>}
    </div>
  );

  if (href) {
    return (
      <Link href={href} className="block hover:opacity-90 transition-opacity">
        {content}
      </Link>
    );
  }
  return content;
}