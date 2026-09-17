import Link from "next/link";
import { GraduationCap } from "lucide-react";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen relative overflow-hidden">
      {/* Gradient background */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute inset-0 bg-gradient-to-br from-[var(--color-brand-from)]/10 via-[var(--color-brand-to)]/5 to-transparent" />
        <div className="absolute top-0 left-1/2 -translate-x-1/2 h-[500px] w-[700px] rounded-full bg-[var(--color-brand-from)]/20 blur-[120px] opacity-60" />
      </div>

      <div className="min-h-screen flex flex-col">
        {/* Simple header */}
        <header className="cx-container py-6">
          <Link href="/" className="inline-flex items-center gap-2 group">
            <div className="grid h-9 w-9 place-items-center rounded-xl bg-brand-gradient shadow-md transition-transform group-hover:scale-105">
              <GraduationCap className="h-5 w-5 text-white" />
            </div>
            <span className="text-display text-xl font-bold tracking-tight">
              Campus<span className="text-brand-gradient">X</span>
            </span>
          </Link>
        </header>

        {/* Content */}
        <main className="flex-1 flex items-center justify-center px-4 py-8">
          {children}
        </main>

        {/* Footer */}
        <footer className="cx-container py-6 text-center text-xs text-[var(--fg-subtle)]">
          © 2026 CampusX · Built for students, by students
        </footer>
      </div>
    </div>
  );
}