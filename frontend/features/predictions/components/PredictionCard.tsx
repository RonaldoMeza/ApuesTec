"use client";

import Link from "next/link";
import { PredictionStatus } from "@/features/predictions/components/PredictionStatus";
import type { PredictionWithMatch } from "@/features/predictions/types/prediction.types";

interface PredictionCardProps {
  prediction: PredictionWithMatch;
}

export function PredictionCard({ prediction }: PredictionCardProps) {
  const match = prediction.match;
  const matchDate = match ? new Date(match.startsAt) : null;
  const formattedDate = matchDate
    ? matchDate.toLocaleDateString("es-PE", {
        day: "numeric",
        month: "short",
        year: "numeric",
      })
    : "";

  return (
    <Link href={`/matches/${prediction.matchId}`}>
      <div className="group rounded-xl border border-border bg-surface p-5 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5">
        <div className="mb-3 flex items-center justify-between">
          {match ? (
            <span className="text-xs text-muted-foreground">{formattedDate}</span>
          ) : (
            <span className="text-xs text-muted-foreground">Partido</span>
          )}
          <PredictionStatus
            canEdit={prediction.canEdit}
            isLocked={prediction.isLocked}
            matchStatus={match?.status || ""}
          />
        </div>

        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-1 flex-col items-center text-center">
            <span className="text-sm font-semibold text-foreground">
              {match?.homeTeam?.name || "Local"}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-2xl font-bold text-primary">{prediction.homeScorePredicted}</span>
            <span className="text-sm text-muted-foreground">-</span>
            <span className="text-2xl font-bold text-primary">{prediction.awayScorePredicted}</span>
          </div>

          <div className="flex flex-1 flex-col items-center text-center">
            <span className="text-sm font-semibold text-foreground">
              {match?.awayTeam?.name || "Visitante"}
            </span>
          </div>
        </div>

        <div className="mt-3 flex items-center justify-between border-t border-border/50 pt-3">
          <span className="text-xs text-muted-foreground">
            Puntos: {prediction.points}
          </span>
          {match?.homeScore != null && match?.awayScore != null && (
            <span className="text-xs text-muted-foreground">
              Real: {match.homeScore} - {match.awayScore}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}
