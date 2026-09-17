"use client";

import * as React from "react";
import { Loader2, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  useCreateCollege,
  useUpdateCollege,
  useCollege,
} from "@/lib/hooks/use-admin";

export function CollegeForm({
  id,
  onClose,
}: {
  id: string | null;
  onClose: () => void;
}) {
  const { data: existing, isLoading: loadingExisting } = useCollege(id || "");
  const createMutation = useCreateCollege();
  const updateMutation = useUpdateCollege();

  const [form, setForm] = React.useState({
    name: "",
    slug: "",
    city: "",
    state: "",
    contact_email: "",
    contact_phone: "",
    website: "",
  });

  React.useEffect(() => {
    if (existing) {
      setForm({
        name: existing.name || "",
        slug: existing.slug || "",
        city: existing.city || "",
        state: existing.state || "",
        contact_email: existing.contact_email || "",
        contact_phone: existing.contact_phone || "",
        website: existing.website || "",
      });
    }
  }, [existing]);

  const isEdit = !!id;
  const isPending = createMutation.isPending || updateMutation.isPending;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isEdit) {
      updateMutation.mutate(
        {
          id: id!,
          input: {
            name: form.name,
            city: form.city,
            state: form.state,
            contact_email: form.contact_email,
            contact_phone: form.contact_phone,
            website: form.website,
          },
        },
        { onSuccess: () => onClose() }
      );
    } else {
      createMutation.mutate(
        {
          name: form.name,
          slug: form.slug,
          city: form.city,
          state: form.state || undefined,
          contact_email: form.contact_email || undefined,
          contact_phone: form.contact_phone || undefined,
          website: form.website || undefined,
        },
        { onSuccess: () => onClose() }
      );
    }
  };

  if (isEdit && loadingExisting) {
    return (
      <Card className="p-6">
        <div className="flex items-center gap-2 text-sm text-[var(--fg-muted)]">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading…
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between mb-5">
        <h2 className="text-display text-lg font-semibold">
          {isEdit ? "Edit college" : "Add new college"}
        </h2>
        <Button variant="ghost" size="icon" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Field label="Name" required>
            <input
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
              placeholder="IIT Bombay"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </Field>

          {!isEdit && (
            <Field label="Slug" required>
              <input
                value={form.slug}
                onChange={(e) =>
                  setForm({
                    ...form,
                    slug: e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "-"),
                  })
                }
                required
                placeholder="iitb"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm font-mono focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>
          )}

          <Field label="City" required>
            <input
              value={form.city}
              onChange={(e) => setForm({ ...form, city: e.target.value })}
              required
              placeholder="Mumbai"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </Field>

          <Field label="State">
            <input
              value={form.state}
              onChange={(e) => setForm({ ...form, state: e.target.value })}
              placeholder="Maharashtra"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </Field>

          <Field label="Contact email">
            <input
              type="email"
              value={form.contact_email}
              onChange={(e) =>
                setForm({ ...form, contact_email: e.target.value })
              }
              placeholder="events@college.edu"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </Field>

          <Field label="Contact phone">
            <input
              value={form.contact_phone}
              onChange={(e) =>
                setForm({ ...form, contact_phone: e.target.value })
              }
              placeholder="+91 98765 43210"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </Field>

          <Field label="Website" full>
            <input
              value={form.website}
              onChange={(e) => setForm({ ...form, website: e.target.value })}
              placeholder="https://college.edu"
              className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
            />
          </Field>
        </div>

        <div className="flex justify-end gap-2 pt-2">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                {isEdit ? "Updating…" : "Creating…"}
              </>
            ) : isEdit ? (
              "Update"
            ) : (
              "Create college"
            )}
          </Button>
        </div>
      </form>
    </Card>
  );
}

function Field({
  label,
  required,
  full,
  children,
}: {
  label: string;
  required?: boolean;
  full?: boolean;
  children: React.ReactNode;
}) {
  return (
    <div className={full ? "md:col-span-2" : ""}>
      <label className="block text-sm font-medium mb-1.5">
        {label}
        {required && <span className="text-[var(--color-danger)] ml-0.5">*</span>}
      </label>
      {children}
    </div>
  );
}