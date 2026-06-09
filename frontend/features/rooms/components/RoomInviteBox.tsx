"use client";

import { useState } from "react";
import { toast } from "sonner";
import { roomService } from "@/features/rooms/services/room.service";
import type { RoomInvite } from "@/features/rooms/types/room.types";

interface RoomInviteBoxProps {
  roomId: string;
  myRole: string;
}

const DURATION_OPTIONS = [1, 3, 5, 10, 15, 20];

export function RoomInviteBox({ roomId, myRole }: RoomInviteBoxProps) {
  const [duration, setDuration] = useState(5);
  const [invite, setInvite] = useState<RoomInvite | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState("");

  const canInvite = myRole === "OWNER" || myRole === "MODERATOR";

  async function handleGenerate() {
    setError("");
    setLoading(true);
    setCopied(false);
    try {
      const result = await roomService.createInvite(roomId, duration);
      setInvite(result);
      toast.success("Invitación generada correctamente");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      const msg = apiErr?.message || "Error al generar la invitación";
      setError(msg);
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }

  async function copyText(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // fallback
    }
  }

  if (!canInvite) {
    return null;
  }

  const inviteLink = invite
    ? `${window.location.origin}/rooms/join?code=${invite.code}`
    : "";

  return (
    <div className="rounded-xl border border-border bg-surface p-5">
      <h4 className="mb-3 text-sm font-semibold text-foreground">Invitaciones</h4>

      <div className="mb-3 flex items-end gap-3">
        <div className="flex-1">
          <label htmlFor="invite-duration" className="block text-xs font-medium text-muted-foreground mb-1">
            Duración
          </label>
          <select
            id="invite-duration"
            value={duration}
            onChange={(e) => setDuration(Number(e.target.value))}
            className="block w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
            disabled={loading}
          >
            {DURATION_OPTIONS.map((min) => (
              <option key={min} value={min}>
                {min} minuto{min !== 1 ? "s" : ""}
              </option>
            ))}
          </select>
        </div>
        <button
          type="button"
          onClick={handleGenerate}
          disabled={loading}
          className="cursor-pointer rounded-xl bg-gradient-to-r from-primary to-amber-500 px-4 py-2 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] hover:shadow-xl hover:shadow-primary/30 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? "Generando..." : "Generar invitación"}
        </button>
      </div>

      {error && (
        <div className="mb-3 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400" role="alert">
          {error}
        </div>
      )}

      {invite && (
        <div className="space-y-3">
          <div>
            <p className="mb-1 text-xs text-muted-foreground">Código de invitación</p>
            <div className="flex items-center gap-2">
              <code className="flex-1 rounded-lg border border-primary/30 bg-surface-muted px-3 py-2 font-mono text-sm font-bold text-primary">
                {invite.code}
              </code>
              <button
                type="button"
                onClick={() => copyText(invite.code)}
                className="rounded-lg border border-border bg-surface-muted px-3 py-2 text-xs text-muted-foreground transition-colors hover:border-primary/30 hover:text-primary"
              >
                {copied ? "Copiado" : "Copiar"}
              </button>
            </div>
          </div>

          <div>
            <p className="mb-1 text-xs text-muted-foreground">Enlace de invitación</p>
            <div className="flex items-center gap-2">
              <code className="flex-1 truncate rounded-lg border border-border bg-surface-muted px-3 py-2 font-mono text-xs text-muted-foreground">
                {inviteLink}
              </code>
              <button
                type="button"
                onClick={() => copyText(inviteLink)}
                className="rounded-lg border border-border bg-surface-muted px-3 py-2 text-xs text-muted-foreground transition-colors hover:border-primary/30 hover:text-primary"
              >
                {copied ? "Copiado" : "Copiar"}
              </button>
            </div>
          </div>

          <div>
            <p className="mb-1 text-xs text-muted-foreground">QR Payload</p>
            <code className="block truncate rounded-lg border border-border bg-surface-muted px-3 py-2 font-mono text-xs text-muted-foreground">
              {invite.qrPayload}
            </code>
          </div>

          <p className="text-xs text-muted-foreground">
            Expira: {new Date(invite.expiresAt).toLocaleString("es-PE")}
          </p>
        </div>
      )}
    </div>
  );
}
