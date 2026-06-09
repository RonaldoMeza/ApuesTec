"use client";

import Link from "next/link";
import type { LeaderboardEntry } from "@/features/leaderboard/types/leaderboard.types";

interface LeaderboardCardProps {
  entry: LeaderboardEntry;
  isCurrentUser?: boolean;
}

export function LeaderboardCard({ entry, isCurrentUser }: LeaderboardCardProps) {
  const rankMedal =
    entry.rank === 1 ? "🥇" : entry.rank === 2 ? "🥈" : entry.rank === 3 ? "🥉" : null;

  return (
    <Link href={`/leaderboard`}>
      <div
        className={`group rounded-xl border p-4 transition-all hover:shadow-lg ${
          isCurrentUser
            ? "border-primary/30 bg-primary/5"
            : "border-border bg-surface hover:border-primary/30 hover:shadow-primary/5"
        }`}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="text-lg font-bold text-muted-foreground">
              {rankMedal || `#${entry.rank}`}
            </span>
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-primary to-amber-500 text-sm font-bold text-black">
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
                  <span className="ml-1 text-xs text-primary">(tú)</span>
                )}
              </p>
              <p className="text-xs text-muted-foreground">
                {entry.predictionsCount} predicciones
              </p>
            </div>
          </div>
          <span className="text-xl font-bold text-primary">{entry.totalPoints}</span>
        </div>
      </div>
    </Link>
  );
}
