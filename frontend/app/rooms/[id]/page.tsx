"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { AppLayout } from "@/shared/components/AppLayout";
import { useAuth } from "@/features/auth/context/AuthContext";
import { roomService } from "@/features/rooms/services/room.service";
import { RoomRoleBadge } from "@/features/rooms/components/RoomRoleBadge";
import { RoomMemberList } from "@/features/rooms/components/RoomMemberList";
import { RoomInviteBox } from "@/features/rooms/components/RoomInviteBox";
import { RoomForm } from "@/features/rooms/components/RoomForm";
import type {
  Room,
  RoomMember,
  RoomLeaderboardEntry,
  UpdateRoomRequest,
} from "@/features/rooms/types/room.types";

type Tab = "ranking" | "members" | "invite" | "settings";

export default function RoomDetailPage() {
  return (
    <ProtectedRoute>
      <RoomDetailContent />
    </ProtectedRoute>
  );
}

function RoomDetailContent() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuth();
  const roomId = params.id as string;

  const [room, setRoom] = useState<Room | null>(null);
  const [members, setMembers] = useState<RoomMember[]>([]);
  const [leaderboard, setLeaderboard] = useState<RoomLeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [activeTab, setActiveTab] = useState<Tab>("ranking");
  const [leaving, setLeaving] = useState(false);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [roomData, membersData, leaderboardData] = await Promise.all([
          roomService.getById(roomId),
          roomService.getMembers(roomId),
          roomService.getLeaderboard(roomId),
        ]);
        if (cancelled) return;
        setRoom(roomData);
        setMembers(membersData.members);
        setLeaderboard(leaderboardData.entries);
      } catch {
        if (!cancelled) setError("Error al cargar la sala");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [roomId]);

  async function handleUpdate(data: UpdateRoomRequest) {
    if (!room) return;
    try {
      const updated = await roomService.update(roomId, data);
      setRoom(updated);
      toast.success("Sala actualizada correctamente");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al actualizar la sala");
    }
  }

  async function handleClose() {
    if (!confirm("¿Estás seguro de cerrar esta sala? Esta acción no se puede deshacer.")) return;
    try {
      const updated = await roomService.close(roomId);
      setRoom(updated);
      toast.success("Sala cerrada correctamente");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al cerrar la sala");
    }
  }

  async function handleLeave() {
    if (!confirm("¿Estás seguro de abandonar la sala?")) return;
    setLeaving(true);
    try {
      await roomService.leave(roomId);
      toast.success("Abandonaste la sala");
      router.push("/rooms");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al abandonar la sala");
      setLeaving(false);
    }
  }

  async function handleRoleChange(userId: string, role: string) {
    try {
      await roomService.changeRole(roomId, userId, role);
      setMembers((prev) =>
        prev.map((m) => (m.userId === userId ? { ...m, role: role as RoomMember["role"] } : m))
      );
      toast.success("Rol actualizado correctamente");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al cambiar el rol");
    }
  }

  async function handleRemoveMember(userId: string) {
    try {
      await roomService.removeMember(roomId, userId);
      setMembers((prev) => prev.filter((m) => m.userId !== userId));
      toast.success("Miembro expulsado de la sala");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al expulsar miembro");
    }
  }

  const isOwner = room?.myRole === "OWNER";
  const isModerator = room?.myRole === "MODERATOR";
  const canManage = isOwner || isModerator;
  const isClosed = room?.status === "CLOSED";

  const tabs: { key: Tab; label: string }[] = [
    { key: "ranking", label: "Ranking" },
    { key: "members", label: "Miembros" },
  ];

  if (canManage && !isClosed) {
    tabs.push({ key: "invite", label: "Invitación" });
  }

  if (isOwner) {
    tabs.push({ key: "settings", label: "Configuración" });
  }

  const statusBadge =
    room?.status === "ACTIVE"
      ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/30"
      : "bg-red-500/10 text-red-400 border-red-500/30";

  const visibilityBadge =
    room?.visibility === "PUBLIC"
      ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/30"
      : "bg-amber-500/10 text-amber-400 border-amber-500/30";

  return (
    <AppLayout>
      <div className="mx-auto max-w-4xl px-4 pt-8">
        <div className="mb-6">
          <Link
            href="/rooms"
            className="text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            &larr; Volver a mis salas
          </Link>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : error ? (
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
            <p className="text-red-400">{error}</p>
            <Link
              href="/rooms"
              className="mt-4 inline-block text-sm text-primary hover:underline"
            >
              Volver a mis salas
            </Link>
          </div>
        ) : room ? (
          <>
            <div className="mb-6">
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <div className="mb-2 flex items-center gap-2 flex-wrap">
                    <span
                      className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${statusBadge}`}
                    >
                      {room.status === "ACTIVE" ? "Activa" : "Cerrada"}
                    </span>
                    <span
                      className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${visibilityBadge}`}
                    >
                      {room.visibility === "PUBLIC" ? "Pública" : "Privada"}
                    </span>
                    <RoomRoleBadge role={room.myRole} />
                  </div>
                  <h1 className="text-3xl font-bold text-foreground">{room.name}</h1>
                  {room.description && (
                    <p className="mt-2 text-muted-foreground">{room.description}</p>
                  )}
                  <p className="mt-2 text-sm text-muted-foreground">
                    {room.memberCount} miembro{room.memberCount !== 1 ? "s" : ""}
                    {room.hasPassword && room.visibility === "PUBLIC" && (
                      <span className="ml-3 text-xs text-amber-400">· Con contraseña</span>
                    )}
                  </p>
                </div>

                {!isOwner && !isClosed && (
                  <button
                    onClick={handleLeave}
                    disabled={leaving}
                    className="cursor-pointer shrink-0 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-2 text-sm font-medium text-red-400 transition-colors hover:bg-red-500/20 disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {leaving ? "Abandonando..." : "Abandonar sala"}
                  </button>
                )}

                {isOwner && (
                  <div className="shrink-0 rounded-lg border border-border bg-surface-muted px-4 py-2 text-xs text-muted-foreground">
                    Eres el propietario
                  </div>
                )}
              </div>
            </div>

            <div className="mb-6 border-b border-border">
              <nav className="flex gap-6" role="tablist">
                {tabs.map((tab) => (
                  <button
                    key={tab.key}
                    role="tab"
                    aria-selected={activeTab === tab.key}
                    onClick={() => setActiveTab(tab.key)}
                    className={`cursor-pointer pb-3 text-sm font-medium transition-colors border-b-2 ${
                      activeTab === tab.key
                        ? "border-primary text-foreground"
                        : "border-transparent text-muted-foreground hover:text-foreground"
                    }`}
                  >
                    {tab.label}
                  </button>
                ))}
              </nav>
            </div>

            {activeTab === "ranking" && (
              <div className="space-y-3">
                {leaderboard.length === 0 ? (
                  <div className="rounded-xl border border-border bg-surface p-12 text-center">
                    <p className="text-sm text-muted-foreground">
                      No hay datos de ranking disponibles para esta sala.
                    </p>
                  </div>
                ) : (
                  leaderboard.map((entry) => {
                    const isCurrentUser = entry.userId === user?.id;
                    return (
                      <div
                        key={entry.userId}
                        className={`flex items-center justify-between rounded-xl border px-5 py-4 transition-all ${
                          isCurrentUser
                            ? "border-primary/30 bg-primary/5"
                            : "border-border bg-surface hover:border-primary/20"
                        }`}
                      >
                        <div className="flex items-center gap-4">
                          <span className="text-lg font-bold text-muted-foreground min-w-[2rem]">
                            #{entry.rank}
                          </span>
                          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-primary to-amber-500 text-sm font-bold text-black">
                            {entry.fullName
                              .split(" ")
                              .map((n) => n[0])
                              .join("")
                              .toUpperCase()
                              .slice(0, 2)}
                          </div>
                          <div>
                            <p className="text-sm font-medium text-foreground">
                              {entry.fullName}
                              {isCurrentUser && (
                                <span className="ml-2 text-xs text-primary">(tú)</span>
                              )}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              {entry.predictionsCount} predicciones &middot;{" "}
                              {entry.exactScoresCount} exactos
                            </p>
                          </div>
                        </div>
                        <span className="text-xl font-bold text-primary">
                          {entry.totalPoints}
                        </span>
                      </div>
                    );
                  })
                )}

                {isClosed && (
                  <div className="rounded-xl border border-amber-500/30 bg-amber-500/5 p-4 text-center">
                    <p className="text-sm text-amber-400">
                      Esta sala está cerrada. No se pueden realizar más cambios.
                    </p>
                  </div>
                )}
              </div>
            )}

            {activeTab === "members" && (
              <RoomMemberList
                members={members}
                currentUserId={user?.id || ""}
                myRole={room.myRole}
                onRoleChange={isOwner && !isClosed ? handleRoleChange : undefined}
                onRemoveMember={canManage && !isClosed ? handleRemoveMember : undefined}
              />
            )}

            {activeTab === "invite" && canManage && !isClosed && (
              <RoomInviteBox roomId={roomId} myRole={room.myRole} />
            )}

            {activeTab === "settings" && isOwner && (
              <div className="space-y-6">
                <div className="rounded-xl border border-border bg-surface p-5">
                  <h4 className="mb-4 text-sm font-semibold text-foreground">
                    {isClosed ? "Configuración (sala cerrada)" : "Editar sala"}
                  </h4>
                  <RoomForm
                    onSubmit={handleUpdate}
                    initialData={{
                      name: room.name,
                      description: room.description,
                      visibility: room.visibility,
                      hasPassword: room.hasPassword,
                    }}
                    disabled={isClosed}
                  />
                </div>

                <div className={`rounded-xl border p-5 ${
                  isClosed
                    ? "border-border bg-surface-muted opacity-60"
                    : "border-red-500/20 bg-red-500/5"
                }`}>
                  <h4 className={`mb-2 text-sm font-semibold ${isClosed ? "text-muted-foreground" : "text-red-400"}`}>
                    {isClosed ? "Sala cerrada" : "Cerrar sala"}
                  </h4>
                  <p className="mb-4 text-xs text-muted-foreground">
                    {isClosed
                      ? "Esta sala ya está cerrada. Los miembros no pueden acceder."
                      : "Al cerrar la sala, los miembros ya no podrán acceder. Esta acción no se puede deshacer."}
                  </p>
                  <button
                    onClick={handleClose}
                    disabled={isClosed}
                    className={`cursor-pointer rounded-lg border px-4 py-2 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                      isClosed
                        ? "border-border bg-surface-muted text-muted-foreground"
                        : "border-red-500/30 bg-red-500/10 text-red-400 hover:bg-red-500/20"
                    }`}
                  >
                    {isClosed ? "Sala cerrada" : "Cerrar sala"}
                  </button>
                </div>
              </div>
            )}
          </>
        ) : null}
      </div>
    </AppLayout>
  );
}
