"use client";

import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { useAuth } from "@/features/auth/context/AuthContext";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { AppLayout } from "@/shared/components/AppLayout";

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <DashboardContent />
    </ProtectedRoute>
  );
}

function DashboardContent() {
  const { user, isLoading } = useAuth();

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando dashboard..." />;
  }

  const summaryCards = [
    { label: "Puntos", value: "—", desc: "Próximamente" },
    { label: "Predicciones", value: "—", desc: "Próximamente" },
    { label: "Salas", value: "—", desc: "Próximamente" },
    { label: "Ranking", value: "—", desc: "Próximamente" },
  ];

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-foreground">
            Bienvenido, {user.fullName}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Panel principal de ApuesTec
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {summaryCards.map((card, i) => (
            <div key={i} className="rounded-xl border border-border bg-surface p-5 transition-all hover:border-primary/30">
              <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">{card.label}</p>
              <p className="mt-2 text-3xl font-bold text-primary">{card.value}</p>
              <p className="mt-1 text-xs text-muted-foreground">{card.desc}</p>
            </div>
          ))}
        </div>

        <div className="mt-8 grid gap-6 lg:grid-cols-2">
          <div className="rounded-xl border border-border bg-surface p-6">
            <h2 className="mb-1 text-lg font-semibold text-foreground">Información de la cuenta</h2>
            <p className="mb-4 text-sm text-muted-foreground">Datos asociados a tu perfil</p>
            <dl className="space-y-3 text-sm">
              <div className="flex justify-between rounded-lg bg-surface-muted px-4 py-3">
                <dt className="text-muted-foreground">Nombre</dt>
                <dd className="font-medium text-foreground">{user.fullName}</dd>
              </div>
              <div className="flex justify-between rounded-lg bg-surface-muted px-4 py-3">
                <dt className="text-muted-foreground">Email</dt>
                <dd className="font-medium text-foreground">{user.email}</dd>
              </div>
              <div className="flex justify-between rounded-lg bg-surface-muted px-4 py-3">
                <dt className="text-muted-foreground">Roles</dt>
                <dd className="font-medium text-foreground">{user.roles.join(", ")}</dd>
              </div>
            </dl>
          </div>

          <div className="rounded-xl border border-border bg-surface p-6">
            <h2 className="mb-1 text-lg font-semibold text-foreground">Próximamente</h2>
            <p className="mb-4 text-sm text-muted-foreground">Nuevas funcionalidades en desarrollo</p>
            <ul className="space-y-3 text-sm">
              <li className="flex items-center gap-3 rounded-lg bg-surface-muted px-4 py-3">
                <span className="text-primary">🏟️</span>
                <span className="text-muted-foreground">Salas privadas con amigos</span>
              </li>
              <li className="flex items-center gap-3 rounded-lg bg-surface-muted px-4 py-3">
                <span className="text-primary">⚽</span>
                <span className="text-muted-foreground">Predicciones de partidos</span>
              </li>
              <li className="flex items-center gap-3 rounded-lg bg-surface-muted px-4 py-3">
                <span className="text-primary">🏆</span>
                <span className="text-muted-foreground">Rankings y estadísticas</span>
              </li>
              <li className="flex items-center gap-3 rounded-lg bg-surface-muted px-4 py-3">
                <span className="text-primary">🎮</span>
                <span className="text-muted-foreground">Logros y gamificación</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
