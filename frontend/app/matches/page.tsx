"use client";

import { useState, useEffect, useTransition } from "react";
import { AppLayout } from "@/shared/components/AppLayout";
import { MatchList } from "@/features/matches/components/MatchList";
import { matchService } from "@/features/matches/services/match.service";
import type { MatchResponse } from "@/features/matches/types/match.types";

type TabType = "upcoming" | "finished";

export default function MatchesPage() {
  const [tab, setTab] = useState<TabType>("upcoming");
  const [matches, setMatches] = useState<MatchResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [, startTransition] = useTransition();

  useEffect(() => {
    let cancelled = false;
    matchService[tab === "upcoming" ? "listUpcoming" : "listFinished"]()
      .then((data) => {
        if (!cancelled) {
          setMatches(data);
          setLoading(false);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setError("Error al cargar partidos");
          setLoading(false);
        }
      });
    return () => { cancelled = true; };
  }, [tab]);

  function handleTabChange(newTab: TabType) {
    startTransition(() => {
      setLoading(true);
      setError("");
      setTab(newTab);
    });
  }

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Partidos</h1>
          <p className="mt-1 text-muted-foreground">Todos los partidos del Mundial</p>
        </div>

        <div className="mb-6 flex gap-2">
          <button
            onClick={() => handleTabChange("upcoming")}
            className={`rounded-lg px-4 py-2 text-sm font-medium transition-all ${
              tab === "upcoming"
                ? "bg-primary text-primary-foreground"
                : "border border-border bg-surface text-muted-foreground hover:text-foreground"
            }`}
          >
            Próximos
          </button>
          <button
            onClick={() => handleTabChange("finished")}
            className={`rounded-lg px-4 py-2 text-sm font-medium transition-all ${
              tab === "finished"
                ? "bg-primary text-primary-foreground"
                : "border border-border bg-surface text-muted-foreground hover:text-foreground"
            }`}
          >
            Finalizados
          </button>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : error ? (
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
            <p className="text-red-400">{error}</p>
          </div>
        ) : (
          <MatchList
            matches={matches}
            emptyMessage={
              tab === "upcoming"
                ? "No hay partidos próximos. Vuelve más tarde."
                : "No hay partidos finalizados todavía."
            }
          />
        )}
      </div>
    </AppLayout>
  );
}
