"use client";

import Link from "next/link";
import { AppLayout } from "@/shared/components/AppLayout";

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
  const previewMatches = [
    { home: "Argentina", away: "Brasil", date: "Próximamente" },
    { home: "España", away: "Francia", date: "Próximamente" },
    { home: "Alemania", away: "Inglaterra", date: "Próximamente" },
  ];

  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <h2 className="text-center text-3xl font-bold text-foreground">Próximos partidos</h2>
        <p className="mt-2 text-center text-muted-foreground">Prepárate para hacer tus predicciones</p>
        <div className="mt-10 grid gap-4 md:grid-cols-3">
          {previewMatches.map((match, i) => (
            <div key={i} className="rounded-xl border border-border bg-surface p-6 text-center transition-all hover:border-primary/30">
              <p className="text-xs text-muted-foreground">{match.date}</p>
              <p className="mt-2 text-lg font-semibold text-foreground">{match.home}</p>
              <p className="text-sm text-muted-foreground">vs</p>
              <p className="text-lg font-semibold text-foreground">{match.away}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function RankingPreview() {
  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <h2 className="text-center text-3xl font-bold text-foreground">Ranking general</h2>
        <p className="mt-2 text-center text-muted-foreground">Los mejores predictors de ApuesTec</p>
        <div className="mt-10 rounded-xl border border-border bg-surface p-8 text-center">
          <p className="text-lg text-muted-foreground">
            El ranking estará disponible cuando comience la temporada de predicciones.
          </p>
          <p className="mt-2 text-sm text-muted-foreground">
            Regístrate y prepárate para competir.
          </p>
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
      <Footer />
    </AppLayout>
  );
}
