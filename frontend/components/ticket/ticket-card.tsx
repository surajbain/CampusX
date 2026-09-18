"use client";

import * as React from "react";
import Link from "next/link";
import { QRCodeSVG } from "qrcode.react";
import {
  Calendar,
  MapPin,
  Ticket as TicketIcon,
  Clock,
  CheckCircle2,
  AlertCircle,
  Ban,
  Download,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { formatDateRange, cn } from "@/lib/utils";
import { categoryMeta } from "@/lib/constants";
import type { Ticket } from "@/lib/api/tickets";

interface Props {
  ticket: Ticket;
  showQR?: boolean;
}

export function TicketCard({ ticket, showQR = false }: Props) {
  const cat = categoryMeta(ticket.event_category);
  const dates = formatDateRange(ticket.event_starts_at, ticket.event_ends_at);

  const statusConfig: Record<
    string,
    { label: string; variant: "success" | "warning" | "danger" | "default"; icon: any }
  > = {
    ACTIVE: { label: "Active", variant: "success", icon: CheckCircle2 },
    USED: { label: "Used", variant: "default", icon: CheckCircle2 },
    REVOKED: { label: "Revoked", variant: "danger", icon: Ban },
    EXPIRED: { label: "Expired", variant: "warning", icon: AlertCircle },
  };

  const status = statusConfig[ticket.status] ?? statusConfig.ACTIVE;
  const StatusIcon = status.icon;

  const now = Date.now();
  const validFrom = ticket.valid_from ? new Date(ticket.valid_from).getTime() : 0;
  const validUntil = ticket.valid_until
    ? new Date(ticket.valid_until).getTime()
    : Infinity;
  const isActive = now >= validFrom && now <= validUntil && ticket.status === "ACTIVE";
  const notYetValid = now < validFrom;

  const qrData = `${ticket.qr_payload}.${ticket.qr_signature}`;

  const downloadQR = () => {
    const svg = document.getElementById(`qr-${ticket.id}`);
    if (!svg) return;
    const svgData = new XMLSerializer().serializeToString(svg);
    const blob = new Blob([svgData], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `ticket-${ticket.ticket_code}.svg`;
    a.click();
    URL.revokeObjectURL(url);
  };

  if (showQR) {
    return (
      <div className="mx-auto max-w-md">
        <div className="relative rounded-3xl border border-[var(--border)] bg-[var(--bg-elevated)] overflow-hidden shadow-2xl">
          {/* Top — gradient banner */}
          <div
            className="relative px-6 pt-6 pb-8 text-white"
            style={{
              background: `linear-gradient(135deg, ${cat.color} 0%, ${cat.color}99 100%)`,
            }}
          >
            <div className="flex items-center gap-2 mb-3">
              <Badge
                variant="outline"
                className="border-white/30 text-white bg-white/10"
              >
                {cat.emoji} {cat.label}
              </Badge>
              <Badge
                variant="outline"
                className="border-white/30 text-white bg-white/10"
              >
                <StatusIcon className="h-3 w-3" />
                {status.label}
              </Badge>
            </div>

            <h2 className="text-display text-2xl font-bold tracking-tight leading-tight">
              {ticket.event_title}
            </h2>

            <div className="mt-2 text-sm text-white/90">
              {ticket.event_college_name}
              {ticket.event_city && ` · ${ticket.event_city}`}
            </div>
          </div>

          {/* Perforation */}
          <div className="relative h-8 bg-[var(--bg-elevated)] -mt-4 rounded-t-3xl">
            <div className="absolute -left-4 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-[var(--bg-subtle)]" />
            <div className="absolute -right-4 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-[var(--bg-subtle)]" />
            <div className="absolute left-8 right-8 top-1/2 border-t-2 border-dashed border-[var(--border)]" />
          </div>

          {/* QR section */}
          <div className="px-6 pt-4 pb-6 text-center">
            {notYetValid ? (
              <div className="py-12 px-6">
                <Clock className="h-12 w-12 mx-auto text-[var(--color-warning)] mb-4" />
                <h3 className="text-display text-lg font-semibold">
                  QR not yet active
                </h3>
                <p className="mt-2 text-sm text-[var(--fg-muted)]">
                  QR will activate on{" "}
                  {new Date(ticket.valid_from!).toLocaleDateString("en-IN", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  })}
                </p>
                <div className="mt-4 text-xs text-[var(--fg-subtle)]">
                  {formatCountdown(validFrom - now)}
                </div>
              </div>
            ) : (
              <>
                <div className="inline-block p-4 rounded-2xl bg-white">
                  <QRCodeSVG
                    id={`qr-${ticket.id}`}
                    value={qrData}
                    size={220}
                    level="H"
                    marginSize={0}
                  />
                </div>
                <p className="mt-4 text-xs text-[var(--fg-muted)]">
                  Show this QR at the venue for check-in
                </p>
              </>
            )}

            {/* Ticket code */}
            <div className="mt-6 pt-6 border-t border-[var(--border)]">
              <div className="text-xs text-[var(--fg-subtle)]">Ticket Code</div>
              <div className="mt-1 text-display text-xl font-bold tracking-widest font-mono">
                {ticket.ticket_code}
              </div>
            </div>

            {/* Info grid */}
            <div className="mt-6 grid grid-cols-2 gap-4 text-left">
              <InfoBlock
                icon={Calendar}
                label="Dates"
                value={dates || "TBA"}
              />
              <InfoBlock
                icon={MapPin}
                label="Venue"
                value={ticket.event_venue || "TBA"}
              />
              <InfoBlock
                icon={TicketIcon}
                label="Attendee"
                value={ticket.user_name || "You"}
              />
              <InfoBlock
                icon={CheckCircle2}
                label="Status"
                value={status.label}
              />
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="mt-6 space-y-2">
          <Button
            variant="secondary"
            size="lg"
            className="w-full"
            onClick={downloadQR}
            disabled={notYetValid}
          >
            <Download className="h-4 w-4" />
            Download QR
          </Button>
          <Link href={`/events/${ticket.event_id}`} className="block">
            <Button variant="ghost" size="lg" className="w-full">
              View event details
            </Button>
          </Link>
        </div>

        <p className="mt-4 text-center text-xs text-[var(--fg-subtle)]">
          Keep this ticket safe. Do not share the QR with others.
        </p>
      </div>
    );
  }

  // ---- COMPACT (list) view ----
  return (
    <Link href={`/tickets/${ticket.id}`} className="block group">
      <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] overflow-hidden hover:border-[var(--border-strong)] hover:shadow-lg transition-all">
        <div className="flex flex-col sm:flex-row">
          <div
            className="relative w-full sm:w-40 shrink-0 aspect-[16/9] sm:aspect-auto overflow-hidden"
            style={{
              background: `linear-gradient(135deg, ${cat.color}30, ${cat.color}10)`,
            }}
          >
            {ticket.event_poster_url ? (
              <img
                src={ticket.event_poster_url}
                alt={ticket.event_title}
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-4xl">
                {cat.emoji}
              </div>
            )}
          </div>

          <div className="flex-1 p-4 sm:p-5">
            <div className="flex items-center gap-2 mb-2">
              <Badge variant={status.variant}>
                <StatusIcon className="h-3 w-3" />
                {status.label}
              </Badge>
              <Badge variant="outline">{cat.label}</Badge>
            </div>

            <h3 className="text-display font-semibold line-clamp-2 group-hover:text-[var(--color-brand)] transition-colors">
              {ticket.event_title}
            </h3>

            <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-[var(--fg-muted)]">
              {dates && (
                <span className="inline-flex items-center gap-1">
                  <Calendar className="h-3.5 w-3.5" />
                  {dates}
                </span>
              )}
              {ticket.event_venue && (
                <span className="inline-flex items-center gap-1">
                  <MapPin className="h-3.5 w-3.5" />
                  {ticket.event_venue}
                </span>
              )}
            </div>

            <div className="mt-3 flex items-center justify-between">
              <span className="text-xs font-mono text-[var(--fg-subtle)]">
                {ticket.ticket_code}
              </span>
              <span className="text-xs text-[var(--color-brand)] font-medium">
                View ticket →
              </span>
            </div>
          </div>
        </div>
      </div>
    </Link>
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
    <div className="flex items-start gap-2">
      <div className="grid h-8 w-8 shrink-0 place-items-center rounded-lg bg-[var(--bg-muted)]">
        <Icon className="h-3.5 w-3.5 text-[var(--fg-muted)]" />
      </div>
      <div className="min-w-0">
        <div className="text-[10px] uppercase tracking-wide text-[var(--fg-subtle)]">
          {label}
        </div>
        <div className="text-xs font-medium truncate">{value}</div>
      </div>
    </div>
  );
}

function formatCountdown(ms: number): string {
  if (ms <= 0) return "Activating now…";
  const days = Math.floor(ms / (1000 * 60 * 60 * 24));
  const hours = Math.floor((ms % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const mins = Math.floor((ms % (1000 * 60 * 60)) / (1000 * 60));
  if (days > 0) return `Activates in ${days}d ${hours}h`;
  if (hours > 0) return `Activates in ${hours}h ${mins}m`;
  return `Activates in ${mins}m`;
}