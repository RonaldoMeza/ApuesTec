"use client";

import { useState, useEffect } from "react";
import { AppLayout } from "@/shared/components/AppLayout";
import { LeaderboardTable } from "@/features/leaderboard/components/LeaderboardTable";
import { leaderboardService } from "@/features/leaderboard/services/leaderboard.service";
import { useAuth } from "@/features/auth/context/AuthContext";
import type { LeaderboardEntry } from "@/features/leaderboard/types/leaderboard.types";

export default function LeaderboardPage() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [limit, setLimit] = useState(50);
  const { user } = useAuth();

  useEffect(() => {
    let cancelled = false;

    leaderboardService
      .getGlobalLeaderboard(limit)
      .then((data) => {
        if (!cancelled) {
          setEntries(data.entries);
          setTotal(data.total);
          setLoading(false);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setError("Error al cargar el ranking");
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [limit]);

  return (
    <AppLayout>
      <div className="mx-auto max-w-5xl px-4 pt-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Ranking Global</h1>
          <p className="mt-1 text-muted-foreground">
            Los mejores predictors de ApuesTec
          </p>
          {total > 0 && (
            <p className="mt-1 text-sm text-muted-foreground">
              {total} participantes
            </p>
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
        ) : (
          <>
            <LeaderboardTable entries={entries} currentUserId={user?.id} />

            {total > limit && (
              <div className="mt-6 text-center">
                <button
                  onClick={() => setLimit(limit + 50)}
                  className="rounded-lg border border-primary px-6 py-2 text-sm font-medium text-primary transition-all hover:bg-primary hover:text-black"
                >
                  Cargar más ({total - limit} restantes)
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </AppLayout>
  );
}
