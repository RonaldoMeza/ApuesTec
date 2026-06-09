"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { AppLayout } from "@/shared/components/AppLayout";
import { matchService } from "@/features/matches/services/match.service";
import { leaderboardService } from "@/features/leaderboard/services/leaderboard.service";
import type { MatchResponse } from "@/features/matches/types/match.types";
import type { LeaderboardEntry } from "@/features/leaderboard/types/leaderboard.types";

function HeroSection() {
  return (
    <section className="relative flex min-h-screen items-center justify-center overflow-hidden">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-primary/20 via-transparent to-transparent" />
      <div className="relative z-10 mx-auto max-w-4xl px-4 text-center">
        <h1 className="bg-gradient-to-r from-primary via-amber-400 to-primary bg-clip-text text-5xl font-bold leading-tight text-transparent md:text-7xl">
          ApuesTec
        </h1>
        <p className="mt-6 text-lg text-muted-foreground md:text-xl">
          La plataforma educativa de predicciones deportivas del Mundial.
          Sin dinero real, solo diversión y competencia.
        </p>
        <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <Link
            href="/register"
            className="rounded-xl bg-gradient-to-r from-primary to-amber-500 px-8 py-3 text-base font-semibold text-primary-foreground shadow-lg shadow-primary/25 transition-all hover:scale-105 hover:shadow-xl hover:shadow-primary/30"
          >
            Participar ahora
          </Link>
          <Link
            href="/login"
            className="rounded-xl border border-border px-8 py-3 text-base font-medium text-foreground transition-all hover:bg-surface-hover"
          >
            Iniciar sesión
          </Link>
        </div>
      </div>
    </section>
  );
}

