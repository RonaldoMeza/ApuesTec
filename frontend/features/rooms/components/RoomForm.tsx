"use client";

import { useState, type FormEvent } from "react";
import type { CreateRoomRequest } from "@/features/rooms/types/room.types";

interface RoomFormProps {
  onSubmit: (data: CreateRoomRequest) => Promise<void>;
  initialData?: { name: string; description?: string; visibility?: string; hasPassword?: boolean };
  isLoading?: boolean;
  disabled?: boolean;
}

export function RoomForm({ onSubmit, initialData, isLoading, disabled }: RoomFormProps) {
  const [name, setName] = useState(initialData?.name || "");
  const [description, setDescription] = useState(initialData?.description || "");
  const [visibility, setVisibility] = useState<"PUBLIC" | "PRIVATE">(
    (initialData?.visibility as "PUBLIC" | "PRIVATE") || "PRIVATE"
  );
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const isEditing = !!initialData;
  const isPublic = visibility === "PUBLIC";
  const isDisabled = disabled || isLoading;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (isDisabled) return;
    setError("");

    if (!name.trim()) {
      setError("El nombre de la sala es obligatorio");
      return;
    }

    if (isPublic && !password && !isEditing) {
      setError("Las salas públicas requieren una contraseña");
      return;
    }

    try {
      await onSubmit({
        name: name.trim(),
        description: description.trim() || undefined,
        visibility,
        password: isPublic ? password || undefined : undefined,
      });
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      setError(apiErr?.message || "Error al guardar la sala");
    }
  }

  const inputClass = `mt-1 block w-full rounded-lg border px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 ${
    isDisabled
      ? "border-border bg-surface-muted text-muted-foreground cursor-not-allowed"
      : "border-border bg-surface-muted text-foreground focus:border-primary focus:ring-primary"
  }`;

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="room-name" className="block text-sm font-medium text-foreground">
          Nombre de la sala
        </label>
        <input
          id="room-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className={inputClass}
          placeholder="Ej: Mundial 2026 - Amigos"
          disabled={isDisabled}
          required
          maxLength={100}
        />
      </div>

      <div>
        <label htmlFor="room-description" className="block text-sm font-medium text-foreground">
          Descripción <span className="text-muted-foreground">(opcional)</span>
        </label>
        <textarea
          id="room-description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className={`${inputClass} resize-none`}
          placeholder="Describe el propósito de la sala..."
          disabled={isDisabled}
          rows={3}
          maxLength={500}
        />
      </div>

      <div className="rounded-xl border border-border bg-surface-muted p-4">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-foreground">Sala pública</p>
            <p className="text-xs text-muted-foreground">
              {isPublic
                ? "Visible para usuarios en tu misma red. Requiere contraseña."
                : "Solo accesible mediante invitación."}
            </p>
          </div>
          <button
            type="button"
            onClick={() => !isDisabled && setVisibility(isPublic ? "PRIVATE" : "PUBLIC")}
            disabled={isDisabled}
            className={`relative inline-flex h-6 w-11 shrink-0 rounded-full border-2 border-transparent transition-colors ${
              isDisabled ? "cursor-not-allowed opacity-50" : "cursor-pointer"
            } ${isPublic ? "bg-primary" : "bg-muted-foreground/30"}`}
            role="switch"
            aria-checked={isPublic}
          >
            <span
              className={`pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow ring-0 transition-transform ${
                isPublic ? "translate-x-5" : "translate-x-0"
              }`}
            />
          </button>
        </div>

        {isPublic && (
          <div className="mt-4">
            <label htmlFor="room-password" className="block text-sm font-medium text-foreground">
              Contraseña de la sala
            </label>
            <input
              id="room-password"
              type="text"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={`mt-1 block w-full rounded-lg border px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 ${
                isDisabled
                  ? "border-border bg-surface-muted text-muted-foreground cursor-not-allowed"
                  : "border-border bg-surface text-foreground focus:border-primary focus:ring-primary"
              }`}
              placeholder={isEditing ? "Dejar vacío para mantener la actual" : "Ej: mundial2026"}
              disabled={isDisabled}
              maxLength={50}
            />
            <p className="mt-1 text-xs text-muted-foreground">
              Quien conozca la contraseña podrá unirse a la sala.
            </p>
          </div>
        )}
      </div>

      {error && (
        <div className="rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400" role="alert">
          {error}
        </div>
      )}

      {isDisabled && !isLoading && (
        <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 px-4 py-3 text-sm text-amber-400">
          La sala está cerrada. No se pueden realizar cambios.
        </div>
      )}

      {!isDisabled && (
        <button
          type="submit"
          disabled={isLoading}
          className="cursor-pointer w-full rounded-xl bg-gradient-to-r from-primary to-amber-500 px-4 py-2.5 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] hover:shadow-xl hover:shadow-primary/30 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading
            ? "Guardando..."
            : isEditing
              ? "Actualizar sala"
              : "Crear sala"}
        </button>
      )}
    </form>
  );
}
