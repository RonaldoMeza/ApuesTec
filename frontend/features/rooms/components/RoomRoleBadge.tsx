"use client";

const roleConfig: Record<string, { label: string; className: string }> = {
  OWNER: {
    label: "Dueño",
    className: "bg-amber-500/10 text-amber-400 border-amber-500/30",
  },
  MODERATOR: {
    label: "Moderador",
    className: "bg-yellow-500/10 text-yellow-400 border-yellow-500/30",
  },
  MEMBER: {
    label: "Miembro",
    className: "bg-gray-500/10 text-gray-400 border-gray-500/30",
  },
};

export function RoomRoleBadge({ role }: { role: string }) {
  const config = roleConfig[role] || {
    label: role,
    className: "bg-gray-500/10 text-gray-400 border-gray-500/30",
  };

  return (
    <span
      className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${config.className}`}
    >
      {config.label}
    </span>
  );
}
