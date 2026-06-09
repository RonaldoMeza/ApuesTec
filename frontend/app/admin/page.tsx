"use client";

import { AdminRoute } from "@/features/auth/components/AdminRoute";
import { useAuth } from "@/features/auth/context/AuthContext";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { AppLayout } from "@/shared/components/AppLayout";
import Link from "next/link";

export default function AdminPage() {
  return (
    <AdminRoute>
      <AdminContent />
    </AdminRoute>
  );
}

function AdminContent() {
  const { user, isLoading } = useAuth();

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando..." />;
  }

  const isSuperAdmin = user.roles.includes("SUPER_ADMIN");

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Panel de Administración</h1>
          <p className="mt-1 text-muted-foreground">Gestiona equipos y partidos</p>
          {isSuperAdmin && (
            <p className="mt-2 text-xs text-primary/70">Acceso completo de SUPER_ADMIN</p>
          )}
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          <Link href="/admin/teams">
            <div className="group rounded-xl border border-border bg-surface p-6 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5">
              <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                <span className="text-xl">🏟️</span>
              </div>
              <h2 className="mb-1 text-lg font-semibold text-foreground">Gestión de Equipos</h2>
              <p className="text-sm text-muted-foreground">
                Crear, editar y eliminar equipos del Mundial.
              </p>
            </div>
          </Link>

          <Link href="/admin/matches">
            <div className="group rounded-xl border border-border bg-surface p-6 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5">
              <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                <span className="text-xl">⚽</span>
              </div>
              <h2 className="mb-1 text-lg font-semibold text-foreground">Gestión de Partidos</h2>
              <p className="text-sm text-muted-foreground">
                Crear partidos, cambiar estados y registrar resultados.
              </p>
            </div>
          </Link>
        </div>
      </div>
    </AppLayout>
  );
}
