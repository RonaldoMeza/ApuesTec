"use client";

import { useState, useEffect } from "react";
import { AdminRoute } from "@/features/auth/components/AdminRoute";
import { useAuth } from "@/features/auth/context/AuthContext";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { AppLayout } from "@/shared/components/AppLayout";
import { MatchForm } from "@/features/admin/components/MatchForm";
import { MatchResultForm } from "@/features/admin/components/MatchResultForm";
import { MatchStatusBadge } from "@/features/matches/components/MatchStatusBadge";
import { matchService } from "@/features/matches/services/match.service";
import type { MatchResponse } from "@/features/matches/types/match.types";
import Link from "next/link";

export default function AdminMatchesPage() {
  return (
    <AdminRoute>
      <AdminMatchesContent />
    </AdminRoute>
  );
}

function AdminMatchesContent() {
  const { user, isLoading } = useAuth();
  const [matches, setMatches] = useState<MatchResponse[]>([]);
  const [loadingMatches, setLoadingMatches] = useState(true);
  const [selectedResultMatch, setSelectedResultMatch] = useState<MatchResponse | null>(null);

  useEffect(() => {
    loadMatches();
  }, []);

  async function loadMatches() {
    setLoadingMatches(true);
    try {
      const data = await matchService.list();
      setMatches(data);
    } catch {
      // handled silently
    } finally {
      setLoadingMatches(false);
    }
  }

  function handleStatusChange(matchId: string, newStatus: string) {
    matchService.updateStatus(matchId, { status: newStatus }).then(() => loadMatches());
  }

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando..." />;
  }

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <Link href="/admin" className="text-sm text-muted-foreground hover:text-foreground">
          ← Admin
        </Link>
        <h1 className="mt-2 mb-6 text-2xl font-bold text-foreground">Gestión de Partidos</h1>

        <div className="grid gap-8 lg:grid-cols-2">
          <div className="rounded-xl border border-border bg-surface p-6">
            <h2 className="mb-4 text-lg font-semibold text-foreground">Crear nuevo partido</h2>
            <MatchForm onSuccess={() => loadMatches()} />
          </div>

          <div>
            <h2 className="mb-4 text-lg font-semibold text-foreground">Partidos registrados</h2>
            {loadingMatches ? (
              <div className="flex items-center justify-center py-10">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : matches.length === 0 ? (
              <div className="rounded-xl border border-border bg-surface p-6 text-center">
                <p className="text-muted-foreground">No hay partidos registrados</p>
              </div>
            ) : (
              <div className="space-y-3">
                {matches.map((match) => (
                  <div key={match.id} className="rounded-lg border border-border bg-surface p-4">
                    <div className="mb-2 flex items-center justify-between">
                      <MatchStatusBadge status={match.status} />
                      <span className="text-xs text-muted-foreground">
                        {new Date(match.startsAt).toLocaleDateString("es-PE")}
                      </span>
                    </div>
                    <p className="font-medium text-foreground">
                      {match.homeTeam?.name || "?"} vs {match.awayTeam?.name || "?"}
                    </p>
                    {match.status === "SCHEDULED" && (
                      <div className="mt-3 flex gap-2">
                        <button
                          onClick={() => handleStatusChange(match.id, "LOCKED")}
                          className="rounded-lg border border-amber-500/30 px-3 py-1 text-xs text-amber-400 transition-all hover:bg-amber-500/10"
                        >
                          Bloquear
                        </button>
                        <button
                          onClick={() => handleStatusChange(match.id, "CANCELLED")}
                          className="rounded-lg border border-red-500/30 px-3 py-1 text-xs text-red-400 transition-all hover:bg-red-500/10"
                        >
                          Cancelar
                        </button>
                      </div>
                    )}
                    {match.status === "LOCKED" && (
                      <div className="mt-3">
                        <button
                          onClick={() => setSelectedResultMatch(match)}
                          className="rounded-lg border border-emerald-500/30 px-3 py-1 text-xs text-emerald-400 transition-all hover:bg-emerald-500/10"
                        >
                          Registrar resultado
                        </button>
                      </div>
                    )}
                    {match.status === "FINISHED" && match.homeScore != null && match.awayScore != null && (
                      <p className="mt-2 text-sm font-bold text-primary">
                        {match.homeScore} - {match.awayScore}
                      </p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {selectedResultMatch && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
            <div className="w-full max-w-md rounded-xl border border-border bg-surface p-6 shadow-xl">
              <h2 className="mb-4 text-lg font-semibold text-foreground">Registrar resultado</h2>
              <MatchResultForm
                match={selectedResultMatch}
                onSuccess={() => {
                  setSelectedResultMatch(null);
                  loadMatches();
                }}
              />
              <button
                onClick={() => setSelectedResultMatch(null)}
                className="mt-3 text-sm text-muted-foreground hover:text-foreground"
              >
                Cancelar
              </button>
            </div>
          </div>
        )}
      </div>
    </AppLayout>
  );
}
