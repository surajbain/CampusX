import type { EventCategory } from "./api/types";

export interface CategoryMeta {
  value: EventCategory | "ALL";
  label: string;
  emoji: string;
  color: string;
}

export const CATEGORIES: CategoryMeta[] = [
  { value: "ALL", label: "All Events", emoji: "✨", color: "#6366F1" },
  { value: "HACKATHON", label: "Hackathon", emoji: "💻", color: "#3B82F6" },
  { value: "CULTURAL", label: "Cultural", emoji: "🎭", color: "#EC4899" },
  { value: "SPORTS", label: "Sports", emoji: "⚽", color: "#F97316" },
  { value: "WORKSHOP", label: "Workshop", emoji: "🛠️", color: "#14B8A6" },
  { value: "TECH_FEST", label: "Tech Fest", emoji: "🚀", color: "#8B5CF6" },
  { value: "OTHER", label: "Other", emoji: "📌", color: "#6B7280" },
];

export function categoryMeta(cat?: string): CategoryMeta {
  return CATEGORIES.find((c) => c.value === cat) ?? CATEGORIES[CATEGORIES.length - 1];
}