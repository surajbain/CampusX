"use client";

import * as React from "react";
import { Loader2, X, Calendar, Trophy, Users, MessageCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  useCreateEvent,
  useUpdateEvent,
  useAdminEvent,
} from "@/lib/hooks/use-admin-events";
import { CATEGORIES } from "@/lib/constants";

interface Props {
  id: string | null; // null = create
  onClose: () => void;
  onSuccess?: (id: string) => void;
}

const CATEGORY_OPTIONS = CATEGORIES.filter((c) => c.value !== "ALL");

export function EventForm({ id, onClose, onSuccess }: Props) {
  const isEdit = !!id;
  const { data: existing, isLoading: loadingExisting } = useAdminEvent(id || "");
  const createMutation = useCreateEvent();
  const updateMutation = useUpdateEvent();

  const [form, setForm] = React.useState({
    title: "",
    slug: "",
    description: "",
    category: "HACKATHON" as (typeof CATEGORY_OPTIONS)[number]["value"],
    poster_url: "",
    venue: "",
    city: "",
    starts_at: "",
    ends_at: "",
    registration_opens: "",
    registration_closes: "",
    price_paise: 0,
    capacity: 0,
    allow_teams: false,
    team_size_min: 2,
    team_size_max: 4,
    prize_pool_paise: 0,
    contact_email: "",
    whatsapp_link: "",
    rules: "",
    is_featured: false,
  });

  // Load existing values for edit
  React.useEffect(() => {
    if (existing) {
      setForm({
        title: existing.title || "",
        slug: existing.slug || "",
        description: existing.description || "",
        category: existing.category as any,
        poster_url: existing.poster_url || "",
        venue: existing.venue || "",
        city: existing.city || "",
        starts_at: toLocalInput(existing.starts_at),
        ends_at: toLocalInput(existing.ends_at),
        registration_opens: toLocalInput(existing.registration_opens),
        registration_closes: toLocalInput(existing.registration_closes),
        price_paise: existing.price_paise || 0,
        capacity: existing.capacity || 0,
        allow_teams: existing.allow_teams || false,
        team_size_min: existing.team_size_min || 2,
        team_size_max: existing.team_size_max || 4,
        prize_pool_paise: existing.prize_pool_paise || 0,
        contact_email: existing.contact_email || "",
        whatsapp_link: existing.whatsapp_link || "",
        rules: existing.rules || "",
        is_featured: existing.is_featured || false,
      });
    }
  }, [existing]);

  const isPending = createMutation.isPending || updateMutation.isPending;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const payload = {
      title: form.title,
      slug: form.slug || undefined,
      description: form.description || undefined,
      category: form.category,
      poster_url: form.poster_url || undefined,
      venue: form.venue || undefined,
      city: form.city || undefined,
      starts_at: form.starts_at ? new Date(form.starts_at).toISOString() : undefined,
      ends_at: form.ends_at ? new Date(form.ends_at).toISOString() : undefined,
      registration_opens: form.registration_opens
        ? new Date(form.registration_opens).toISOString()
        : undefined,
      registration_closes: form.registration_closes
        ? new Date(form.registration_closes).toISOString()
        : undefined,
      price_paise: form.price_paise,
      capacity: form.capacity || undefined,
      allow_teams: form.allow_teams,
      team_size_min: form.allow_teams ? form.team_size_min : undefined,
      team_size_max: form.allow_teams ? form.team_size_max : undefined,
      prize_pool_paise: form.prize_pool_paise,
      contact_email: form.contact_email || undefined,
      whatsapp_link: form.whatsapp_link || undefined,
      rules: form.rules || undefined,
      is_featured: form.is_featured,
    };

    if (isEdit) {
      updateMutation.mutate(
        { id: id!, input: payload },
        {
          onSuccess: (event) => {
            onSuccess?.(event.id);
            onClose();
          },
        }
      );
    } else {
      createMutation.mutate(payload, {
        onSuccess: (event) => {
          onSuccess?.(event.id);
          onClose();
        },
      });
    }
  };

  if (isEdit && loadingExisting) {
    return (
      <Card className="p-6">
        <div className="flex items-center gap-2 text-sm text-[var(--fg-muted)]">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading event…
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between mb-5">
        <h2 className="text-display text-lg font-semibold">
          {isEdit ? "Edit event" : "Create new event"}
        </h2>
        <Button variant="ghost" size="icon" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* ---- BASIC INFO ---- */}
        <Section title="Basic information" icon={Calendar}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Field label="Event title" required full>
              <input
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                required
                minLength={3}
                maxLength={250}
                placeholder="TechFest 2026 — Annual Tech Festival"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Custom slug" hint="Leave empty to auto-generate">
              <input
                value={form.slug}
                onChange={(e) =>
                  setForm({
                    ...form,
                    slug: e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "-"),
                  })
                }
                placeholder="techfest-2026"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm font-mono focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Category" required>
              <select
                value={form.category}
                onChange={(e) =>
                  setForm({ ...form, category: e.target.value as any })
                }
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              >
                {CATEGORY_OPTIONS.map((c) => (
                  <option key={c.value} value={c.value}>
                    {c.emoji} {c.label}
                  </option>
                ))}
              </select>
            </Field>

            <Field label="Description" full>
              <textarea
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                rows={4}
                maxLength={20000}
                placeholder="Describe your event — what, who, why…"
                className="w-full rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] p-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)] resize-y"
              />
            </Field>

            <Field label="Poster URL" hint="Paste image URL (later: upload)" full>
              <input
                value={form.poster_url}
                onChange={(e) => setForm({ ...form, poster_url: e.target.value })}
                placeholder="https://example.com/poster.jpg"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>
          </div>
        </Section>

        {/* ---- WHEN & WHERE ---- */}
        <Section title="When & where" icon={Calendar}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Field label="Starts at" required>
              <input
                type="datetime-local"
                value={form.starts_at}
                onChange={(e) => setForm({ ...form, starts_at: e.target.value })}
                required
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Ends at">
              <input
                type="datetime-local"
                value={form.ends_at}
                onChange={(e) => setForm({ ...form, ends_at: e.target.value })}
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Registration opens">
              <input
                type="datetime-local"
                value={form.registration_opens}
                onChange={(e) =>
                  setForm({ ...form, registration_opens: e.target.value })
                }
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Registration closes">
              <input
                type="datetime-local"
                value={form.registration_closes}
                onChange={(e) =>
                  setForm({ ...form, registration_closes: e.target.value })
                }
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Venue">
              <input
                value={form.venue}
                onChange={(e) => setForm({ ...form, venue: e.target.value })}
                placeholder="Main Auditorium"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="City">
              <input
                value={form.city}
                onChange={(e) => setForm({ ...form, city: e.target.value })}
                placeholder="Mumbai"
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>
          </div>
        </Section>

        {/* ---- PRICING & CAPACITY ---- */}
        <Section title="Pricing & capacity" icon={Trophy}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Field label="Price (in ₹)" hint="0 = FREE">
              <input
                type="number"
                min={0}
                value={form.price_paise / 100}
                onChange={(e) =>
                  setForm({
                    ...form,
                    price_paise: Math.round(Number(e.target.value) * 100),
                  })
                }
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Capacity" hint="Leave 0 for unlimited">
              <input
                type="number"
                min={0}
                value={form.capacity}
                onChange={(e) =>
                  setForm({ ...form, capacity: Number(e.target.value) })
                }
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Prize pool (in ₹)">
              <input
                type="number"
                min={0}
                value={form.prize_pool_paise / 100}
                onChange={(e) =>
                  setForm({
                    ...form,
                    prize_pool_paise: Math.round(Number(e.target.value) * 100),
                  })
                }
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
          </div>

          {/* Teams */}
          <div className="mt-4 pt-4 border-t border-[var(--border)]">
            <label className="flex items-center gap-2 cursor-pointer select-none">
              <input
                type="checkbox"
                checked={form.allow_teams}
                onChange={(e) =>
                  setForm({ ...form, allow_teams: e.target.checked })
                }
                className="h-4 w-4 accent-[var(--color-brand)]"
              />
              <Users className="h-4 w-4 text-[var(--fg-muted)]" />
              <span className="text-sm font-medium">
                Allow team participation
              </span>
            </label>

            {form.allow_teams && (
              <div className="mt-3 grid grid-cols-2 gap-4">
                <Field label="Min team size">
                  <input
                    type="number"
                    min={1}
                    max={50}
                    value={form.team_size_min}
                    onChange={(e) =>
                      setForm({
                        ...form,
                        team_size_min: Number(e.target.value),
                      })
                    }
                    className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
                  />
                </Field>
                <Field label="Max team size">
                  <input
                    type="number"
                    min={1}
                    max={50}
                    value={form.team_size_max}
                    onChange={(e) =>
                      setForm({
                        ...form,
                        team_size_max: Number(e.target.value),
                      })
                    }
                    className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
                  />
                </Field>
              </div>
            )}
          </div>
        </Section>

        {/* ---- WHATSAPP & RULES ---- */}
        <Section title="Communication & rules" icon={MessageCircle}>
          <div className="grid grid-cols-1 gap-4">
            <Field
              label="WhatsApp group link"
              hint="Students will see this after registering"
              full
            >
              <input
                value={form.whatsapp_link}
                onChange={(e) =>
                  setForm({ ...form, whatsapp_link: e.target.value })
                }
                placeholder="https://chat.whatsapp.com/..."
                className="w-full h-11 rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] px-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)]"
              />
            </Field>

            <Field label="Rules" full>
              <textarea
                value={form.rules}
                onChange={(e) => setForm({ ...form, rules: e.target.value })}
                rows={5}
                placeholder="Event rules, eligibility, code of conduct…"
                className="w-full rounded-xl border border-[var(--border)] bg-[var(--bg-subtle)] p-3.5 text-sm focus:outline-none focus:border-[var(--color-brand)] resize-y"
              />
            </Field>

            <label className="flex items-center gap-2 cursor-pointer select-none">
              <input
                type="checkbox"
                checked={form.is_featured}
                onChange={(e) =>
                  setForm({ ...form, is_featured: e.target.checked })
                }
                className="h-4 w-4 accent-[var(--color-brand)]"
              />
              <span className="text-sm font-medium">
                Feature this event on homepage
              </span>
            </label>
          </div>
        </Section>

        {/* ---- SUBMIT ---- */}
        <div className="flex justify-end gap-2 pt-2 border-t border-[var(--border)]">
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
              "Save changes"
            ) : (
              "Create event"
            )}
          </Button>
        </div>
      </form>
    </Card>
  );
}

// ---- helpers ----

function Field({
  label,
  hint,
  required,
  full,
  children,
}: {
  label: string;
  hint?: string;
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
      {hint && (
        <div className="mt-1 text-xs text-[var(--fg-subtle)]">{hint}</div>
      )}
    </div>
  );
}

function Section({
  title,
  icon: Icon,
  children,
}: {
  title: string;
  icon: React.ComponentType<{ className?: string }>;
  children: React.ReactNode;
}) {
  return (
    <div>
      <div className="flex items-center gap-2 mb-4">
        <div className="grid h-8 w-8 place-items-center rounded-lg bg-[var(--color-brand)]/10">
          <Icon className="h-4 w-4 text-[var(--color-brand)]" />
        </div>
        <h3 className="text-display text-base font-semibold">{title}</h3>
      </div>
      {children}
    </div>
  );
}

// Convert ISO string → "YYYY-MM-DDTHH:MM" for datetime-local input
function toLocalInput(iso?: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(
    d.getHours()
  )}:${pad(d.getMinutes())}`;
}