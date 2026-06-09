"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { useAuth } from "@/features/auth/context/AuthContext";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { AppLayout } from "@/shared/components/AppLayout";
import { matchService } from "@/features/matches/services/match.service";
import { MatchStatusBadge } from "@/features/matches/components/MatchStatusBadge";
import { predictionService } from "@/features/predictions/services/prediction.service";
import { leaderboardService } from "@/features/leaderboard/services/leaderboard.service";
import { ScoreEventList } from "@/features/leaderboard/components/ScoreEventList";
import type { MatchResponse } from "@/features/matches/types/match.types";
import type { PredictionResponse } from "@/features/predictions/types/prediction.types";
import type { UserStats, ScoreEvent, LeaderboardEntry } from "@/features/leaderboard/types/leaderboard.types";

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <DashboardContent />
    </ProtectedRoute>
  );
}

function DashboardContent() {
  const { user, isLoading } = useAuth();
  const [upcomingMatches, setUpcomingMatches] = useState<MatchResponse[]>([]);
  const [loadingMatches, setLoadingMatches] = useState(true);
  const [predictions, setPredictions] = useState<PredictionResponse[]>([]);
  const [loadingPredictions, setLoadingPredictions] = useState(true);
  const [stats, setStats] = useState<UserStats | null>(null);
  const [scoreEvents, setScoreEvents] = useState<ScoreEvent[]>([]);
  const [leaderboardPreview, setLeaderboardPreview] = useState<LeaderboardEntry[]>([]);
  const [loadingLeaderboard, setLoadingLeaderboard] = useState(true);

  useEffect(() => {
    matchService
      .listUpcoming()
      .then(setUpcomingMatches)
      .catch(() => setUpcomingMatches([]))
      .finally(() => setLoadingMatches(false));
  }, []);

  useEffect(() => {
    predictionService
      .listMyPredictions()
      .then(setPredictions)
      .catch(() => setPredictions([]))
      .finally(() => setLoadingPredictions(false));
  }, []);

  useEffect(() => {
    leaderboardService
      .getMyStats()
      .then(setStats)
      .catch(() => setStats(null));
  }, []);

  useEffect(() => {
    leaderboardService
      .getMyScoreEvents()
      .then(setScoreEvents)
      .catch(() => setScoreEvents([]));
  }, []);

  useEffect(() => {
    leaderboardService
      .getGlobalLeaderboard(5)
      .then((data) => setLeaderboardPreview(data.entries))
      .catch(() => setLeaderboardPreview([]))
      .finally(() => setLoadingLeaderboard(false));
  }, []);

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando dashboard..." />;
  }

  const isAdmin = user.roles.some((r) => r === "ADMIN" || r === "SUPER_ADMIN");

  const totalPoints = stats?.totalPoints ?? 0;
  const globalRank = stats?.globalRank ?? 0;
  const exactScores = stats?.exactScoresCount ?? 0;

  const summaryCards = [
    { label: "Puntos", value: totalPoints.toString(), desc: "Totales" },
    { label: "Predicciones", value: predictions.length.toString(), desc: "Registradas" },
    { label: "Ranking", value: globalRank > 0 ? `#${globalRank}` : "—", desc: "Global" },
    { label: "Exactos", value: exactScores.toString(), desc: "Marcadores exactos" },
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
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">Próximos partidos</h2>
              <Link href="/matches" className="text-xs text-primary hover:underline">
                Ver todos
              </Link>
            </div>
            {loadingMatches ? (
              <div className="flex items-center justify-center py-6">
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : upcomingMatches.length === 0 ? (
              <p className="text-sm text-muted-foreground">No hay partidos próximos.</p>
            ) : (
              <div className="space-y-3">
                {upcomingMatches.slice(0, 5).map((match) => (
                  <Link key={match.id} href={`/matches/${match.id}`}>
                    <div className="flex items-center justify-between rounded-lg bg-surface-muted px-4 py-3 transition-all hover:bg-surface-hover">
                      <div>
                        <p className="text-sm font-medium text-foreground">
                          {match.homeTeam?.name || "?"} vs {match.awayTeam?.name || "?"}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {new Date(match.startsAt).toLocaleDateString("es-PE", {
                            day: "numeric",
                            month: "short",
                            hour: "2-digit",
                            minute: "2-digit",
                          })}
                        </p>
                      </div>
                      <MatchStatusBadge status={match.status} />
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-xl border border-border bg-surface p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">Mis predicciones recientes</h2>
              <Link href="/predictions" className="text-xs text-primary hover:underline">
                Ver todas
              </Link>
            </div>
            {loadingPredictions ? (
              <div className="flex items-center justify-center py-6">
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : predictions.length === 0 ? (
              <div className="text-center py-6">
                <p className="text-sm text-muted-foreground mb-4">
                  No tienes predicciones registradas
                </p>
                <Link
                  href="/matches"
                  className="inline-block rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90"
                >
                  Ver partidos
                </Link>
              </div>
            ) : (
              <div className="space-y-3">
                {predictions.slice(0, 5).map((p) => (
                  <div key={p.id}>
                    <Link href={`/matches/${p.matchId}`}>
                      <div className="flex items-center justify-between rounded-lg bg-surface-muted px-4 py-3 transition-all hover:bg-surface-hover">
                        <div className="flex items-center gap-3">
                          <span className="text-lg font-bold text-primary">{p.homeScorePredicted}</span>
                          <span className="text-xs text-muted-foreground">-</span>
                          <span className="text-lg font-bold text-primary">{p.awayScorePredicted}</span>
                          {p.points > 0 && (
                            <span className="ml-2 text-xs font-medium text-emerald-400">
                              +{p.points}
                            </span>
                          )}
                        </div>
                        <span className="text-xs text-muted-foreground">
                          {p.points > 0 ? `${p.points} pts` : "Pendiente"}
                        </span>
                      </div>
                    </Link>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="mt-8 grid gap-6 lg:grid-cols-2">
          <div className="rounded-xl border border-border bg-surface p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">Ranking global</h2>
              <Link href="/leaderboard" className="text-xs text-primary hover:underline">
                Ver completo
              </Link>
            </div>
            {loadingLeaderboard ? (
              <div className="flex items-center justify-center py-6">
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : leaderboardPreview.length === 0 ? (
              <p className="py-6 text-center text-sm text-muted-foreground">
                El ranking aparecerá cuando los partidos finalicen.
              </p>
            ) : (
              <div className="space-y-2">
                {leaderboardPreview.map((entry) => (
                  <div
                    key={entry.userId}
                    className={`flex items-center justify-between rounded-lg px-4 py-2 ${
                      entry.userId === user.id
                        ? "bg-primary/5 ring-1 ring-primary/20"
                        : "bg-surface-muted"
                    }`}
                  >
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-bold text-muted-foreground">
                        #{entry.rank}
                      </span>
                      <span className="text-sm text-foreground">
                        {entry.fullName}
                        {entry.userId === user.id && (
                          <span className="ml-1 text-xs text-primary">(tú)</span>
                        )}
                      </span>
                    </div>
                    <span className="text-sm font-bold text-primary">
                      {entry.totalPoints}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-xl border border-border bg-surface p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">
                Eventos de puntuación recientes
              </h2>
              <Link href="/stats" className="text-xs text-primary hover:underline">
                Ver todas
              </Link>
            </div>
            <ScoreEventList events={scoreEvents.slice(0, 5)} />
          </div>
        </div>

        {isAdmin && (
          <div className="mt-6">
            <Link
              href="/admin"
              className="inline-flex items-center gap-1 text-sm text-primary hover:underline"
            >
              Panel de administración →
            </Link>
          </div>
        )}
      </div>
    </AppLayout>
  );
}
