"use client";

import * as React from "react";
import Link from "next/link";
import {
  Plus,
  Calendar,
  MapPin,
  Pencil,
  Trash2,
  Send,
  Ban,
  Loader2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import {
  useAdminEvents,
  useDeleteEvent,
  usePublishEvent,
  useCancelEvent,
} from "@/lib/hooks/use-admin-events";
import { EventForm } from "@/components/admin/event-form";
import { categoryMeta } from "@/lib/constants";
import { formatPrice, formatDateRange } from "@/lib/utils";

const statusVariants: Record<
  string,
  "default" | "success" | "warning" | "danger" | "brand"
> = {
  DRAFT: "default",
  PUBLISHED: "success",
  ONGOING: "brand",
  COMPLETED: "outline" as any,
  CANCELLED: "danger",
};

export default function AdminEventsPage() {
  const [editingId, setEditingId] = React.useState<string | null>(null);
  const [showCreate, setShowCreate] = React.useState(false);

  const { data, isLoading } = useAdminEvents({ limit: 100 });
  const deleteMutation = useDeleteEvent();
  const publishMutation = usePublishEvent();
  const cancelMutation = useCancelEvent();

  const events = ((data as any)?.data ?? []) as any[];

  return (
    <div>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h1 className="text-display text-2xl font-bold tracking-tight">
            Events
          </h1>
          <p className="mt-1 text-sm text-[var(--fg-muted)]">
            {events.length} events you manage
          </p>
        </div>
        <Button
          onClick={() => {
            setShowCreate(true);
            setEditingId(null);
          }}
        >
          <Plus className="h-4 w-4" />
          Create event
        </Button>
      </div>

      {(showCreate || editingId) && (
        <div className="mb-6">
          <EventForm
            id={editingId}
            onClose={() => {
              setShowCreate(false);
              setEditingId(null);
            }}
          />
        </div>
      )}

      {isLoading && (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-24 rounded-xl" />
          ))}
        </div>
      )}

      {!isLoading && events.length === 0 && (
        <EmptyState
          icon={Calendar}
          title="No events yet"
          description="Create your first event to start accepting registrations."
          action={
            <Button onClick={() => setShowCreate(true)}>
              <Plus className="h-4 w-4" />
              Create event
            </Button>
          }
        />
      )}

      {!isLoading && events.length > 0 && (
        <div className="space-y-3">
          {events.map((e) => {
            const cat = categoryMeta(e.category);
            const dates = formatDateRange(e.starts_at, e.ends_at);
            const price = formatPrice(e.price_paise, e.currency);

            return (
              <Card key={e.id} className="p-4">
                <div className="flex flex-col sm:flex-row sm:items-center gap-4">
                  {/* Icon */}
                  <div
                    className="grid h-12 w-12 shrink-0 place-items-center rounded-xl text-2xl"
                    style={{
                      background: `linear-gradient(135deg, ${cat.color}30, ${cat.color}10)`,
                    }}
                  >
                    {cat.emoji}
                  </div>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <h3 className="text-display font-semibold truncate">
                        {e.title}
                      </h3>
                      <Badge variant={statusVariants[e.status] || "default"}>
                        {e.status}
                      </Badge>
                      <Badge variant="outline">{cat.label}</Badge>
                      <Badge
                        variant={e.price_paise === 0 ? "success" : "brand"}
                      >
                        {price}
                      </Badge>
                    </div>
                    <div className="mt-1 flex flex-wrap gap-x-3 gap-y-1 text-xs text-[var(--fg-muted)]">
                      {dates && (
                        <span className="inline-flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {dates}
                        </span>
                      )}
                      {e.venue && (
                        <span className="inline-flex items-center gap-1">
                          <MapPin className="h-3 w-3" />
                          {e.venue}
                        </span>
                      )}
                      <span className="font-mono text-[var(--fg-subtle)]">
                        /{e.slug}
                      </span>
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-1 shrink-0">
                    {e.status === "DRAFT" && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => publishMutation.mutate(e.id)}
                        disabled={publishMutation.isPending}
                        title="Publish"
                      >
                        {publishMutation.isPending ? (
                          <Loader2 className="h-3.5 w-3.5 animate-spin" />
                        ) : (
                          <Send className="h-3.5 w-3.5 text-[var(--color-success)]" />
                        )}
                      </Button>
                    )}
                    {e.status === "PUBLISHED" && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => {
                          if (confirm("Cancel this event?")) {
                            cancelMutation.mutate(e.id);
                          }
                        }}
                        title="Cancel event"
                      >
                        <Ban className="h-3.5 w-3.5 text-[var(--color-warning)]" />
                      </Button>
                    )}
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        setEditingId(e.id);
                        setShowCreate(false);
                      }}
                      title="Edit"
                    >
                      <Pencil className="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        if (confirm(`Delete "${e.title}"? This cannot be undone.`)) {
                          deleteMutation.mutate(e.id);
                        }
                      }}
                      title="Delete"
                    >
                      <Trash2 className="h-3.5 w-3.5 text-[var(--color-danger)]" />
                    </Button>
                  </div>
                </div>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}