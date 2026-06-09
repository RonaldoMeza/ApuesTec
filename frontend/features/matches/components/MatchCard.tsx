"use client";

import Link from "next/link";
import { MatchStatusBadge } from "@/features/matches/components/MatchStatusBadge";
import type { MatchResponse } from "@/features/matches/types/match.types";

export function MatchCard({ match }: { match: MatchResponse }) {
  const matchDate = new Date(match.startsAt);
  const formattedDate = matchDate.toLocaleDateString("es-PE", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
  const formattedTime = matchDate.toLocaleTimeString("es-PE", {
    hour: "2-digit",
    minute: "2-digit",
  });

  const isFinished = match.status === "FINISHED" || match.status === "CANCELLED";

  return (
    <Link href={`/matches/${match.id}`}>
      <div className="group rounded-xl border border-border bg-surface p-5 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5">
        <div className="mb-3 flex items-center justify-between">
          <span className="text-xs text-muted-foreground">{formattedDate}</span>
          <MatchStatusBadge status={match.status} />
        </div>
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-1 flex-col items-center text-center">
            <span className="text-sm font-semibold text-foreground">
              {match.homeTeam?.name || "Local"}
            </span>
            <span className="mt-1 text-xs text-muted-foreground">{match.homeTeam?.country || ""}</span>
          </div>

          <div className="flex items-center gap-2">
            {isFinished && match.homeScore != null && match.awayScore != null ? (
              <>
                <span className="text-2xl font-bold text-primary">{match.homeScore}</span>
                <span className="text-sm text-muted-foreground">-</span>
                <span className="text-2xl font-bold text-primary">{match.awayScore}</span>
              </>
            ) : (
              <span className="text-xs text-muted-foreground">vs</span>
            )}
          </div>

          <div className="flex flex-1 flex-col items-center text-center">
            <span className="text-sm font-semibold text-foreground">
              {match.awayTeam?.name || "Visitante"}
            </span>
            <span className="mt-1 text-xs text-muted-foreground">{match.awayTeam?.country || ""}</span>
          </div>
        </div>
        <div className="mt-3 text-center">
          <span className="text-xs text-muted-foreground">{formattedTime}</span>
        </div>
      </div>
    </Link>
  );
}
