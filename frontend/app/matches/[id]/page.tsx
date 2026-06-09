"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { AppLayout } from "@/shared/components/AppLayout";
import { MatchStatusBadge } from "@/features/matches/components/MatchStatusBadge";
import { matchService } from "@/features/matches/services/match.service";
import type { MatchResponse } from "@/features/matches/types/match.types";

export default function MatchDetailPage() {
  const params = useParams();
  const id = params?.id as string;
  const [match, setMatch] = useState<MatchResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!id) return;

    let cancelled = false;
    matchService
      .getByID(id)
      .then((data) => {
        if (!cancelled) {
          setMatch(data);
          setLoading(false);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setError("Partido no encontrado");
          setLoading(false);
        }
      });

    return () => { cancelled = true; };
  }, [id]);

  if (loading) {
    return (
      <AppLayout>
        <div className="mx-auto max-w-4xl px-4 pt-8">
          <div className="flex items-center justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        </div>
      </AppLayout>
    );
  }

  if (error || !match) {
    return (
      <AppLayout>
        <div className="mx-auto max-w-4xl px-4 pt-8">
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
            <p className="text-red-400">{error || "Partido no encontrado"}</p>
            <Link
              href="/matches"
              className="mt-4 inline-block text-sm text-primary hover:underline"
            >
              Volver a partidos
            </Link>
          </div>
        </div>
      </AppLayout>
    );
  }

  const matchDate = new Date(match.startsAt);
  const formattedDate = matchDate.toLocaleDateString("es-PE", {
    weekday: "long",
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
    <AppLayout>
      <div className="mx-auto max-w-4xl px-4 pt-8">
        <Link
          href="/matches"
          className="mb-6 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          ← Volver a partidos
        </Link>

        <div className="rounded-xl border border-border bg-surface p-8">
          <div className="mb-6 flex items-center justify-between">
            <MatchStatusBadge status={match.status} />
            <span className="text-sm text-muted-foreground">{formattedDate}</span>
          </div>

          <div className="flex items-center justify-between gap-6">
            <div className="flex flex-1 flex-col items-center text-center">
              <div className="mb-2 flex h-20 w-20 items-center justify-center rounded-full bg-surface-muted">
                <span className="text-2xl font-bold text-foreground">
                  {match.homeTeam?.name?.charAt(0) || "?"}
                </span>
              </div>
              <h2 className="text-xl font-bold text-foreground">
                {match.homeTeam?.name || "Local"}
              </h2>
              <p className="text-sm text-muted-foreground">{match.homeTeam?.country || ""}</p>
            </div>

            <div className="text-center">
              {isFinished && match.homeScore != null && match.awayScore != null ? (
                <div className="flex items-center gap-4">
                  <span className="text-5xl font-bold text-primary">{match.homeScore}</span>
                  <span className="text-2xl text-muted-foreground">-</span>
                  <span className="text-5xl font-bold text-primary">{match.awayScore}</span>
                </div>
              ) : (
                <div className="text-center">
                  <span className="text-lg text-muted-foreground">vs</span>
                  <p className="mt-1 text-sm text-muted-foreground">{formattedTime}</p>
                </div>
              )}
            </div>

            <div className="flex flex-1 flex-col items-center text-center">
              <div className="mb-2 flex h-20 w-20 items-center justify-center rounded-full bg-surface-muted">
                <span className="text-2xl font-bold text-foreground">
                  {match.awayTeam?.name?.charAt(0) || "?"}
                </span>
              </div>
              <h2 className="text-xl font-bold text-foreground">
                {match.awayTeam?.name || "Visitante"}
              </h2>
              <p className="text-sm text-muted-foreground">{match.awayTeam?.country || ""}</p>
            </div>
          </div>

          <div className="mt-8 rounded-lg border border-border/50 bg-surface-muted p-4 text-center">
            <p className="text-sm italic text-muted-foreground">
              Las predicciones estarán disponibles en una próxima fase.
            </p>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
