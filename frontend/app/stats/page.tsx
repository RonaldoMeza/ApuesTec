"use client";

import { useState, useEffect } from "react";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { AppLayout } from "@/shared/components/AppLayout";
import { UserStatsCard } from "@/features/leaderboard/components/UserStatsCard";
import { ScoreEventList } from "@/features/leaderboard/components/ScoreEventList";
import { leaderboardService } from "@/features/leaderboard/services/leaderboard.service";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { useAuth } from "@/features/auth/context/AuthContext";
import type { UserStats, ScoreEvent } from "@/features/leaderboard/types/leaderboard.types";

export default function StatsPage() {
  return (
    <ProtectedRoute>
      <StatsContent />
    </ProtectedRoute>
  );
}

function StatsContent() {
  const { user, isLoading } = useAuth();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [events, setEvents] = useState<ScoreEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadStats() {
      try {
        const [statsData, eventsData] = await Promise.all([
          leaderboardService.getMyStats(),
          leaderboardService.getMyScoreEvents(),
        ]);
        if (!cancelled) {
          setStats(statsData);
          setEvents(eventsData);
          setLoading(false);
        }
      } catch {
        if (!cancelled) {
          setError("Error al cargar estadísticas");
          setLoading(false);
        }
      }
    }

    loadStats();
    return () => { cancelled = true; };
  }, []);

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando..." />;
  }

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Mis estadísticas</h1>
          <p className="mt-1 text-muted-foreground">
            Resumen de tu rendimiento en ApuesTec
          </p>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : error ? (
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
            <p className="text-red-400">{error}</p>
          </div>
        ) : (
          <>
            {stats && <UserStatsCard stats={stats} />}

            <div className="mt-8">
              <h2 className="mb-4 text-xl font-semibold text-foreground">
                Historial de puntuación
              </h2>
              <ScoreEventList events={events} />
            </div>
          </>
        )}
      </div>
    </AppLayout>
  );
}
