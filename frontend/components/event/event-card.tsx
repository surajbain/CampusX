import Link from "next/link";
import { Calendar, MapPin, Trophy, Users } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { formatPrice, formatDateRange, truncate } from "@/lib/utils";
import { categoryMeta } from "@/lib/constants";
import type { Event } from "@/lib/api/types";

const categoryVariantMap: Record<string, "hackathon" | "cultural" | "sports" | "workshop" | "techfest" | "other"> = {
  HACKATHON: "hackathon",
  CULTURAL: "cultural",
  SPORTS: "sports",
  WORKSHOP: "workshop",
  TECH_FEST: "techfest",
  OTHER: "other",
};

export function EventCard({ event }: { event: Event }) {
  const cat = categoryMeta(event.category);
  const price = formatPrice(event.price_paise, event.currency);
  const dates = formatDateRange(event.starts_at, event.ends_at);

  return (
    <Link href={`/events/${event.id}`} className="group block">
      <Card interactive className="h-full flex flex-col">
        {/* Poster */}
        <div className="relative aspect-[16/9] overflow-hidden bg-gradient-to-br from-[var(--color-brand-from)]/20 to-[var(--color-brand-to)]/20">
          {event.poster_url ? (
            <img
              src={event.poster_url}
              alt={event.title}
              className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
          ) : (
            <div
              className="flex h-full w-full items-center justify-center text-5xl"
              style={{ background: `linear-gradient(135deg, ${cat.color}30, ${cat.color}10)` }}
            >
              {cat.emoji}
            </div>
          )}
          {/* Category badge top-left */}
          <div className="absolute top-3 left-3">
            <Badge variant={categoryVariantMap[event.category] || "other"}>
              {cat.emoji} {cat.label}
            </Badge>
          </div>
          {/* Price badge top-right */}
          <div className="absolute top-3 right-3">
            <Badge
              variant={event.price_paise === 0 ? "success" : "brand"}
            >
              {price}
            </Badge>
          </div>
          {/* Featured badge bottom-left */}
          {event.is_featured && (
            <div className="absolute bottom-3 left-3">
              <Badge variant="warning">🔥 Featured</Badge>
            </div>
          )}
        </div>

        {/* Body */}
        <div className="flex flex-1 flex-col p-5">
          <h3 className="text-display text-base font-semibold leading-snug line-clamp-2 group-hover:text-[var(--color-brand)] transition-colors">
            {event.title}
          </h3>

          {event.description && (
            <p className="mt-2 text-sm text-[var(--fg-muted)] line-clamp-2">
              {truncate(event.description, 100)}
            </p>
          )}

          {/* College + City */}
          <div className="mt-3 flex items-center gap-1.5 text-xs text-[var(--fg-muted)]">
            <span className="font-medium text-[var(--fg)]">{event.college_name}</span>
            {event.college_city && (
              <>
                <span className="text-[var(--fg-subtle)]">•</span>
                <span>{event.college_city}</span>
              </>
            )}
          </div>

          {/* Meta */}
          <div className="mt-auto pt-4 flex flex-wrap gap-x-3 gap-y-1.5 text-xs text-[var(--fg-muted)]">
            {dates && (
              <span className="inline-flex items-center gap-1">
                <Calendar className="h-3.5 w-3.5" />
                {dates}
              </span>
            )}
            {event.venue && (
              <span className="inline-flex items-center gap-1">
                <MapPin className="h-3.5 w-3.5" />
                {truncate(event.venue, 20)}
              </span>
            )}
            {event.prize_pool_paise > 0 && (
              <span className="inline-flex items-center gap-1 text-[var(--color-warning)] font-medium">
                <Trophy className="h-3.5 w-3.5" />
                {formatPrice(event.prize_pool_paise, event.currency)}
              </span>
            )}
            {event.allow_teams && event.team_size_max && (
              <span className="inline-flex items-center gap-1">
                <Users className="h-3.5 w-3.5" />
                {event.team_size_min}–{event.team_size_max}
              </span>
            )}
          </div>
        </div>
      </Card>
    </Link>
  );
}