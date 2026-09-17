import type { Metadata } from "next";
import { TopNav } from "@/components/layout/top-nav";
import { BottomNav } from "@/components/layout/bottom-nav";
import { Footer } from "@/components/layout/footer";
import { EventsBrowser } from "@/components/event/events-browser";

export const metadata: Metadata = {
  title: "All events — CampusX",
  description:
    "Browse hackathons, cultural fests, sports meets, and workshops across colleges.",
};

export default function EventsPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />
      <main className="flex-1 pb-20 md:pb-0">
        <EventsBrowser />
      </main>
      <Footer />
      <BottomNav />
    </div>
  );
}