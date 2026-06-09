"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import { useAuth } from "@/features/auth/context/AuthContext";
import { AppLayout } from "@/shared/components/AppLayout";
import { inviteService } from "@/features/invites/services/invite.service";
import { InvitePreviewCard } from "@/features/invites/components/InvitePreviewCard";
import type { InvitePreview } from "@/features/invites/types/invite.types";

export default function JoinPage() {
  const params = useParams();
  const router = useRouter();
  const { user, isLoading: authLoading } = useAuth();
  const code = params.code as string;

  const [invite, setInvite] = useState<InvitePreview | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [joining, setJoining] = useState(false);
  const [joinError, setJoinError] = useState("");
  const [alreadyMemberRoomId, setAlreadyMemberRoomId] = useState("");

  useEffect(() => {
    if (authLoading) return;

    let cancelled = false;

    inviteService
      .preview(code)
      .then((data) => {
        if (!cancelled) setInvite(data);
      })
      .catch(() => {
        if (!cancelled) setError("Invitación no encontrada");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => { cancelled = true; };
  }, [code, authLoading]);

  async function handleJoin() {
    setJoinError("");
    setJoining(true);
    try {
      const result = await inviteService.join(code);
      toast.success(`Te uniste a "${result.roomName}"`);
      router.push(`/rooms/${result.roomId}`);
    } catch (err: unknown) {
      const apiErr = err as { message?: string; code?: string };
      if (apiErr?.code === "ALREADY_MEMBER" || apiErr?.message?.includes("ya eres miembro")) {
        setJoinError("Ya eres miembro de esta sala");
        setAlreadyMemberRoomId(invite?.roomName ? "redirect" : "");
        try {
          const data = await inviteService.join(code);
          router.push(`/rooms/${data.roomId}`);
        } catch {
          setJoinError("Ya eres miembro de esta sala");
        }
      } else {
        setJoinError(apiErr?.message || "Error al unirte a la sala");
      }
    } finally {
      setJoining(false);
    }
  }

  const isLoggedIn = !!user;

  return (
    <AppLayout>
      <div className="mx-auto max-w-lg px-4 pt-20">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-foreground">Invitación a sala</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Te han invitado a unirte a una sala privada
          </p>
        </div>

        {loading || authLoading ? (
          <div className="flex items-center justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : error ? (
          <div className="rounded-xl border border-border bg-surface p-12 text-center">
            <p className="text-lg text-muted-foreground mb-4">{error}</p>
            <Link
              href="/"
              className="inline-block text-sm text-primary hover:underline"
            >
              Volver al inicio
            </Link>
          </div>
        ) : invite ? (
          <div className="space-y-4">
            <InvitePreviewCard
              invite={invite}
              onJoin={handleJoin}
              isJoining={joining}
              isLoggedIn={isLoggedIn}
            />

            {joinError && !alreadyMemberRoomId && (
              <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-center">
                <p className="text-sm text-red-400">{joinError}</p>
              </div>
            )}

            {joinError && alreadyMemberRoomId && (
              <div className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-4 text-center">
                <p className="text-sm text-amber-400 mb-3">{joinError}</p>
                <Link
                  href="/rooms"
                  className="inline-block rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90"
                >
                  Ir a mis salas
                </Link>
              </div>
            )}

            {!isLoggedIn && !invite.isExpired && (
              <div className="rounded-xl border border-border bg-surface p-4 text-center">
                <p className="text-sm text-muted-foreground mb-3">
                  ¿Ya tienes cuenta?
                </p>
                <Link
                  href="/login"
                  className="inline-block rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90"
                >
                  Inicia sesión
                </Link>
              </div>
            )}
          </div>
        ) : null}
      </div>
    </AppLayout>
  );
}
