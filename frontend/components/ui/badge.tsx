import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium whitespace-nowrap transition-colors",
  {
    variants: {
      variant: {
        default: "bg-[var(--bg-muted)] text-[var(--fg-muted)]",
        brand: "bg-brand-gradient text-white",
        success: "bg-[var(--color-success)]/15 text-[var(--color-success)]",
        warning: "bg-[var(--color-warning)]/15 text-[var(--color-warning)]",
        danger: "bg-[var(--color-danger)]/15 text-[var(--color-danger)]",
        outline: "border border-[var(--border)] text-[var(--fg-muted)]",
        // Category-specific (colored by --color-cat-* tokens)
        hackathon: "bg-[var(--color-cat-hackathon)]/15 text-[var(--color-cat-hackathon)]",
        cultural: "bg-[var(--color-cat-cultural)]/15 text-[var(--color-cat-cultural)]",
        sports: "bg-[var(--color-cat-sports)]/15 text-[var(--color-cat-sports)]",
        workshop: "bg-[var(--color-cat-workshop)]/15 text-[var(--color-cat-workshop)]",
        techfest: "bg-[var(--color-cat-techfest)]/15 text-[var(--color-cat-techfest)]",
        other: "bg-[var(--color-cat-other)]/15 text-[var(--color-cat-other)]",
      },
    },
    defaultVariants: { variant: "default" },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />;
}

export { badgeVariants };