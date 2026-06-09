"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { AppLayout } from "@/shared/components/AppLayout";
import { predictionService } from "@/features/predictions/services/prediction.service";
import { PredictionCard } from "@/features/predictions/components/PredictionCard";
import { matchService } from "@/features/matches/services/match.service";
import type { PredictionWithMatch, PredictionResponse } from "@/features/predictions/types/prediction.types";
import type { MatchResponse } from "@/features/matches/types/match.types";

export default function PredictionsPage() {
  return (
    <ProtectedRoute>
      <PredictionsContent />
    </ProtectedRoute>
  );
}

function PredictionsContent() {
  const [predictions, setPredictions] = useState<PredictionWithMatch[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadPredictions() {
      try {
        const preds = await predictionService.listMyPredictions();

        const matchIds = [...new Set(preds.map((p) => p.matchId))];
        const matchMap = new Map<string, MatchResponse>();

        const matches = await matchService.list();
        for (const m of matches) {
          if (matchIds.includes(m.id)) {
            matchMap.set(m.id, m);
          }
        }

        const enriched: PredictionWithMatch[] = preds.map((p: PredictionResponse) => ({
          ...p,
          match: matchMap.get(p.matchId),
        }));

        if (!cancelled) {
          setPredictions(enriched);
          setLoading(false);
        }
      } catch {
        if (!cancelled) {
          setError("Error al cargar predicciones");
          setLoading(false);
        }
      }
    }

    loadPredictions();
    return () => { cancelled = true; };
  }, []);

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Mis predicciones</h1>
          <p className="mt-1 text-muted-foreground">
            Historial de todas tus predicciones
          </p>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : error ? (
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
            <p className="text-red-400">{error}</p>
          </div>
        ) : predictions.length === 0 ? (
          <div className="rounded-xl border border-border bg-surface p-12 text-center">
            <p className="text-lg text-muted-foreground mb-4">
              No tienes predicciones registradas
            </p>
            <Link
              href="/matches"
              className="inline-block rounded-lg bg-primary px-6 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90"
            >
              Ver partidos disponibles
            </Link>
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {predictions.map((prediction) => (
              <PredictionCard key={prediction.id} prediction={prediction} />
            ))}
          </div>
        )}
      </div>
    </AppLayout>
  );
}
