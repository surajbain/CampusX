"use client";

import * as React from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Plus, Trash2, Loader2, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { UserForm } from "@/components/admin/user-form";
import {
  useCollegeUsers,
  useDeleteCollegeUser,
} from "@/lib/hooks/use-admin";
import { useCollege } from "@/lib/hooks/use-admin";

export default function CollegeUsersPage() {
  const params = useParams();
  const collegeId = params.id as string;
  const [showCreate, setShowCreate] = React.useState(false);

  const { data: college } = useCollege(collegeId);
  const { data, isLoading } = useCollegeUsers(collegeId, { limit: 100 });
  const deleteMutation = useDeleteCollegeUser();

  const users = ((data as any)?.data ?? []) as any[];

  return (
    <div>
      <Link
        href={`/admin/colleges/${collegeId}`}
        className="inline-flex items-center gap-2 text-sm text-[var(--fg-muted)] hover:text-[var(--fg)] mb-6 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to {college?.name || "college"}
      </Link>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-display text-2xl font-bold tracking-tight">
            Users
          </h1>
          <p className="mt-1 text-sm text-[var(--fg-muted)]">
            {users.length} users in this college
          </p>
        </div>
        <Button onClick={() => setShowCreate((v) => !v)}>
          <Plus className="h-4 w-4" />
          Add user
        </Button>
      </div>

      {showCreate && (
        <div className="mb-6">
          <UserForm
            collegeId={collegeId}
            onClose={() => setShowCreate(false)}
          />
        </div>
      )}

      {isLoading && (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-16 rounded-xl" />
          ))}
        </div>
      )}

      {!isLoading && users.length === 0 && (
        <EmptyState
          icon={Users}
          title="No users yet"
          description="Add organizers and volunteers to help run events."
        />
      )}

      {!isLoading && users.length > 0 && (
        <div className="space-y-2">
          {users.map((u) => (
            <Card key={u.id} className="p-4">
              <div className="flex items-center gap-3">
                <div className="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-brand-gradient text-white text-xs font-bold">
                  {u.full_name
                    .split(" ")
                    .map((w: string) => w[0])
                    .slice(0, 2)
                    .join("")
                    .toUpperCase()}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium truncate">{u.full_name}</span>
                    <Badge
                      variant={
                        u.role === "COLLEGE_ADMIN"
                          ? "brand"
                          : u.role === "ORGANIZER"
                          ? "hackathon"
                          : "outline"
                      }
                    >
                      {u.role}
                    </Badge>
                  </div>
                  <div className="text-xs text-[var(--fg-muted)] truncate">
                    {u.email}
                  </div>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    if (confirm(`Delete ${u.full_name}?`)) {
                      deleteMutation.mutate({ collegeId, userId: u.id });
                    }
                  }}
                >
                  <Trash2 className="h-3.5 w-3.5 text-[var(--color-danger)]" />
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}