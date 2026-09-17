import { TopNav } from "@/components/layout/top-nav";
import { AdminGuard } from "@/components/admin/admin-guard";
import { AdminSidebar, AdminBottomNav } from "@/components/admin/sidebar";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col">
      <TopNav />
      <AdminGuard>
        <div className="flex flex-1">
          <AdminSidebar />
          <main className="flex-1 pb-20 lg:pb-0">
            <div className="cx-container py-8">{children}</div>
          </main>
        </div>
        <AdminBottomNav />
      </AdminGuard>
    </div>
  );
}