"use client";

import * as React from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Users, MapPin, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useCollege } from "@/lib/hooks/use-admin";

export default function CollegeDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { data: college, isLoading } = useCollege(id);

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 text-sm text-[var(--fg-muted)]">
        <Loader2 className="h-4 w-4 animate-spin" /> Loading…
      </div>
    );
  }

  if (!college) {
    return (
      <div className="text-center py-16">
        <p className="text-[var(--fg-muted)]">College not found</p>
        <Link href="/admin/colleges" className="mt-4 inline-block">
          <Button variant="secondary">Back to colleges</Button>
        </Link>
      </div>
    );
  }

  return (
    <div>
      <Link
        href="/admin/colleges"
        className="inline-flex items-center gap-2 text-sm text-[var(--fg-muted)] hover:text-[var(--fg)] mb-6 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to colleges
      </Link>

      <Card className="p-6 mb-6">
        <div className="flex items-start gap-4">
          <div className="grid h-16 w-16 shrink-0 place-items-center rounded-2xl bg-brand-gradient text-white font-bold text-2xl">
            {college.name?.[0]?.toUpperCase() ?? "?"}
          </div>
          <div className="flex-1">
            <h1 className="text-display text-2xl font-bold tracking-tight">
              {college.name}
            </h1>
            <div className="mt-1.5 flex items-center gap-2 text-sm text-[var(--fg-muted)]">
              <MapPin className="h-3.5 w-3.5" />
              {college.city}
              {college.state && `, ${college.state}`}
            </div>
            <div className="mt-3 flex gap-2">
              <Badge variant="success">Active</Badge>
              <Badge variant="outline" className="font-mono">
                {college.slug}
              </Badge>
            </div>
          </div>

          <Link href={`/admin/colleges/${id}/users`}>
            <Button>
              <Users className="h-4 w-4" />
              Manage users
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  );
}