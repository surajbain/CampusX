import { Suspense } from "react";
import type { Metadata } from "next";
import { TopNav } from "@/components/layout/top-nav";
import { BottomNav } from "@/components/layout/bottom-nav";
import { Footer } from "@/components/layout/footer";
import { EventsBrowser } from "@/components/event/events-browser";
import { EventGridSkeleton } from "@/components/ui/skeleton";

export const metadata: Metadata = {
  title: "Search events — CampusX",
  description: "Search hackathons, fests, sports meets, and workshops.",
};

export default function SearchPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />
      <main className="flex-1 pb-20 md:pb-0">
        <Suspense fallback={<div className="cx-container py-12"><EventGridSkeleton count={6} /></div>}>
          <EventsBrowser />
        </Suspense>
      </main>
      <Footer />
      <BottomNav />
    </div>
  );
}