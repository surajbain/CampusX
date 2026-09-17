import type { Metadata } from "next";
import { Suspense } from "react";
import { LoginForm } from "@/components/auth/login-form";

export const metadata: Metadata = {
  title: "Sign in — CampusX",
  description: "Sign in to your CampusX account",
};

export default function LoginPage() {
  return (
    <Suspense fallback={<div className="text-sm text-[var(--fg-muted)]">Loading…</div>}>
      <LoginForm />
    </Suspense>
  );
}