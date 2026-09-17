"use client";

import * as React from "react";
import Link from "next/link";
import { Plus, Search, MapPin, Building2, MoreHorizontal, Trash2, Pencil } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { listColleges } from "@/lib/api/colleges";
import { useDeleteCollege } from "@/lib/hooks/use-admin";
import { CollegeForm } from "@/components/admin/college-form";

export default function AdminCollegesPage() {
  const [q, setQ] = React.useState("");
  const [showCreate, setShowCreate] = React.useState(false);
  const [editingId, setEditingId] = React.useState<string | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ["colleges", "admin"],
    queryFn: () => listColleges({ limit: 100 }),
  });

  const deleteMutation = useDeleteCollege();
  const colleges = ((data as any)?.data ?? []) as any[];

  const filtered = q
    ? colleges.filter(
        (c) =>
          c.name.toLowerCase().includes(q.toLowerCase()) ||
          c.city.toLowerCase().includes(q.toLowerCase())
      )
    : colleges;

  return (
    <div>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h1 className="text-display text-2xl font-bold tracking-tight">
            Colleges
          </h1>
          <p className="mt-1 text-sm text-[var(--fg-muted)]">
            {colleges.length} colleges on CampusX
          </p>
        </div>
        <Button onClick={() => { setShowCreate(true); setEditingId(null); }} size="md">
          <Plus className="h-4 w-4" />
          Add college
        </Button>
      </div>

      {/* Create/Edit form */}
      {(showCreate || editingId) && (
        <div className="mb-6">
          <CollegeForm
            id={editingId}
            onClose={() => {
              setShowCreate(false);
              setEditingId(null);
            }}
          />
        </div>
      )}

      {/* Search */}
      <div className="mb-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[var(--fg-subtle)]" />
          <input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="Search colleges…"
            className="w-full h-10 rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)] pl-10 pr-4 text-sm focus:outline-none focus:border-[var(--color-brand)]"
          />
        </div>
      </div>

      {/* List */}
      {isLoading && (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-20 rounded-xl" />
          ))}
        </div>
      )}

      {!isLoading && filtered.length === 0 && (
        <EmptyState
          icon={Building2}
          title="No colleges yet"
          description="Add your first college to get started."
          action={
            <Button onClick={() => setShowCreate(true)}>
              <Plus className="h-4 w-4" />
              Add college
            </Button>
          }
        />
      )}

      {!isLoading && filtered.length > 0 && (
        <div className="space-y-3">
          {filtered.map((c) => (
            <Card key={c.id} className="p-4">
              <div className="flex items-center gap-4">
                <div className="grid h-11 w-11 shrink-0 place-items-center rounded-xl bg-brand-gradient text-white font-bold">
                  {c.name?.[0]?.toUpperCase() ?? "?"}
                </div>

                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h3 className="text-display font-semibold truncate">
                      {c.name}
                    </h3>
                    <Badge variant="success">Active</Badge>
                  </div>
                  <div className="mt-0.5 flex items-center gap-1 text-xs text-[var(--fg-muted)]">
                    <MapPin className="h-3 w-3" />
                    {c.city}
                    {c.state && `, ${c.state}`}
                    <span className="text-[var(--fg-subtle)]">·</span>
                    <span className="font-mono">{c.slug}</span>
                  </div>
                </div>

                <div className="flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      setEditingId(c.id);
                      setShowCreate(false);
                    }}
                  >
                    <Pencil className="h-3.5 w-3.5" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      if (confirm(`Delete ${c.name}? This cannot be undone.`)) {
                        deleteMutation.mutate(c.id);
                      }
                    }}
                  >
                    <Trash2 className="h-3.5 w-3.5 text-[var(--color-danger)]" />
                  </Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}