"use client";

import { useState, useEffect, type FormEvent } from "react";
import { useAuth } from "@/features/auth/context/AuthContext";
import { predictionService } from "@/features/predictions/services/prediction.service";
import { PredictionStatus } from "@/features/predictions/components/PredictionStatus";
import { toast } from "sonner";
import Link from "next/link";
import type { PredictionResponse } from "@/features/predictions/types/prediction.types";

interface PredictionFormProps {
  matchId: string;
  matchStatus: string;
  matchDate: Date;
}

export function PredictionForm({ matchId, matchStatus, matchDate }: PredictionFormProps) {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return (
      <div className="mt-8 rounded-xl border border-border bg-surface p-6 text-center">
        <p className="text-sm text-muted-foreground">
          Inicia sesión para registrar tu predicción
        </p>
        <Link
          href="/login"
          className="mt-3 inline-block rounded-lg bg-primary px-6 py-2 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90"
        >
          Iniciar sesión
        </Link>
      </div>
    );
  }

  return <PredictionFormAuthenticated matchId={matchId} matchStatus={matchStatus} matchDate={matchDate} />;
}

function PredictionFormAuthenticated({ matchId, matchStatus, matchDate }: PredictionFormProps) {
  const [prediction, setPrediction] = useState<PredictionResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [homeScore, setHomeScore] = useState("");
  const [awayScore, setAwayScore] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const isFinished = matchStatus === "FINISHED" || matchStatus === "CANCELLED";
  const isLocked = matchStatus === "LOCKED";
  const isTooClose = new Date().getTime() > matchDate.getTime() - 10 * 60 * 1000;

  const canEditPrediction = !isFinished && !isLocked && !isTooClose;

  useEffect(() => {
    let cancelled = false;
    predictionService
      .getMyPrediction(matchId)
      .then((data) => {
        if (!cancelled) {
          setPrediction(data);
          setHomeScore(data.homeScorePredicted.toString());
          setAwayScore(data.awayScorePredicted.toString());
          setLoading(false);
        }
      })
      .catch(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [matchId]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();

    const home = parseInt(homeScore, 10);
    const away = parseInt(awayScore, 10);

    if (isNaN(home) || isNaN(away) || home < 0 || away < 0) {
      toast.error("Los marcadores deben ser números enteros mayores o iguales a 0");
      return;
    }

    setSubmitting(true);
    try {
      const result = await predictionService.upsertPrediction(matchId, {
        homeScorePredicted: home,
        awayScorePredicted: away,
      });
      setPrediction(result);
      toast.success("Predicción guardada exitosamente");
    } catch {
      toast.error("Error al guardar la predicción");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <div className="mt-8 rounded-xl border border-border bg-surface p-6">
        <div className="flex items-center justify-center py-4">
          <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        </div>
      </div>
    );
  }

  return (
    <div className="mt-8 rounded-xl border border-border bg-surface p-6">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-lg font-semibold text-foreground">
          {prediction ? "Tu predicción" : "Registrar predicción"}
        </h3>
        {prediction && (
          <PredictionStatus
            canEdit={prediction.canEdit}
            isLocked={prediction.isLocked}
            matchStatus={matchStatus}
          />
        )}
      </div>

      {isFinished && (
        <div className="mb-4 rounded-lg bg-amber-500/10 p-3 text-sm text-amber-400">
          Este partido ya ha finalizado. No se pueden registrar predicciones.
        </div>
      )}

      {isLocked && !isFinished && (
        <div className="mb-4 rounded-lg bg-amber-500/10 p-3 text-sm text-amber-400">
          Este partido está bloqueado. No se pueden registrar predicciones.
        </div>
      )}

      {!isFinished && !isLocked && isTooClose && (
        <div className="mb-4 rounded-lg bg-red-500/10 p-3 text-sm text-red-400">
          Faltan menos de 10 minutos para el inicio del partido. Ya no se pueden registrar predicciones.
        </div>
      )}

      {prediction && !canEditPrediction && !isFinished && !isLocked && !isTooClose && prediction.isLocked && (
        <div className="mb-4 rounded-lg bg-amber-500/10 p-3 text-sm text-amber-400">
          Tu predicción ha sido bloqueada y no se puede editar.
        </div>
      )}

      {!canEditPrediction && prediction ? (
        <div className="flex items-center justify-center gap-4 py-4">
          <div className="text-center">
            <p className="text-sm text-muted-foreground">Local</p>
            <p className="text-4xl font-bold text-primary">{prediction.homeScorePredicted}</p>
          </div>
          <span className="text-2xl text-muted-foreground">-</span>
          <div className="text-center">
            <p className="text-sm text-muted-foreground">Visitante</p>
            <p className="text-4xl font-bold text-primary">{prediction.awayScorePredicted}</p>
          </div>
        </div>
      ) : null}

      {(canEditPrediction || (!prediction && canEditPrediction)) && (
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="flex items-center justify-center gap-4">
            <div className="w-24 text-center">
              <label className="block text-sm text-muted-foreground">Local</label>
              <input
                type="number"
                min="0"
                value={homeScore}
                onChange={(e) => setHomeScore(e.target.value)}
                className="mt-1 w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-center text-2xl font-bold text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                disabled={submitting}
                required
              />
            </div>
            <span className="mt-6 text-lg text-muted-foreground">-</span>
            <div className="w-24 text-center">
              <label className="block text-sm text-muted-foreground">Visitante</label>
              <input
                type="number"
                min="0"
                value={awayScore}
                onChange={(e) => setAwayScore(e.target.value)}
                className="mt-1 w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-center text-2xl font-bold text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                disabled={submitting}
                required
              />
            </div>
          </div>

          {prediction && (
            <p className="text-center text-xs text-muted-foreground">
              Estás editando tu predicción existente.
            </p>
          )}

          <button
            type="submit"
            disabled={submitting}
            className="w-full rounded-xl bg-gradient-to-r from-primary to-amber-500 px-4 py-2.5 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] hover:shadow-xl hover:shadow-primary/30 disabled:opacity-50"
          >
            {submitting
              ? "Guardando..."
              : prediction
                ? "Actualizar predicción"
                : "Registrar predicción"}
          </button>
        </form>
      )}

      <p className="mt-4 text-center text-xs text-muted-foreground">
        La puntuación se calculará cuando el partido finalice.
      </p>
    </div>
  );
}
