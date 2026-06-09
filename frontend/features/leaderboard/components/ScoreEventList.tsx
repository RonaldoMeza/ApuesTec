"use client";

import type { ScoreEvent } from "@/features/leaderboard/types/leaderboard.types";

interface ScoreEventListProps {
  events: ScoreEvent[];
}

const eventConfig: Record<string, { label: string; color: string }> = {
  EXACT_SCORE: { label: "Marcador exacto", color: "text-emerald-400" },
  WINNER_CORRECT: { label: "Ganador correcto", color: "text-blue-400" },
  GOAL_DIFFERENCE_CORRECT: { label: "Diferencia correcta", color: "text-purple-400" },
  EARLY_BONUS: { label: "Bonus anticipado", color: "text-amber-400" },
  STREAK_BONUS: { label: "Bonus racha", color: "text-rose-400" },
};

export function ScoreEventList({ events }: ScoreEventListProps) {
  if (events.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-surface p-8 text-center">
        <p className="text-sm text-muted-foreground">
          No tienes eventos de puntuación aún.
        </p>
        <p className="mt-1 text-xs text-muted-foreground">
          Los eventos aparecerán cuando los partidos que has predicho finalicen.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {events.map((event) => {
        const config = eventConfig[event.eventType] || {
          label: event.eventType,
          color: "text-muted-foreground",
        };
        return (
          <div
            key={event.id}
            className="flex items-center justify-between rounded-lg border border-border bg-surface px-4 py-3 transition-all hover:bg-surface-muted"
          >
            <div className="flex items-center gap-3">
              <span className={`text-sm font-medium ${config.color}`}>
                {config.label}
              </span>
              {event.description && (
                <span className="hidden text-xs text-muted-foreground sm:inline">
                  {event.description}
                </span>
              )}
            </div>
            <div className="flex items-center gap-3">
              <span className={`text-sm font-bold ${config.color}`}>
                +{event.points}
              </span>
              <span className="text-xs text-muted-foreground">
                {new Date(event.createdAt).toLocaleDateString("es-PE", {
                  day: "numeric",
                  month: "short",
                })}
              </span>
            </div>
          </div>
        );
      })}
    </div>
  );
}
