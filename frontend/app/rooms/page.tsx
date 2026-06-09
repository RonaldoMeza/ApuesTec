"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { AppLayout } from "@/shared/components/AppLayout";
import { roomService } from "@/features/rooms/services/room.service";
import { inviteService } from "@/features/invites/services/invite.service";
import { RoomCard } from "@/features/rooms/components/RoomCard";
import type { Room, PublicRoom } from "@/features/rooms/types/room.types";

type Tab = "my-rooms" | "discover";

export default function RoomsPage() {
  return (
    <ProtectedRoute>
      <RoomsContent />
    </ProtectedRoute>
  );
}

function RoomsContent() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState<Tab>("my-rooms");
  const [rooms, setRooms] = useState<Room[]>([]);
  const [publicRooms, setPublicRooms] = useState<PublicRoom[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [searchingPublic, setSearchingPublic] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(null);

  const [inviteCode, setInviteCode] = useState("");
  const [joiningCode, setJoiningCode] = useState(false);
  const [inviteError, setInviteError] = useState("");

  const [passwordModal, setPasswordModal] = useState<{
    room: PublicRoom;
    password: string;
  } | null>(null);
  const [joiningPublic, setJoiningPublic] = useState(false);

  useEffect(() => {
    loadMyRooms();
  }, []);

  async function loadMyRooms() {
    setLoading(true);
    setError("");
    try {
      const data = await roomService.listMyRooms();
      setRooms(data.rooms);
    } catch {
      setError("Error al cargar tus salas");
    } finally {
      setLoading(false);
    }
  }

  const searchPublic = useCallback(async (query: string) => {
    setSearchingPublic(true);
    setError("");
    try {
      const data = await roomService.searchPublic(query);
      setPublicRooms(data.rooms);
    } catch {
      setError("Error al buscar salas");
    } finally {
      setSearchingPublic(false);
    }
  }, []);

  function handleSearchChange(value: string) {
    setSearchQuery(value);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => {
      searchPublic(value);
    }, 400);
  }

  function handleTabChange(tab: Tab) {
    setActiveTab(tab);
    setError("");
    if (tab === "discover") {
      searchPublic(searchQuery);
    } else {
      loadMyRooms();
    }
  }

  async function handleJoinByCode() {
    if (!inviteCode.trim()) return;
    setJoiningCode(true);
    setInviteError("");
    try {
      const result = await inviteService.join(inviteCode.trim());
      toast.success(`Te uniste a "${result.roomName}"`);
      router.push(`/rooms/${result.roomId}`);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      const msg = apiErr?.message || "Código inválido o expirado";
      setInviteError(msg);
      toast.error(msg);
    } finally {
      setJoiningCode(false);
    }
  }

  async function handleJoinPublic(room: PublicRoom) {
    if (room.hasPassword) {
      setPasswordModal({ room, password: "" });
      return;
    }
    setJoiningPublic(true);
    try {
      const result = await roomService.joinPublic(room.id);
      toast.success(`Te uniste a "${result.name}"`);
      router.push(`/rooms/${result.id}`);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al unirte a la sala");
    } finally {
      setJoiningPublic(false);
    }
  }

  async function handleJoinWithPassword() {
    if (!passwordModal) return;
    setJoiningPublic(true);
    try {
      const result = await roomService.joinPublic(passwordModal.room.id, {
        password: passwordModal.password,
      });
      toast.success(`Te uniste a "${result.name}"`);
      setPasswordModal(null);
      router.push(`/rooms/${result.id}`);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al unirte a la sala");
    } finally {
      setJoiningPublic(false);
    }
  }

  const tabs: { key: Tab; label: string }[] = [
    { key: "my-rooms", label: "Mis Salas" },
    { key: "discover", label: "Salas públicas" },
  ];

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8">
        <div className="relative mb-8 overflow-hidden rounded-2xl border border-border bg-gradient-to-br from-primary/10 via-background to-amber-500/5 p-6">
          <div className="relative z-10 flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Salas</h1>
              <p className="mt-1 text-muted-foreground">
                Administra tus salas o descubre salas públicas en tu red
              </p>
            </div>
            <Link
              href="/rooms/create"
              className="cursor-pointer rounded-xl bg-gradient-to-r from-primary to-amber-500 px-5 py-2.5 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] hover:shadow-xl hover:shadow-primary/30"
            >
              + Crear sala
            </Link>
          </div>
        </div>

        <div className="mb-6">
          <div className="rounded-xl border border-border bg-surface-muted p-4">
            <h3 className="mb-3 text-sm font-semibold text-foreground flex items-center gap-2">
              <svg className="h-4 w-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
              Unirse por código de invitación
            </h3>
            <form
              onSubmit={(e) => { e.preventDefault(); handleJoinByCode(); }}
              className="flex gap-3"
            >
              <input
                type="text"
                value={inviteCode}
                onChange={(e) => { setInviteCode(e.target.value); setInviteError(""); }}
                placeholder="Ej: a1b2c3d4e5f6..."
                className="flex-1 rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                disabled={joiningCode}
              />
              <button
                type="submit"
                disabled={joiningCode || !inviteCode.trim()}
                className="cursor-pointer rounded-lg bg-primary px-6 py-2 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {joiningCode ? "Uniéndose..." : "Unirse"}
              </button>
            </form>
            {inviteError && (
              <p className="mt-2 text-xs text-red-400">{inviteError}</p>
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
                onClick={() => handleTabChange(tab.key)}
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

        {activeTab === "my-rooms" && (
          <>
            {loading ? (
              <div className="flex items-center justify-center py-20">
                <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : error ? (
              <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
                <p className="text-red-400">{error}</p>
              </div>
            ) : rooms.length === 0 ? (
              <div className="rounded-xl border border-border bg-surface p-12 text-center">
                <svg className="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
                <p className="text-lg text-muted-foreground mb-4">
                  No perteneces a ninguna sala
                </p>
                <div className="flex flex-col items-center gap-3">
                  <Link
                    href="/rooms/create"
                    className="cursor-pointer inline-block rounded-xl bg-gradient-to-r from-primary to-amber-500 px-6 py-3 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] hover:shadow-xl hover:shadow-primary/30"
                  >
                    Crear sala
                  </Link>
                  <button
                    onClick={() => handleTabChange("discover")}
                    className="cursor-pointer text-sm text-primary hover:underline"
                  >
                    Buscar salas públicas
                  </button>
                </div>
              </div>
            ) : (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {rooms.map((room) => (
                  <RoomCard key={room.id} room={room} />
                ))}
              </div>
            )}
          </>
        )}

        {activeTab === "discover" && (
          <>
            <div className="mb-6">
              <div className="relative">
                <svg className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => handleSearchChange(e.target.value)}
                  placeholder="Buscar salas públicas en tu red..."
                  className="w-full rounded-xl border border-border bg-surface py-3 pl-10 pr-4 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-all"
                  autoFocus
                />
              </div>
              {searchingPublic && (
                <div className="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
                  <div className="h-3 w-3 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                  Buscando...
                </div>
              )}
            </div>

            {loading ? (
              <div className="flex items-center justify-center py-20">
                <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : error ? (
              <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-8 text-center">
                <p className="text-red-400">{error}</p>
              </div>
            ) : publicRooms.length === 0 ? (
              <div className="rounded-xl border border-border bg-surface p-12 text-center">
                <svg className="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <p className="text-lg text-muted-foreground">
                  {searchQuery
                    ? "No se encontraron salas con ese nombre en tu red"
                    : "No hay salas públicas en tu red actualmente"}
                </p>
                <p className="mt-2 text-sm text-muted-foreground">
                  {searchQuery
                    ? "Prueba con otro término de búsqueda"
                    : "Pide a alguien que cree una sala pública"}
                </p>
              </div>
            ) : (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {publicRooms.map((room) => (
                  <div
                    key={room.id}
                    className="rounded-xl border border-border bg-surface p-5 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5"
                  >
                    <div className="mb-2 flex items-center gap-2">
                      <span className="inline-flex items-center rounded-full bg-emerald-500/10 px-2.5 py-0.5 text-xs font-medium text-emerald-400 border border-emerald-500/30">
                        Pública
                      </span>
                      {room.hasPassword && (
                        <span className="inline-flex items-center rounded-full bg-amber-500/10 px-2.5 py-0.5 text-xs font-medium text-amber-400 border border-amber-500/30">
                          Con contraseña
                        </span>
                      )}
                    </div>

                    <h3 className="text-lg font-semibold text-foreground">{room.name}</h3>
                    {room.description && (
                      <p className="mt-1 text-sm text-muted-foreground line-clamp-2">
                        {room.description}
                      </p>
                    )}
                    <p className="mt-2 text-xs text-muted-foreground">
                      Creada por {room.ownerName} &middot; {room.memberCount} miembro
                      {room.memberCount !== 1 ? "s" : ""}
                    </p>

                    <button
                      onClick={() => handleJoinPublic(room)}
                      disabled={room.isMember || joiningPublic}
                      className={`cursor-pointer mt-4 w-full rounded-lg py-2 text-sm font-medium transition-all disabled:cursor-not-allowed ${
                        room.isMember
                          ? "border border-border bg-surface-muted text-muted-foreground"
                          : "bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                      }`}
                    >
                      {room.isMember ? "Ya eres miembro" : "Unirse"}
                    </button>
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>

      {passwordModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-2xl border border-border bg-surface p-6 shadow-2xl animate-in fade-in zoom-in-95">
            <h3 className="text-lg font-semibold text-foreground mb-1">
              Unirse a &ldquo;{passwordModal.room.name}&rdquo;
            </h3>
            <p className="text-sm text-muted-foreground mb-4">
              Esta sala requiere contraseña
            </p>

            <input
              type="text"
              value={passwordModal.password}
              onChange={(e) =>
                setPasswordModal({ ...passwordModal, password: e.target.value })
              }
              placeholder="Contraseña de la sala"
              className="w-full rounded-lg border border-border bg-surface-muted px-3 py-2.5 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary mb-4"
              autoFocus
            />

            <div className="flex gap-3">
              <button
                onClick={() => setPasswordModal(null)}
                className="cursor-pointer flex-1 rounded-lg border border-border px-4 py-2.5 text-sm font-medium text-foreground transition-all hover:bg-surface-muted"
              >
                Cancelar
              </button>
              <button
                onClick={handleJoinWithPassword}
                disabled={joiningPublic || !passwordModal.password}
                className="cursor-pointer flex-1 rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {joiningPublic ? "Uniéndose..." : "Unirse"}
              </button>
            </div>
          </div>
        </div>
      )}
    </AppLayout>
  );
}
