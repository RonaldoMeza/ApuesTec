"use client";

const statusConfig: Record<string, { label: string; className: string }> = {
  SCHEDULED: {
    label: "Programado",
    className: "bg-blue-500/10 text-blue-400 border-blue-500/30",
  },
  LOCKED: {
    label: "Bloqueado",
    className: "bg-amber-500/10 text-amber-400 border-amber-500/30",
  },
  FINISHED: {
    label: "Finalizado",
    className: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
  },
  CANCELLED: {
    label: "Cancelado",
    className: "bg-red-500/10 text-red-400 border-red-500/30",
  },
};

export function MatchStatusBadge({ status }: { status: string }) {
  const config = statusConfig[status] || {
    label: status,
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
