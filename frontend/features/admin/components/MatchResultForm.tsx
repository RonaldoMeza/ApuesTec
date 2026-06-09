"use client";

import { useState, type FormEvent } from "react";
import { toast } from "sonner";
import { matchService } from "@/features/matches/services/match.service";
import type { MatchResponse } from "@/features/matches/types/match.types";

interface MatchResultFormProps {
  match: MatchResponse;
  onSuccess?: (match: MatchResponse) => void;
}

export function MatchResultForm({ match, onSuccess }: MatchResultFormProps) {
  const [homeScore, setHomeScore] = useState(match.homeScore?.toString() || "0");
  const [awayScore, setAwayScore] = useState(match.awayScore?.toString() || "0");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const result = await matchService.updateResult(match.id, {
        homeScore: parseInt(homeScore, 10),
        awayScore: parseInt(awayScore, 10),
      });
      toast.success("Resultado registrado exitosamente");
      onSuccess?.(result);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      setError(apiErr?.message || "Error al registrar resultado");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400">
          {error}
        </div>
      )}

      <p className="text-sm text-muted-foreground">
        Registrar resultado para:{" "}
        <span className="font-semibold text-foreground">
          {match.homeTeam?.name || "Local"} vs {match.awayTeam?.name || "Visitante"}
        </span>
      </p>

      <div className="flex items-center gap-4">
        <div className="flex-1">
          <label className="mb-1 block text-sm font-medium text-foreground">
            {match.homeTeam?.name || "Local"}
          </label>
          <input
            type="number"
            min="0"
            value={homeScore}
            onChange={(e) => setHomeScore(e.target.value)}
            required
            className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
        </div>

        <span className="mt-6 text-lg font-bold text-muted-foreground">-</span>

        <div className="flex-1">
          <label className="mb-1 block text-sm font-medium text-foreground">
            {match.awayTeam?.name || "Visitante"}
          </label>
          <input
            type="number"
            min="0"
            value={awayScore}
            onChange={(e) => setAwayScore(e.target.value)}
            required
            className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
        </div>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-lg bg-gradient-to-r from-emerald-500 to-emerald-600 px-4 py-2.5 text-sm font-semibold text-white transition-all hover:opacity-90 disabled:opacity-50"
      >
        {loading ? "Registrando..." : "Registrar resultado"}
      </button>
    </form>
  );
}
