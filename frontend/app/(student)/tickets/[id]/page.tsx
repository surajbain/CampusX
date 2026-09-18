"use client";

import * as React from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Loader2, TicketX } from "lucide-react";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { TicketCard } from "@/components/ticket/ticket-card";
import { useTicket } from "@/lib/hooks/use-tickets";
import { WhatsAppButton } from "@/components/ticket/whatsapp-button";

export default function TicketDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { data: ticket, isLoading, isError, error } = useTicket(id);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <Loader2 className="h-6 w-6 animate-spin text-[var(--fg-muted)]" />
      </div>
    );
  }

  if (isError || !ticket) {
    return (
      <div className="max-w-md mx-auto">
        <EmptyState
          icon={TicketX}
          title="Ticket not found"
          description={
            error instanceof Error
              ? error.message
              : "This ticket doesn't exist or you don't have access."
          }
          action={
            <Link href="/tickets">
              <Button>
                <ArrowLeft className="h-4 w-4" />
                Back to tickets
              </Button>
            </Link>
          }
        />
      </div>
    );
  }

    return (
    <div>
      <Link
        href="/tickets"
        className="inline-flex items-center gap-2 text-sm text-[var(--fg-muted)] hover:text-[var(--fg)] mb-6 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to tickets
      </Link>

      <TicketCard ticket={ticket} showQR />

      {/* WhatsApp group button — only if link exists */}
      {ticket.event_whatsapp_link && (
        <div className="mt-4 max-w-md mx-auto">
          <WhatsAppButton link={ticket.event_whatsapp_link} fullWidth />
          <p className="mt-2 text-center text-xs text-[var(--fg-subtle)]">
            Join the event WhatsApp group for updates
          </p>
        </div>
      )}
    </div>
  );
}