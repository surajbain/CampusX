import { TopNav } from "@/components/layout/top-nav";
import { BottomNav } from "@/components/layout/bottom-nav";
import { Footer } from "@/components/layout/footer";
import { Hero } from "@/components/home/hero";
import { EventsSection } from "@/components/home/events-section";

export default function HomePage() {
  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />

      <main className="flex-1 pb-20 md:pb-0">
        <Hero />
        <EventsSection />
      </main>

      <Footer />
      <BottomNav />
    </div>
  );
}