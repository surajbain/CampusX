"use client";

import * as React from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Eye, EyeOff, Loader2, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useLogin } from "@/lib/hooks/use-auth";

const schema = z.object({
  email: z.string().email("Enter a valid email"),
  password: z.string().min(1, "Password is required"),
});

type FormValues = z.infer<typeof schema>;

export function LoginForm() {
  const params = useSearchParams();
  const redirect = params.get("redirect") || undefined;
  const login = useLogin();
  const [showPassword, setShowPassword] = React.useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = (values: FormValues) => {
    login.mutate(values);
  };

  return (
    <div className="w-full max-w-md">
      <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-xl p-8">
        <div className="text-center mb-8">
          <h1 className="text-display text-2xl font-bold tracking-tight">
            Welcome back
          </h1>
          <p className="mt-1.5 text-sm text-[var(--fg-muted)]">
            Sign in to your CampusX account
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Email */}
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

          {/* Password */}
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <label className="text-sm font-medium">Password</label>
              <Link
                href="/forgot-password"
                className="text-xs text-[var(--color-brand)] hover:underline"
              >
                Forgot?
              </Link>
            </div>
            <div className="relative">
              <input
                {...register("password")}
                type={showPassword ? "text" : "password"}
                autoComplete="current-password"
                placeholder="••••••••"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 pr-10 text-sm placeholder:text-[var(--fg-subtle)] focus:outline-none focus:border-[var(--color-brand)] focus:ring-2 focus:ring-[var(--color-brand)]/20 transition-all"
              />
              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--fg-subtle)] hover:text-[var(--fg)] transition-colors"
                aria-label="Toggle password visibility"
              >
                {showPassword ? (
                  <EyeOff className="h-4 w-4" />
                ) : (
                  <Eye className="h-4 w-4" />
                )}
              </button>
            </div>
            {errors.password && (
              <p className="mt-1.5 text-xs text-[var(--color-danger)] flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.password.message}
              </p>
            )}
          </div>

          {/* Submit */}
          <Button
            type="submit"
            size="lg"
            className="w-full"
            disabled={isSubmitting || login.isPending}
          >
            {login.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Signing in…
              </>
            ) : (
              "Sign in"
            )}
          </Button>
        </form>

        {/* Demo credentials */}
        <div className="mt-6 pt-6 border-t border-[var(--border)]">
          <p className="text-xs text-[var(--fg-muted)] text-center mb-3">
            Try with demo credentials
          </p>
          <div className="grid grid-cols-1 gap-2 text-xs">
            <button
              type="button"
              onClick={() =>
                onSubmit({
                  email: "student@iitb.edu",
                  password: "Student@123",
                })
              }
              className="rounded-lg border border-[var(--border)] px-3 py-2 text-left hover:bg-[var(--bg-muted)] transition-colors"
            >
              <span className="font-medium">Student:</span>{" "}
              student@iitb.edu / Student@123
            </button>
            <button
              type="button"
              onClick={() =>
                onSubmit({
                  email: "admin@iitb.edu",
                  password: "Admin@123",
                })
              }
              className="rounded-lg border border-[var(--border)] px-3 py-2 text-left hover:bg-[var(--bg-muted)] transition-colors"
            >
              <span className="font-medium">College Admin:</span>{" "}
              admin@iitb.edu / Admin@123
            </button>
          </div>
        </div>
      </div>

      <p className="mt-6 text-center text-sm text-[var(--fg-muted)]">
        New to CampusX?{" "}
        <Link
          href={redirect ? `/register?redirect=${redirect}` : "/register"}
          className="font-medium text-[var(--color-brand)] hover:underline"
        >
          Create an account
        </Link>
      </p>
    </div>
  );
}