function HowItWorks() {
  const steps = [
    { icon: "🎯", title: "Predice marcadores", desc: "Anticipa los resultados de cada partido del Mundial." },
    { icon: "⭐", title: "Gana puntos", desc: "Acumula puntos por cada acierto y sube en la clasificación." },
    { icon: "🏆", title: "Compite en salas", desc: "Crea o únete a salas privadas con tus amigos." },
    { icon: "📈", title: "Sube en el ranking", desc: "Compara tu puntuación con otros jugadores." },
  ];

  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <h2 className="text-center text-3xl font-bold text-foreground">Cómo funciona</h2>
        <p className="mt-2 text-center text-muted-foreground">Cuatro pasos para empezar a disfrutar</p>
        <div className="mt-12 grid gap-6 md:grid-cols-4">
          {steps.map((step, i) => (
            <div key={i} className="group rounded-xl border border-border bg-surface p-6 transition-all hover:border-primary/50 hover:shadow-lg hover:shadow-primary/5">
              <div className="mb-4 text-3xl">{step.icon}</div>
              <h3 className="mb-2 font-semibold text-foreground">{step.title}</h3>
              <p className="text-sm text-muted-foreground">{step.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function ScoringRules() {
  const rules = [
    { label: "Marcador exacto", points: "+5", desc: "Aciertas el resultado exacto del partido" },
    { label: "Ganador / Empate", points: "+3", desc: "Aciertas quién gana o si hay empate" },
    { label: "Diferencia de goles", points: "+2", desc: "Aciertas la diferencia de goles" },
    { label: "Predicción anticipada", points: "+1", desc: "Predices antes del plazo límite" },
    { label: "Racha de aciertos", points: "+2", desc: "Cada 3 aciertos consecutivos" },
  ];

  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <h2 className="text-center text-3xl font-bold text-foreground">Sistema de puntuación</h2>
        <p className="mt-2 text-center text-muted-foreground">Así se calculan los puntos</p>
        <div className="mt-10 space-y-3">
          {rules.map((rule, i) => (
            <div key={i} className="flex items-center justify-between rounded-xl border border-border bg-surface px-6 py-4 transition-all hover:border-primary/30">
              <div>
                <span className="font-medium text-foreground">{rule.label}</span>
                <p className="text-sm text-muted-foreground">{rule.desc}</p>
              </div>
              <span className="text-2xl font-bold text-primary">{rule.points}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function PreviewSection() {
  const [matches, setMatches] = useState<MatchResponse[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    matchService
      .listUpcoming()
      .then(setMatches)
      .catch(() => setMatches([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold text-foreground">Próximos partidos</h2>
            <p className="mt-2 text-muted-foreground">Prepárate para hacer tus predicciones</p>
          </div>
          <Link
            href="/matches"
            className="rounded-lg border border-primary px-4 py-2 text-sm font-medium text-primary transition-all hover:bg-primary hover:text-black"
          >
            Ver partidos
          </Link>
        </div>
        <div className="mt-10 grid gap-4 md:grid-cols-3">
          {loading ? (
            <div className="col-span-full flex items-center justify-center py-10">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
            </div>
          ) : matches.length === 0 ? (
            <div className="col-span-full rounded-xl border border-border bg-surface p-8 text-center">
              <p className="text-lg text-muted-foreground">
                No hay partidos próximos. Vuelve más tarde.
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                Los administradores estarán cargando los partidos del Mundial pronto.
              </p>
            </div>
          ) : (
            matches.slice(0, 3).map((match) => (
              <Link key={match.id} href={`/matches/${match.id}`}>
                <div className="rounded-xl border border-border bg-surface p-6 text-center transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5">
                  <p className="text-xs text-muted-foreground">
                    {new Date(match.startsAt).toLocaleDateString("es-PE", {
                      day: "numeric",
                      month: "long",
                    })}
                  </p>
                  <p className="mt-2 text-lg font-semibold text-foreground">
                    {match.homeTeam?.name || "?"}
                  </p>
                  <p className="text-sm text-muted-foreground">vs</p>
                  <p className="text-lg font-semibold text-foreground">
                    {match.awayTeam?.name || "?"}
                  </p>
                </div>
              </Link>
            ))
          )}
        </div>
        {!loading && matches.length > 0 && (
          <div className="mt-6 text-center">
            <Link
              href="/matches"
              className="text-sm text-primary hover:underline"
            >
              Ver todos los partidos →
            </Link>
          </div>
        )}
      </div>
    </section>
  );
}

function RankingPreview() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    leaderboardService
      .getGlobalLeaderboard(5)
      .then((data) => setEntries(data.entries))
      .catch(() => setEntries([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold text-foreground">Ranking general</h2>
            <p className="mt-2 text-muted-foreground">Los mejores predictors de ApuesTec</p>
          </div>
          <Link
            href="/leaderboard"
            className="rounded-lg border border-primary px-4 py-2 text-sm font-medium text-primary transition-all hover:bg-primary hover:text-black"
          >
            Ver ranking
          </Link>
        </div>
        <div className="mt-10">
          {loading ? (
            <div className="flex items-center justify-center py-10">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
            </div>
          ) : entries.length === 0 ? (
            <div className="rounded-xl border border-border bg-surface p-8 text-center">
              <p className="text-lg text-muted-foreground">
                El ranking aparecerá cuando los partidos finalicen y se calculen las puntuaciones.
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                Regístrate y prepárate para competir.
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {entries.map((entry) => (
                <div
                  key={entry.userId}
                  className="flex items-center justify-between rounded-xl border border-border bg-surface px-6 py-4 transition-all hover:border-primary/30"
                >
                  <div className="flex items-center gap-4">
                    <span className="text-lg font-bold text-muted-foreground">
                      {entry.rank === 1 ? "🥇" : entry.rank === 2 ? "🥈" : entry.rank === 3 ? "🥉" : `#${entry.rank}`}
                    </span>
                    <div>
                      <p className="font-medium text-foreground">{entry.fullName}</p>
                      <p className="text-xs text-muted-foreground">{entry.predictionsCount} predicciones</p>
                    </div>
                  </div>
                  <span className="text-2xl font-bold text-primary">{entry.totalPoints}</span>
                </div>
              ))}
            </div>
          )}
        </div>
        {!loading && entries.length > 0 && (
          <div className="mt-6 text-center">
            <Link
              href="/leaderboard"
              className="text-sm text-primary hover:underline"
            >
              Ver ranking completo →
            </Link>
          </div>
        )}
      </div>
    </section>
  );
}

function PrivateRoomsSection() {
  return (
    <section className="relative py-24">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-primary/10 via-transparent to-transparent" />
      <div className="relative z-10 mx-auto max-w-6xl px-4 text-center">
        <h2 className="text-3xl font-bold text-foreground">Salas</h2>
        <p className="mx-auto mt-4 max-w-2xl text-muted-foreground">
          Crea salas privadas o públicas, invita con códigos temporales, busca salas en tu red o
          únete con contraseña.
        </p>
        <div className="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <Link
            href="/rooms"
            className="rounded-xl bg-gradient-to-r from-primary to-amber-500 px-8 py-3 text-base font-semibold text-primary-foreground shadow-lg shadow-primary/25 transition-all hover:scale-105 hover:shadow-xl hover:shadow-primary/30"
          >
            Explorar salas
          </Link>
          <Link
            href="/register"
            className="rounded-xl border border-border px-8 py-3 text-base font-medium text-foreground transition-all hover:bg-surface-hover"
          >
            Registrarse
          </Link>
        </div>
      </div>
    </section>
  );
}

function Footer() {
  return (
    <footer className="border-t border-border py-8">
      <div className="mx-auto max-w-6xl px-4 text-center">
        <p className="text-sm text-muted-foreground">
          ApuesTec — Plataforma educativa. No utiliza dinero real, no realiza apuestas monetarias,
          no usa cuotas reales ni se integra con casas de apuestas.
        </p>
      </div>
    </footer>
  );
}

export default function HomePage() {
  return (
    <AppLayout>
      <HeroSection />
      <HowItWorks />
      <ScoringRules />
      <PreviewSection />
      <RankingPreview />
      <PrivateRoomsSection />
      <Footer />
    </AppLayout>
  );
}
