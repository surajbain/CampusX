"use client";

import { MessageCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface Props {
  link?: string;
  size?: "sm" | "md" | "lg";
  fullWidth?: boolean;
}

export function WhatsAppButton({ link, size = "md", fullWidth }: Props) {
  if (!link) return null;

  return (
    <a
      href={link}
      target="_blank"
      rel="noopener noreferrer"
      className={fullWidth ? "block" : "inline-block"}
    >
      <Button
        size={size}
        variant="secondary"
        className={`gap-2 ${fullWidth ? "w-full" : ""}`}
        style={{
          borderColor: "#25D366",
          color: "#25D366",
        }}
      >
        <MessageCircle className="h-4 w-4" />
        Join WhatsApp group
      </Button>
    </a>
  );
}