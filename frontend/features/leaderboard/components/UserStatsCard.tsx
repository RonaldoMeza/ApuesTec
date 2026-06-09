"use client";

import type { UserStats } from "@/features/leaderboard/types/leaderboard.types";

interface UserStatsCardProps {
  stats: UserStats;
}

export function UserStatsCard({ stats }: UserStatsCardProps) {
  const items = [
    { label: "Puntos totales", value: stats.totalPoints, color: "text-primary" },
    { label: "Ranking global", value: `#${stats.globalRank || "-"}`, color: "text-yellow-400" },
    { label: "Predicciones", value: stats.predictionsCount, color: "text-foreground" },
    { label: "Marcadores exactos", value: stats.exactScoresCount, color: "text-emerald-400" },
    { label: "Ganador/empate correcto", value: stats.winnerCorrectCount, color: "text-blue-400" },
    { label: "Diferencia de goles", value: stats.goalDifferenceCorrectCount, color: "text-purple-400" },
    { label: "Racha actual", value: stats.currentStreak, color: "text-amber-400" },
    { label: "Mejor racha", value: stats.bestStreak, color: "text-rose-400" },
  ];

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {items.map((item, i) => (
        <div
          key={i}
          className="rounded-xl border border-border bg-surface p-5 transition-all hover:border-primary/30"
        >
          <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
            {item.label}
          </p>
          <p className={`mt-2 text-3xl font-bold ${item.color}`}>{item.value}</p>
        </div>
      ))}
    </div>
  );
}
