"use client";

import { Chip } from "@/components/ui/chip";
import { CATEGORIES } from "@/lib/constants";
import type { EventCategory } from "@/lib/api/types";

interface Props {
  value: string | "ALL";
  onChange: (value: EventCategory | "ALL") => void;
  className?: string;
}

export function CategoryChips({ value, onChange, className }: Props) {
  return (
    <div className={`flex gap-2 overflow-x-auto no-scrollbar pb-1 ${className ?? ""}`}>
      {CATEGORIES.map((cat) => (
        <Chip
          key={cat.value}
          active={value === cat.value}
          emoji={cat.emoji}
          color={cat.color}
          onClick={() => onChange(cat.value)}
        >
          {cat.label}
        </Chip>
      ))}
    </div>
  );
}