"use client";

import type { LeaderboardEntry } from "@/features/leaderboard/types/leaderboard.types";

interface LeaderboardTableProps {
  entries: LeaderboardEntry[];
  currentUserId?: string;
}

export function LeaderboardTable({ entries, currentUserId }: LeaderboardTableProps) {
  if (entries.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-surface p-12 text-center">
        <p className="text-lg text-muted-foreground">
          No hay datos de ranking disponibles.
        </p>
        <p className="mt-2 text-sm text-muted-foreground">
          Los puntos aparecerán cuando los partidos finalicen y se calculen las puntuaciones.
        </p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto rounded-xl border border-border">
      <table className="w-full">
        <thead>
          <tr className="border-b border-border bg-surface-muted">
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
              #
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
              Usuario
            </th>
            <th className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground">
              Puntos
            </th>
            <th className="hidden px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground md:table-cell">
              Predicciones
            </th>
            <th className="hidden px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground md:table-cell">
              Exactos
            </th>
            <th className="hidden px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground lg:table-cell">
              Aciertos
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {entries.map((entry) => {
            const isCurrentUser = entry.userId === currentUserId;
            const rankColor =
              entry.rank === 1
                ? "text-yellow-400"
                : entry.rank === 2
                  ? "text-gray-300"
                  : entry.rank === 3
                    ? "text-amber-600"
                    : "text-muted-foreground";

            return (
              <tr
                key={entry.userId}
                className={`transition-colors hover:bg-surface-muted ${
                  isCurrentUser ? "bg-primary/5 ring-1 ring-primary/20" : ""
                }`}
              >
                <td className={`px-4 py-4 text-sm font-bold ${rankColor}`}>
                  #{entry.rank}
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2">
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-gradient-to-br from-primary to-amber-500 text-xs font-bold text-black">
                      {entry.fullName
                        .split(" ")
                        .map((n) => n[0])
                        .join("")
                        .toUpperCase()
                        .slice(0, 2)}
                    </div>
                    <div>
                      <p className="text-sm font-medium text-foreground">
                        {entry.fullName}
                        {isCurrentUser && (
                          <span className="ml-2 text-xs text-primary">(tú)</span>
                        )}
                      </p>
                    </div>
                  </div>
                </td>
                <td className="px-4 py-4 text-right">
                  <span className="text-lg font-bold text-primary">
                    {entry.totalPoints}
                  </span>
                </td>
                <td className="hidden px-4 py-4 text-right text-sm text-muted-foreground md:table-cell">
                  {entry.predictionsCount}
                </td>
                <td className="hidden px-4 py-4 text-right text-sm text-muted-foreground md:table-cell">
                  {entry.exactScores}
                </td>
                <td className="hidden px-4 py-4 text-right text-sm text-muted-foreground lg:table-cell">
                  {entry.winnerCorrect}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
