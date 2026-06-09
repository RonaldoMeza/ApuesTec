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

  const isFinished = match?.status === "FINISHED";
  const isCancelled = match?.status === "CANCELLED";
  const hasPoints = prediction.points > 0;

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
          <div className="flex flex-wrap gap-1">
            {isFinished && hasPoints && (
              <>
                {prediction.isExactScore && (
                  <span className="rounded-full bg-emerald-500/10 px-2 py-0.5 text-xs font-medium text-emerald-400">
                    Exacto
                  </span>
                )}
                {prediction.isWinnerCorrect && !prediction.isExactScore && (
                  <span className="rounded-full bg-blue-500/10 px-2 py-0.5 text-xs font-medium text-blue-400">
                    Ganador
                  </span>
                )}
                {prediction.isGoalDifferenceCorrect && !prediction.isExactScore && (
                  <span className="rounded-full bg-purple-500/10 px-2 py-0.5 text-xs font-medium text-purple-400">
                    Diferencia
                  </span>
                )}
                {prediction.earlyBonusPoints > 0 && (
                  <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-xs font-medium text-amber-400">
                    +Anticipado
                  </span>
                )}
                {prediction.streakBonusPoints > 0 && (
                  <span className="rounded-full bg-rose-500/10 px-2 py-0.5 text-xs font-medium text-rose-400">
                    +Racha
                  </span>
                )}
              </>
            )}
            {isCancelled && (
              <span className="text-xs text-muted-foreground">Sin puntos</span>
            )}
            {isFinished && !hasPoints && !isCancelled && (
              <span className="text-xs text-muted-foreground">Sin aciertos</span>
            )}
          </div>
          <span className="text-xs text-muted-foreground">
            {isFinished || isCancelled
              ? `Puntos: ${prediction.points}`
              : "Pendiente de puntuación"}
          </span>
        </div>

        {match?.homeScore != null && match?.awayScore != null && (
          <div className="mt-2 text-center text-xs text-muted-foreground">
            Real: {match.homeScore} - {match.awayScore}
          </div>
        )}
      </div>
    </Link>
  );
}
