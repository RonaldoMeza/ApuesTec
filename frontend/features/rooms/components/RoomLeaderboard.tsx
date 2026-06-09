"use client";

import { useEffect, useState } from "react";
import { roomService } from "@/features/rooms/services/room.service";
import { RoomRoleBadge } from "@/features/rooms/components/RoomRoleBadge";
import type { RoomLeaderboardEntry } from "@/features/rooms/types/room.types";

interface RoomLeaderboardProps {
  roomId: string;
  roomName: string;
  currentUserId?: string;
}

export function RoomLeaderboard({ roomId, currentUserId }: RoomLeaderboardProps) {
  const [entries, setEntries] = useState<RoomLeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function fetchLeaderboard() {
      try {
        const data = await roomService.getLeaderboard(roomId);
        setEntries(data.entries);
      } catch (err: unknown) {
        const apiErr = err as { message?: string };
        setError(apiErr?.message || "Error al cargar el ranking");
      } finally {
        setLoading(false);
      }
    }
    fetchLeaderboard();
  }, [roomId]);

  if (loading) {
    return (
      <div className="rounded-xl border border-border bg-surface p-12 text-center">
        <p className="text-sm text-muted-foreground">Cargando ranking...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-12 text-center">
        <p className="text-sm text-red-400">{error}</p>
      </div>
    );
  }

  if (entries.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-surface p-12 text-center">
        <p className="text-lg text-muted-foreground">
          No hay datos de ranking disponibles.
        </p>
        <p className="mt-2 text-sm text-muted-foreground">
          Los puntos aparecerán cuando los partidos finalicen y se calculen las puntuaciones.
        </p>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-4">
        <p className="text-xs text-muted-foreground italic">
          El ranking de sala usa el puntaje global acumulado
        </p>
      </div>

      <div className="overflow-x-auto rounded-xl border border-border">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border bg-surface-muted">
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                #
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                Usuario
              </th>
              <th className="px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-muted-foreground">
                Rol
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground">
                Puntos
              </th>
              <th className="hidden px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground md:table-cell">
                Predicciones
              </th>
              <th className="hidden px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground md:table-cell">
                Exactos
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {entries.map((entry) => {
              const isCurrentUser = entry.userId === currentUserId;
              const rankColor =
                entry.rank === 1
                  ? "text-yellow-400"
                  : entry.rank === 2
                    ? "text-gray-300"
                    : entry.rank === 3
                      ? "text-amber-600"
                      : "text-muted-foreground";

              return (
                <tr
                  key={entry.userId}
                  className={`transition-colors hover:bg-surface-muted ${
                    isCurrentUser ? "bg-primary/5 ring-1 ring-primary/20" : ""
                  }`}
                >
                  <td className={`px-4 py-4 text-sm font-bold ${rankColor}`}>
                    #{entry.rank}
                  </td>
                  <td className="px-4 py-4">
                    <div className="flex items-center gap-2">
                      <div className="flex h-8 w-8 items-center justify-center rounded-full bg-gradient-to-br from-primary to-amber-500 text-xs font-bold text-black">
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
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-4 text-center">
                    <RoomRoleBadge role={entry.roomRole} />
                  </td>
                  <td className="px-4 py-4 text-right">
                    <span className="text-lg font-bold text-primary">
                      {entry.totalPoints}
                    </span>
                  </td>
                  <td className="hidden px-4 py-4 text-right text-sm text-muted-foreground md:table-cell">
                    {entry.predictionsCount}
                  </td>
                  <td className="hidden px-4 py-4 text-right text-sm text-muted-foreground md:table-cell">
                    {entry.exactScoresCount}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
