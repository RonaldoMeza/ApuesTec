"use client";

import { useState, useEffect, type FormEvent } from "react";
import { toast } from "sonner";
import { teamService } from "@/features/teams/services/team.service";
import { matchService } from "@/features/matches/services/match.service";
import type { TeamResponse } from "@/features/teams/types/team.types";
import type { MatchResponse } from "@/features/matches/types/match.types";

interface MatchFormProps {
  onSuccess?: (match: MatchResponse) => void;
}

export function MatchForm({ onSuccess }: MatchFormProps) {
  const [teams, setTeams] = useState<TeamResponse[]>([]);
  const [homeTeamId, setHomeTeamId] = useState("");
  const [awayTeamId, setAwayTeamId] = useState("");
  const [startsAt, setStartsAt] = useState("");
  const [loading, setLoading] = useState(false);
  const [loadingTeams, setLoadingTeams] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    teamService
      .list()
      .then(setTeams)
      .catch(() => setError("Error al cargar equipos"))
      .finally(() => setLoadingTeams(false));
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    if (homeTeamId === awayTeamId) {
      setError("El equipo local y visitante deben ser diferentes");
      setLoading(false);
      return;
    }

    try {
      const startsAtISO = new Date(startsAt).toISOString();
      const result = await matchService.create({
        homeTeamId,
        awayTeamId,
        startsAt: startsAtISO,
      });
      toast.success("Partido creado exitosamente");
      setHomeTeamId("");
      setAwayTeamId("");
      setStartsAt("");
      onSuccess?.(result);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      setError(apiErr?.message || "Error al crear el partido");
    } finally {
      setLoading(false);
    }
  }

  if (loadingTeams) {
    return <p className="text-sm text-muted-foreground">Cargando equipos...</p>;
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400">
          {error}
        </div>
      )}

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">Equipo local</label>
        <select
          value={homeTeamId}
          onChange={(e) => setHomeTeamId(e.target.value)}
          required
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        >
          <option value="">Seleccionar equipo local</option>
          {teams.map((team) => (
            <option key={team.id} value={team.id}>
              {team.name} ({team.code})
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">Equipo visitante</label>
        <select
          value={awayTeamId}
          onChange={(e) => setAwayTeamId(e.target.value)}
          required
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        >
          <option value="">Seleccionar equipo visitante</option>
          {teams.map((team) => (
            <option key={team.id} value={team.id}>
              {team.name} ({team.code})
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">Fecha y hora del partido</label>
        <input
          type="datetime-local"
          value={startsAt}
          onChange={(e) => setStartsAt(e.target.value)}
          required
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-lg bg-gradient-to-r from-primary to-amber-500 px-4 py-2.5 text-sm font-semibold text-primary-foreground transition-all hover:opacity-90 disabled:opacity-50"
      >
        {loading ? "Creando partido..." : "Crear partido"}
      </button>
    </form>
  );
}
