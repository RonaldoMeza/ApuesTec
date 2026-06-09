"use client";

import { useState, useEffect } from "react";
import { AdminRoute } from "@/features/auth/components/AdminRoute";
import { useAuth } from "@/features/auth/context/AuthContext";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { AppLayout } from "@/shared/components/AppLayout";
import { TeamForm } from "@/features/admin/components/TeamForm";
import { teamService } from "@/features/teams/services/team.service";
import type { TeamResponse } from "@/features/teams/types/team.types";
import Link from "next/link";

export default function AdminTeamsPage() {
  return (
    <AdminRoute>
      <AdminTeamsContent />
    </AdminRoute>
  );
}

function AdminTeamsContent() {
  const { user, isLoading } = useAuth();
  const [teams, setTeams] = useState<TeamResponse[]>([]);
  const [loadingTeams, setLoadingTeams] = useState(true);
  const [selectedTeam, setSelectedTeam] = useState<TeamResponse | null>(null);

  useEffect(() => {
    loadTeams();
  }, []);

  async function loadTeams() {
    setLoadingTeams(true);
    try {
      const data = await teamService.list();
      setTeams(data);
    } catch {
      // handled silently
    } finally {
      setLoadingTeams(false);
    }
  }

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando..." />;
  }

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="mb-6 flex items-center justify-between">
          <div>
            <Link href="/admin" className="text-sm text-muted-foreground hover:text-foreground">
              ← Admin
            </Link>
            <h1 className="mt-2 text-2xl font-bold text-foreground">Gestión de Equipos</h1>
          </div>
        </div>

        <div className="grid gap-8 lg:grid-cols-2">
          <div className="rounded-xl border border-border bg-surface p-6">
            <h2 className="mb-4 text-lg font-semibold text-foreground">
              {selectedTeam ? "Editar equipo" : "Crear nuevo equipo"}
            </h2>
            <TeamForm
              key={selectedTeam?.id || "new"}
              initialData={selectedTeam}
              onSuccess={() => {
                setSelectedTeam(null);
                loadTeams();
              }}
            />
            {selectedTeam && (
              <button
                onClick={() => setSelectedTeam(null)}
                className="mt-3 text-sm text-muted-foreground hover:text-foreground"
              >
                Cancelar edición
              </button>
            )}
          </div>

          <div>
            <h2 className="mb-4 text-lg font-semibold text-foreground">Equipos registrados</h2>
            {loadingTeams ? (
              <div className="flex items-center justify-center py-10">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : teams.length === 0 ? (
              <div className="rounded-xl border border-border bg-surface p-6 text-center">
                <p className="text-muted-foreground">No hay equipos registrados</p>
              </div>
            ) : (
              <div className="space-y-2">
                {teams.map((team) => (
                  <div
                    key={team.id}
                    className="flex items-center justify-between rounded-lg border border-border bg-surface p-4"
                  >
                    <div>
                      <p className="font-medium text-foreground">{team.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {team.code} · {team.country}
                      </p>
                    </div>
                    <button
                      onClick={() => setSelectedTeam(team)}
                      className="rounded-lg border border-border px-3 py-1.5 text-xs text-muted-foreground transition-all hover:border-primary/30 hover:text-foreground"
                    >
                      Editar
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
