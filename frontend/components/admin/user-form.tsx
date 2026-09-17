"use client";

import * as React from "react";
import { Loader2, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useCreateCollegeUser } from "@/lib/hooks/use-admin";

export function UserForm({
  collegeId,
  onClose,
}: {
  collegeId: string;
  onClose: () => void;
}) {
  const createMutation = useCreateCollegeUser();
  const [form, setForm] = React.useState({
    full_name: "",
    email: "",
    password: "",
    phone: "",
    role: "ORGANIZER" as "ORGANIZER" | "VOLUNTEER" | "COLLEGE_ADMIN",
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createMutation.mutate(
      {
        collegeId,
        input: {
          email: form.email,
          password: form.password,
          full_name: form.full_name,
          phone: form.phone || undefined,
          role: form.role,
        },
      },
      {
        onSuccess: () => {
          onClose();
        },
      }
    );
  };

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between mb-5">
        <h2 className="text-display text-lg font-semibold">Add new user</h2>
        <Button variant="ghost" size="icon" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">
              Full name <span className="text-[var(--color-danger)]">*</span>
            </label>
            <input
              value={form.full_name}
              onChange={(e) => setForm({ ...form, full_name: e.target.value })}
              required
              placeholder="Aarav Sharma"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">
              Email <span className="text-[var(--color-danger)]">*</span>
            </label>
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
              placeholder="aarav@college.edu"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">
              Password <span className="text-[var(--color-danger)]">*</span>
            </label>
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              required
              minLength={8}
              placeholder="At least 8 characters"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">Phone</label>
            <input
              value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })}
              placeholder="+91 98765 43210"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </div>

          <div className="md:col-span-2">
            <label className="block text-sm font-medium mb-1.5">
              Role <span className="text-[var(--color-danger)]">*</span>
            </label>
            <select
              value={form.role}
              onChange={(e) =>
                setForm({
                  ...form,
                  role: e.target.value as typeof form.role,
                })
              }
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            >
              <option value="ORGANIZER">Organizer — can create events</option>
              <option value="VOLUNTEER">Volunteer — can scan tickets</option>
              <option value="COLLEGE_ADMIN">College Admin — full access</option>
            </select>
          </div>
        </div>

        <div className="flex justify-end gap-2 pt-2">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={createMutation.isPending}>
            {createMutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Creating…
              </>
            ) : (
              "Create user"
            )}
          </Button>
        </div>
      </form>
    </Card>
  );
}