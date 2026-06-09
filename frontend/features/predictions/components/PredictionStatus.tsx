"use client";

interface PredictionStatusProps {
  canEdit: boolean;
  isLocked: boolean;
  matchStatus: string;
}

export function PredictionStatus({ canEdit, isLocked, matchStatus }: PredictionStatusProps) {
  if (matchStatus === "FINISHED" || matchStatus === "CANCELLED") {
    return (
      <span className="inline-flex items-center rounded-full border border-gray-500/30 bg-gray-500/10 px-2.5 py-0.5 text-xs font-medium text-gray-400">
        Partido {matchStatus === "FINISHED" ? "finalizado" : "cancelado"}
      </span>
    );
  }

  if (matchStatus === "LOCKED" || isLocked) {
    return (
      <span className="inline-flex items-center rounded-full border border-amber-500/30 bg-amber-500/10 px-2.5 py-0.5 text-xs font-medium text-amber-400">
        Bloqueada
      </span>
    );
  }

  if (!canEdit) {
    return (
      <span className="inline-flex items-center rounded-full border border-red-500/30 bg-red-500/10 px-2.5 py-0.5 text-xs font-medium text-red-400">
        Edición cerrada
      </span>
    );
  }

  return (
    <span className="inline-flex items-center rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2.5 py-0.5 text-xs font-medium text-emerald-400">
      Editable
    </span>
  );
}
