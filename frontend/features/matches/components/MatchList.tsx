"use client";

import { MatchCard } from "@/features/matches/components/MatchCard";
import type { MatchResponse } from "@/features/matches/types/match.types";

interface MatchListProps {
  matches: MatchResponse[];
  emptyMessage?: string;
}

export function MatchList({ matches, emptyMessage = "No hay partidos disponibles" }: MatchListProps) {
  if (matches.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-surface p-8 text-center">
        <p className="text-lg text-muted-foreground">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {matches.map((match) => (
        <MatchCard key={match.id} match={match} />
      ))}
    </div>
  );
}
