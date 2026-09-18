"use client";

import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { EventForm } from "@/components/admin/event-form";

export default function EditEventPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

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
        id={id}
        onClose={() => router.push("/admin/events")}
      />
    </div>
  );
}