"use client";

import * as React from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useQuery } from "@tanstack/react-query";
import { Eye, EyeOff, Loader2, AlertCircle, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useRegister } from "@/lib/hooks/use-auth";
import { listColleges } from "@/lib/api/colleges";

const schema = z.object({
  full_name: z.string().min(2, "Name must be at least 2 characters").max(200),
  email: z.string().email("Enter a valid email"),
  phone: z
    .string()
    .min(10, "Phone must be at least 10 digits")
    .max(30, "Too long")
    .regex(/^[0-9+\-\s()]+$/, "Invalid phone format")
    .optional()
    .or(z.literal("")),
  college_id: z.string().min(1, "Select a college"),
  password: z.string().min(8, "Password must be at least 8 characters").max(128),
});

type FormValues = z.infer<typeof schema>;

function passwordStrength(pw: string): { score: number; label: string; color: string } {
  let score = 0;
  if (pw.length >= 8) score++;
  if (/[a-z]/.test(pw)) score++;
  if (/[A-Z]/.test(pw)) score++;
  if (/[0-9]/.test(pw)) score++;
  if (/[^A-Za-z0-9]/.test(pw)) score++;

  if (score <= 2) return { score: 1, label: "Weak", color: "var(--color-danger)" };
  if (score === 3) return { score: 2, label: "Fair", color: "var(--color-warning)" };
  if (score === 4) return { score: 3, label: "Good", color: "var(--color-info)" };
  return { score: 4, label: "Strong", color: "var(--color-success)" };
}

export function RegisterForm() {
  const params = useSearchParams();
  const redirect = params.get("redirect") || undefined;
  const registerMutation = useRegister();
  const [showPassword, setShowPassword] = React.useState(false);
  const [agree, setAgree] = React.useState(false);

  const { data: collegesData } = useQuery({
    queryKey: ["colleges-all"],
    queryFn: () => listColleges({ limit: 100 }),
    staleTime: 10 * 60 * 1000,
  });

  const colleges = ((collegesData as any)?.data ?? []) as any[];

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      full_name: "",
      email: "",
      phone: "",
      college_id: "",
      password: "",
    },
  });

  const password = watch("password");
  const strength = password ? passwordStrength(password) : null;

  const onSubmit = (values: FormValues) => {
    if (!agree) return;
    registerMutation.mutate({
      email: values.email,
      password: values.password,
      full_name: values.full_name,
      phone: values.phone || undefined,
      college_id: values.college_id,
    });
  };

  return (
    <div className="w-full max-w-md">
      <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-xl p-8">
        <div className="text-center mb-8">
          <h1 className="text-display text-2xl font-bold tracking-tight">
            Create your account
          </h1>
          <p className="mt-1.5 text-sm text-[var(--fg-muted)]">
            Join CampusX and discover events near you
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">Full name</label>
            <input
              {...register("full_name")}
              type="text"
              autoComplete="name"
              placeholder="Rohan Sharma"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
            />
            {errors.full_name && (
              <p className="mt-1.5 text-xs text-[var(--color-danger)] flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.full_name.message}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">Email</label>
            <input
              {...register("email")}
              type="email"
              autoComplete="email"
              placeholder="you@college.edu"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
            />
            {errors.email && (
              <p className="mt-1.5 text-xs text-[var(--color-danger)] flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.email.message}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">
              Phone <span className="text-[var(--fg-subtle)] font-normal">(optional)</span>
            </label>
            <input
              {...register("phone")}
              type="tel"
              autoComplete="tel"
              placeholder="+91 98765 43210"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
            />
            {errors.phone && (
              <p className="mt-1.5 text-xs text-[var(--color-danger)] flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.phone.message}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">College</label>
            <select
              {...register("college_id")}
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
            >
              <option value="">Select your college</option>
              {colleges.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name} — {c.city}
                </option>
              ))}
            </select>
            {errors.college_id && (
              <p className="mt-1.5 text-xs text-[var(--color-danger)] flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.college_id.message}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">Password</label>
            <div className="relative">
              <input
                {...register("password")}
                type={showPassword ? "text" : "password"}
                autoComplete="new-password"
                placeholder="At least 8 characters"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 pr-10 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
              />
              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--fg-subtle)] hover:text-[var(--fg)] transition-colors"
                aria-label="Toggle password visibility"
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
            {strength && (
              <div className="mt-2 flex items-center gap-2">
                <div className="flex-1 h-1.5 rounded-full bg-[var(--bg-muted)] overflow-hidden">
                  <div
                    className="h-full transition-all duration-300"
                    style={{
                      width: `${(strength.score / 4) * 100}%`,
                      backgroundColor: strength.color,
                    }}
                  />
                </div>
                <span className="text-xs font-medium" style={{ color: strength.color }}>
                  {strength.label}
                </span>
              </div>
            )}
            {errors.password && (
              <p className="mt-1.5 text-xs text-[var(--color-danger)] flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.password.message}
              </p>
            )}
          </div>

          <label className="flex items-start gap-2 cursor-pointer select-none">
            <input
              type="checkbox"
              checked={agree}
              onChange={(e) => setAgree(e.target.checked)}
              className="mt-0.5 h-4 w-4 rounded border-[var(--border)] accent-[var(--color-brand)]"
            />
            <span className="text-xs text-[var(--fg-muted)]">
              I agree to the{" "}
              <Link href="/terms" className="text-[var(--color-brand)] hover:underline">
                Terms
              </Link>{" "}
              and{" "}
              <Link href="/privacy" className="text-[var(--color-brand)] hover:underline">
                Privacy Policy
              </Link>
            </span>
          </label>

          <Button
            type="submit"
            size="lg"
            className="w-full"
            disabled={!agree || isSubmitting || registerMutation.isPending}
          >
            {registerMutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Creating account…
              </>
            ) : (
              <>
                <Check className="h-4 w-4" />
                Create account
              </>
            )}
          </Button>
        </form>
      </div>

      <p className="mt-6 text-center text-sm text-[var(--fg-muted)]">
        Already have an account?{" "}
        <Link
          href={redirect ? `/login?redirect=${redirect}` : "/login"}
          className="font-medium text-[var(--color-brand)] hover:underline"
        >
          Sign in
        </Link>
      </p>
    </div>
  );
}