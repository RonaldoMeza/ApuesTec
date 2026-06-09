"use client";

import { RoomRoleBadge } from "@/features/rooms/components/RoomRoleBadge";
import type { RoomMember } from "@/features/rooms/types/room.types";

interface RoomMemberListProps {
  members: RoomMember[];
  currentUserId: string;
  myRole: string;
  onRoleChange?: (userId: string, role: string) => Promise<void>;
  onRemoveMember?: (userId: string) => Promise<void>;
}

export function RoomMemberList({
  members,
  currentUserId,
  myRole,
  onRoleChange,
  onRemoveMember,
}: RoomMemberListProps) {
  if (members.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-surface p-12 text-center">
        <p className="text-sm text-muted-foreground">No hay miembros en esta sala.</p>
      </div>
    );
  }

  const canManage = myRole === "OWNER" || myRole === "MODERATOR";

  return (
    <div className="space-y-2">
      {members.map((member) => {
        const isCurrentUser = member.userId === currentUserId;
        const isOwner = member.role === "OWNER";

        const showRoleChange = myRole === "OWNER" && !isOwner && onRoleChange;
        const showRemove = !isOwner && onRemoveMember && (
          myRole === "OWNER" || (myRole === "MODERATOR" && member.role === "MEMBER")
        );
        const targetRole = member.role === "MEMBER" ? "MODERATOR" : "MEMBER";

        return (
          <div
            key={member.id}
            className={`flex items-center justify-between rounded-xl border border-border bg-surface p-4 transition-colors ${
              isCurrentUser ? "ring-1 ring-primary/20 bg-primary/5" : ""
            }`}
          >
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-primary to-amber-500 text-sm font-bold text-black">
                {member.fullName
                  .split(" ")
                  .map((n) => n[0])
                  .join("")
                  .toUpperCase()
                  .slice(0, 2)}
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">
                  {member.fullName}
                  {isCurrentUser && (
                    <span className="ml-2 text-xs text-primary">(tú)</span>
                  )}
                </p>
                <p className="text-xs text-muted-foreground">{member.email}</p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <RoomRoleBadge role={member.role} />

              {canManage && !isCurrentUser && (
                <div className="flex gap-1">
                    {showRoleChange && (
                      <button
                        type="button"
                        onClick={() => onRoleChange(member.userId, targetRole)}
                        className="cursor-pointer rounded-md border border-border bg-surface-muted px-2 py-1 text-xs text-muted-foreground transition-colors hover:border-primary/30 hover:text-primary"
                        title={`Cambiar a ${targetRole === "MODERATOR" ? "Moderador" : "Miembro"}`}
                      >
                        {targetRole === "MODERATOR" ? "Ascender" : "Descender"}
                      </button>
                    )}
                    {showRemove && (
                      <button
                        type="button"
                        onClick={() => onRemoveMember(member.userId)}
                        className="cursor-pointer rounded-md border border-red-500/30 bg-red-500/10 px-2 py-1 text-xs text-red-400 transition-colors hover:bg-red-500/20"
                        title="Eliminar miembro"
                      >
                        Eliminar
                      </button>
                    )}
                </div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
