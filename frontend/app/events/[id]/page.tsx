import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { TopNav } from "@/components/layout/top-nav";
import { BottomNav } from "@/components/layout/bottom-nav";
import { Footer } from "@/components/layout/footer";
import { EventDetailView } from "@/components/event/event-detail-view";
import { getEvent } from "@/lib/api/events";
import { formatDateRange } from "@/lib/utils";
import type { Event } from "@/lib/api/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchEventServer(id: string): Promise<Event | null> {
  try {
    const res = await fetch(`${API_URL}/api/v1/events/${id}`, {
      next: { revalidate: 60 },
    });
    if (!res.ok) return null;
    const body = await res.json();
    return body?.data ?? null;
  } catch {
    return null;
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}): Promise<Metadata> {
  const { id } = await params;
  const event = await fetchEventServer(id);
  if (!event) return { title: "Event not found — CampusX" };

  const dates = formatDateRange(event.starts_at, event.ends_at);
  const description = event.description
    ? event.description.slice(0, 155)
    : `${event.category} event at ${event.college_name}. ${dates}`;

  return {
    title: `${event.title} — CampusX`,
    description,
    openGraph: {
      title: event.title,
      description,
      images: event.poster_url ? [event.poster_url] : [],
      type: "website",
    },
    twitter: {
      card: "summary_large_image",
      title: event.title,
      description,
    },
  };
}

export default async function EventDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const event = await fetchEventServer(id);

  if (!event) notFound();

  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />

      <main className="flex-1 pb-20 md:pb-0">
        <EventDetailView event={event} />
      </main>

      <Footer />
      <BottomNav />
    </div>
  );
}