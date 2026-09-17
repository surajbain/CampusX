import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatPrice(paise: number, currency = "INR"): string {
  if (paise === 0) return "FREE";
  const rupees = paise / 100;
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(rupees);
}

export function formatDateRange(
  start?: string | Date | null,
  end?: string | Date | null
): string {
  if (!start) return "";
  const s = typeof start === "string" ? new Date(start) : start;
  const e = end ? (typeof end === "string" ? new Date(end) : end) : null;

  const opts: Intl.DateTimeFormatOptions = { month: "short", day: "numeric" };
  if (!e) return s.toLocaleDateString("en-IN", { ...opts, year: "numeric" });

  const sameYear = s.getFullYear() === e.getFullYear();
  const sameMonth = sameYear && s.getMonth() === e.getMonth();
  const sameDay = sameMonth && s.getDate() === e.getDate();

  if (sameDay) return s.toLocaleDateString("en-IN", { ...opts, year: "numeric" });
  if (sameMonth)
    return `${s.toLocaleDateString("en-IN", opts)} – ${e.toLocaleDateString("en-IN", { ...opts, year: "numeric" })}`;
  return `${s.toLocaleDateString("en-IN", opts)} – ${e.toLocaleDateString("en-IN", { ...opts, year: "numeric" })}`;
}

export function formatRelative(date: string | Date): string {
  const d = typeof date === "string" ? new Date(date) : date;
  const diff = d.getTime() - Date.now();
  const absSec = Math.abs(diff) / 1000;
  const future = diff > 0;

  const fmt = (n: number, unit: string) =>
    future ? `in ${n} ${unit}${n > 1 ? "s" : ""}` : `${n} ${unit}${n > 1 ? "s" : ""} ago`;

  if (absSec < 60) return "just now";
  if (absSec < 3600) return fmt(Math.floor(absSec / 60), "min");
  if (absSec < 86400) return fmt(Math.floor(absSec / 3600), "hour");
  if (absSec < 2592000) return fmt(Math.floor(absSec / 86400), "day");
  if (absSec < 31536000) return fmt(Math.floor(absSec / 2592000), "month");
  return fmt(Math.floor(absSec / 31536000), "year");
}

export function truncate(s: string, n: number): string {
  if (s.length <= n) return s;
  return s.slice(0, n).trimEnd() + "…";
}

export function initials(name: string): string {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase() ?? "")
    .join("");
}