"use client";

import Link from "next/link";
import Image from "next/image";
import { useAuth } from "@/features/auth/context/AuthContext";
import { UserNav } from "@/shared/components/UserNav";

interface AppLayoutProps {
  children: React.ReactNode;
  hideHeader?: boolean;
}

const navLinks = [
  { href: "/", label: "Inicio" },
  { href: "/matches", label: "Partidos" },
];

const authNavLinks = [
  { href: "/predictions", label: "Mis predicciones" },
];

export function AppLayout({ children, hideHeader = false }: AppLayoutProps) {
  const { isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen bg-background">
      {!hideHeader && (
        <header className="fixed top-0 z-50 w-full border-b border-border/50 bg-background/80 backdrop-blur-lg">
          <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
            <Link href="/" className="flex items-center gap-2">
              <Image
                src="/logo_apuestec.png"
                alt="ApuesTec"
                width={0}
                height={0}
                sizes="100vw"
                className="w-42 h-auto"
                priority
              />
            </Link>
            <nav className="flex items-center gap-6">
              {navLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                >
                  {link.label}
                </Link>
              ))}
              {isAuthenticated && authNavLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                >
                  {link.label}
                </Link>
              ))}
              {isAuthenticated ? (
                <UserNav />
              ) : (
                <>
                  <Link
                    href="/login"
                    className="rounded-lg border border-primary px-4 py-2 text-sm font-medium text-primary transition-all hover:bg-primary hover:text-black"
                  >
                    Iniciar sesión
                  </Link>
                  <Link
                    href="/register"
                    className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90"
                  >
                    Registrarse
                  </Link>
                </>
              )}
            </nav>
          </div>
        </header>
      )}
      <main className={hideHeader ? "" : "pt-16"}>
        {children}
      </main>
    </div>
  );
}
