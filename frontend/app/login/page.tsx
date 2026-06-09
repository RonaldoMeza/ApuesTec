"use client";

import Image from "next/image";
import Link from "next/link";
import { LoginForm } from "@/features/auth/components/LoginForm";

export default function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <Link href="/" className="inline-flex items-center gap-2">
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
          <h1 className="mt-6 text-xl font-semibold text-foreground">Iniciar sesión</h1>

        </div>
        <div className="rounded-xl border border-border bg-surface p-6 shadow-lg shadow-black/20">
          <LoginForm />
        </div>
      </div>
    </div>
  );
}
