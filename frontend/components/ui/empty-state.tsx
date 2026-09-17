import * as React from "react";
import { cn } from "@/lib/utils";
import type { LucideIcon } from "lucide-react";

export interface EmptyStateProps {
  icon?: LucideIcon;
  title: string;
  description?: string;
  action?: React.ReactNode;
  className?: string;
}

export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
  className,
}: EmptyStateProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center text-center py-16 px-6",
        className
      )}
    >
      {Icon && (
        <div className="mb-5 grid h-16 w-16 place-items-center rounded-2xl bg-[var(--bg-muted)]">
          <Icon className="h-8 w-8 text-[var(--fg-subtle)]" />
        </div>
      )}
      <h3 className="text-display text-lg font-semibold text-[var(--fg)]">{title}</h3>
      {description && (
        <p className="mt-2 max-w-sm text-sm text-[var(--fg-muted)]">{description}</p>
      )}
      {action && <div className="mt-6">{action}</div>}
    </div>
  );
}