"use client";

import Link from "next/link";
import type { InvitePreview } from "@/features/invites/types/invite.types";

interface InvitePreviewCardProps {
  invite: InvitePreview;
  onJoin?: () => void;
  isJoining?: boolean;
  isLoggedIn: boolean;
}

export function InvitePreviewCard({ invite, onJoin, isJoining, isLoggedIn }: InvitePreviewCardProps) {
  const expiresAt = new Date(invite.expiresAt);
  const now = new Date();
  const diff = expiresAt.getTime() - now.getTime();
  const minutesLeft = Math.max(0, Math.floor(diff / 60000));
  const secondsLeft = Math.max(0, Math.floor((diff % 60000) / 1000));

  return (
    <div className="rounded-xl border border-border bg-surface p-6 transition-all hover:border-primary/30">
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <h3 className="text-lg font-semibold text-foreground">{invite.roomName}</h3>
          {invite.roomDescription && (
            <p className="mt-1 text-sm text-muted-foreground">{invite.roomDescription}</p>
          )}
        </div>
        <span
          className={`shrink-0 rounded-full px-3 py-1 text-xs font-semibold ${
            invite.isExpired
              ? "bg-red-500/10 text-red-400"
              : "bg-emerald-500/10 text-emerald-400"
          }`}
        >
          {invite.isExpired ? "Expirada" : "Válida"}
        </span>
      </div>

      {!invite.isExpired && (
        <div className="mt-3">
          <p className="text-xs text-muted-foreground">
            Expira en {minutesLeft > 0 ? `${minutesLeft}m ${secondsLeft}s` : `${secondsLeft}s`}
          </p>
        </div>
      )}

      <div className="mt-5">
        {invite.isExpired ? (
          <div className="rounded-lg border border-red-500/20 bg-red-500/5 px-4 py-3 text-sm text-red-400">
            Esta invitación ha expirado
          </div>
        ) : isLoggedIn ? (
          <button
            onClick={onJoin}
            disabled={isJoining}
            className="cursor-pointer w-full rounded-xl bg-gradient-to-r from-primary to-amber-500 px-4 py-2.5 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] hover:shadow-xl hover:shadow-primary/30 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {isJoining ? "Uniéndote..." : "Unirme a la sala"}
          </button>
        ) : (
          <Link
            href="/login"
            className="block w-full rounded-xl border border-primary/30 bg-primary/5 px-4 py-2.5 text-center text-sm font-semibold text-primary transition-all hover:bg-primary/10"
          >
            Inicia sesión para unirte
          </Link>
        )}
      </div>
    </div>
  );
}
