"use client";

import { useState, type FormEvent } from "react";
import { toast } from "sonner";
import { teamService } from "@/features/teams/services/team.service";
import type { TeamResponse } from "@/features/teams/types/team.types";

interface TeamFormProps {
  onSuccess?: (team: TeamResponse) => void;
  initialData?: TeamResponse | null;
}

export function TeamForm({ onSuccess, initialData }: TeamFormProps) {
  const [name, setName] = useState(initialData?.name || "");
  const [code, setCode] = useState(initialData?.code || "");
  const [country, setCountry] = useState(initialData?.country || "");
  const [flagUrl, setFlagUrl] = useState(initialData?.flagUrl || "");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const isEditing = !!initialData;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const payload = {
        name: name.trim(),
        code: code.trim().toUpperCase(),
        country: country.trim(),
        flagUrl: flagUrl.trim() || null,
      };

      let result: TeamResponse;
      if (isEditing) {
        result = await teamService.update(initialData!.id, payload);
        toast.success("Equipo actualizado exitosamente");
      } else {
        result = await teamService.create(payload);
        toast.success("Equipo creado exitosamente");
      }

      if (!isEditing) {
        setName("");
        setCode("");
        setCountry("");
        setFlagUrl("");
      }

      onSuccess?.(result);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      setError(apiErr?.message || "Error al guardar el equipo");
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

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">Nombre del equipo</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          minLength={1}
          maxLength={150}
          placeholder="Ej: Argentina"
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">Código</label>
        <input
          type="text"
          value={code}
          onChange={(e) => setCode(e.target.value)}
          required
          minLength={2}
          maxLength={10}
          placeholder="Ej: ARG"
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">País</label>
        <input
          type="text"
          value={country}
          onChange={(e) => setCountry(e.target.value)}
          required
          minLength={2}
          maxLength={100}
          placeholder="Ej: Argentina"
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-foreground">URL de bandera (opcional)</label>
        <input
          type="url"
          value={flagUrl}
          onChange={(e) => setFlagUrl(e.target.value)}
          placeholder="https://ejemplo.com/bandera.png"
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-lg bg-gradient-to-r from-primary to-amber-500 px-4 py-2.5 text-sm font-semibold text-primary-foreground transition-all hover:opacity-90 disabled:opacity-50"
      >
        {loading ? "Guardando..." : isEditing ? "Actualizar equipo" : "Crear equipo"}
      </button>
    </form>
  );
}
