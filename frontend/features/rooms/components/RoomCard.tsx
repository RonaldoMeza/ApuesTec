"use client";

import Link from "next/link";
import { RoomRoleBadge } from "@/features/rooms/components/RoomRoleBadge";
import type { Room } from "@/features/rooms/types/room.types";

interface RoomCardProps {
  room: Room;
}

const statusConfig: Record<string, { label: string; className: string }> = {
  ACTIVE: {
    label: "Activa",
    className: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
  },
  CLOSED: {
    label: "Cerrada",
    className: "bg-red-500/10 text-red-400 border-red-500/30",
  },
};

export function RoomCard({ room }: RoomCardProps) {
  const status = statusConfig[room.status] || {
    label: room.status,
    className: "bg-gray-500/10 text-gray-400 border-gray-500/30",
  };

  return (
    <Link href={`/rooms/${room.id}`}>
      <div className="group rounded-xl border border-border bg-surface p-5 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5">
        <div className="mb-3 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span
              className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${status.className}`}
            >
              {status.label}
            </span>
            {room.visibility === "PUBLIC" && (
              <span className="inline-flex items-center rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2.5 py-0.5 text-xs font-medium text-emerald-400">
                Pública
              </span>
            )}
          </div>
          <RoomRoleBadge role={room.myRole} />
        </div>

        <h3 className="mb-1 text-lg font-semibold text-foreground group-hover:text-primary transition-colors">
          {room.name}
        </h3>

        <p className="mb-4 text-sm text-muted-foreground line-clamp-2">
          {room.description || "Sin descripción"}
        </p>

        <div className="flex items-center justify-between border-t border-border/50 pt-3 text-xs text-muted-foreground">
          <span>{room.memberCount} miembro{room.memberCount !== 1 ? "s" : ""}</span>
          <span className="text-primary">
            Ver sala &rarr;
          </span>
        </div>
      </div>
    </Link>
  );
}
