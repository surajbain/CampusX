import { TopNav } from "@/components/layout/top-nav";
import { ProtectedRoute } from "@/components/auth/protected-route";
import {
  StudentSidebar,
  StudentBottomNav,
} from "@/components/student/student-sidebar";

export default function StudentLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />
      <ProtectedRoute>
        <div className="flex flex-1">
          <StudentSidebar />
          <main className="flex-1 pb-20 lg:pb-0">
            <div className="cx-container py-8">{children}</div>
          </main>
        </div>
        <StudentBottomNav />
      </ProtectedRoute>
    </div>
  );
}