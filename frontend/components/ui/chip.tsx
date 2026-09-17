"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface ChipProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  active?: boolean;
  emoji?: string;
  color?: string; // hex like "#3B82F6"
}

export const Chip = React.forwardRef<HTMLButtonElement, ChipProps>(
  ({ className, active = false, emoji, color, children, ...props }, ref) => {
    return (
      <button
        ref={ref}
        type="button"
        style={
          active && color
            ? {
                backgroundColor: `${color}1A`, // ~10% alpha
                color: color,
                borderColor: `${color}55`,
              }
            : undefined
        }
        className={cn(
          "inline-flex items-center gap-1.5 rounded-full border px-3.5 py-1.5 text-sm font-medium whitespace-nowrap transition-all duration-200 hover:scale-[1.03] active:scale-[0.97]",
          active
            ? "border-[var(--border-strong)] bg-[var(--bg-muted)] text-[var(--fg)] shadow-sm"
            : "border-[var(--border)] bg-[var(--bg-elevated)] text-[var(--fg-muted)] hover:border-[var(--border-strong)] hover:text-[var(--fg)]",
          className
        )}
        {...props}
      >
        {emoji && <span className="text-base leading-none">{emoji}</span>}
        {children}
      </button>
    );
  }
);
Chip.displayName = "Chip";