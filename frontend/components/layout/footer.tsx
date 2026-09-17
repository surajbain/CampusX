import Link from "next/link";
import { GraduationCap, Globe, MessageCircle, AtSign } from "lucide-react";

export function Footer() {
  return (
    <footer className="mt-24 border-t border-[var(--border)] bg-[var(--bg-subtle)]">
      <div className="cx-container py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-10">
          <div className="md:col-span-2">
            <Link href="/" className="flex items-center gap-2">
              <div className="grid h-8 w-8 place-items-center rounded-lg bg-brand-gradient">
                <GraduationCap className="h-4 w-4 text-white" />
              </div>
              <span className="text-display text-lg font-bold">
                Campus<span className="text-brand-gradient">X</span>
              </span>
            </Link>
            <p className="mt-4 max-w-sm text-sm text-[var(--fg-muted)]">
              Discover hackathons, cultural fests, sports meets, and workshops across colleges. Free for students, free for colleges.
            </p>
            <div className="mt-5 flex gap-2">
              <a href="#" aria-label="Social" className="grid h-9 w-9 place-items-center rounded-lg border border-[var(--border)] hover:bg-[var(--bg-muted)] transition-colors">
                <AtSign className="h-4 w-4" />
              </a>
              <a href="#" aria-label="Community" className="grid h-9 w-9 place-items-center rounded-lg border border-[var(--border)] hover:bg-[var(--bg-muted)] transition-colors">
                <MessageCircle className="h-4 w-4" />
              </a>
              <a href="#" aria-label="Website" className="grid h-9 w-9 place-items-center rounded-lg border border-[var(--border)] hover:bg-[var(--bg-muted)] transition-colors">
                <Globe className="h-4 w-4" />
              </a>
            </div>
          </div>

          <div>
            <h4 className="text-sm font-semibold mb-3">Product</h4>
            <ul className="space-y-2 text-sm text-[var(--fg-muted)]">
              <li><Link href="/events" className="hover:text-[var(--fg)] transition-colors">Events</Link></li>
              <li><Link href="/colleges" className="hover:text-[var(--fg)] transition-colors">Colleges</Link></li>
              <li><Link href="/search" className="hover:text-[var(--fg)] transition-colors">Search</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-semibold mb-3">Company</h4>
            <ul className="space-y-2 text-sm text-[var(--fg-muted)]">
              <li><Link href="/about" className="hover:text-[var(--fg)] transition-colors">About</Link></li>
              <li><Link href="/contact" className="hover:text-[var(--fg)] transition-colors">Contact</Link></li>
              <li><Link href="/privacy" className="hover:text-[var(--fg)] transition-colors">Privacy</Link></li>
            </ul>
          </div>
        </div>

        <div className="mt-10 pt-6 border-t border-[var(--border)] flex flex-col sm:flex-row justify-between gap-3 text-xs text-[var(--fg-subtle)]">
          <p>© 2026 CampusX. Built for students, by students.</p>
          <p>Made with ❤️ in India</p>
        </div>
      </div>
    </footer>
  );
}