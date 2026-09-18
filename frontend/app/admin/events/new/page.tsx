"use client";

import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { EventForm } from "@/components/admin/event-form";

export default function NewEventPage() {
  const router = useRouter();

  return (
    <div>
      <Link
        href="/admin/events"
        className="inline-flex items-center gap-2 text-sm text-[var(--fg-muted)] hover:text-[var(--fg)] mb-6 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to events
      </Link>

      <EventForm
        id={null}
        onClose={() => router.push("/admin/events")}
        onSuccess={(id) => router.push(`/admin/events/${id}/edit`)}
      />
    </div>
  );
}