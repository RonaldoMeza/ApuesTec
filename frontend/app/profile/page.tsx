"use client";

import { useState, type FormEvent } from "react";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { useAuth } from "@/features/auth/context/AuthContext";
import { LoadingScreen } from "@/shared/components/LoadingScreen";
import { AppLayout } from "@/shared/components/AppLayout";
import { Button } from "@/components/ui/button";

export default function ProfilePage() {
  return (
    <ProtectedRoute>
      <ProfileContent />
    </ProtectedRoute>
  );
}

function ProfileContent() {
  const { user, isLoading, changePassword } = useAuth();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleChangePassword(e: FormEvent) {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (!currentPassword || !newPassword) {
      setError("Todos los campos son obligatorios");
      return;
    }

    if (newPassword.length < 8) {
      setError("La nueva contraseña debe tener al menos 8 caracteres");
      return;
    }

    setSubmitting(true);
    try {
      await changePassword({ currentPassword, newPassword });
      setSuccess("Contraseña cambiada exitosamente.");
      setCurrentPassword("");
      setNewPassword("");
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      setError(apiErr?.message || "Error al cambiar la contraseña");
    } finally {
      setSubmitting(false);
    }
  }

  if (isLoading || !user) {
    return <LoadingScreen message="Cargando perfil..." />;
  }

  return (
    <AppLayout>
      <div className="mx-auto max-w-6xl px-4 pt-8 pb-12">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-foreground">Mi Perfil</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Información personal y configuración de la cuenta
          </p>
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <div className="rounded-xl border border-border bg-surface p-6 shadow-lg shadow-black/20">
            <h2 className="mb-1 text-lg font-semibold text-foreground">Información personal</h2>
            <p className="mb-4 text-sm text-muted-foreground">Datos asociados a tu cuenta</p>
            <dl className="space-y-3 text-sm">
              <div className="flex justify-between rounded-lg bg-surface-muted px-4 py-3">
                <dt className="text-muted-foreground">Nombre</dt>
                <dd className="font-medium text-foreground">{user.fullName}</dd>
              </div>
              <div className="flex justify-between rounded-lg bg-surface-muted px-4 py-3">
                <dt className="text-muted-foreground">Email</dt>
                <dd className="font-medium text-foreground">{user.email}</dd>
              </div>
              <div className="flex justify-between rounded-lg bg-surface-muted px-4 py-3">
                <dt className="text-muted-foreground">Roles</dt>
                <dd className="font-medium text-foreground">{user.roles.join(", ")}</dd>
              </div>
            </dl>
          </div>

          <div className="rounded-xl border border-border bg-surface p-6 shadow-lg shadow-black/20">
            <h2 className="mb-1 text-lg font-semibold text-foreground">Cambiar contraseña</h2>
            <p className="mb-4 text-sm text-muted-foreground">Actualiza tu contraseña de acceso</p>

            {error && (
              <div className="mb-4 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400" role="alert">
                {error}
              </div>
            )}
            {success && (
              <div className="mb-4 rounded-lg border border-green-500/30 bg-green-500/10 px-4 py-3 text-sm text-green-400" role="alert">
                {success}
              </div>
            )}

            <form onSubmit={handleChangePassword} className="space-y-4">
              <div>
                <label htmlFor="currentPassword" className="block text-sm font-medium text-foreground">
                  Contraseña actual
                </label>
                <input
                  id="currentPassword"
                  type="password"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  className="mt-1 block w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-sm text-foreground shadow-sm placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  disabled={submitting}
                  required
                />
              </div>

              <div>
                <label htmlFor="newPassword" className="block text-sm font-medium text-foreground">
                  Nueva contraseña
                </label>
                <input
                  id="newPassword"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="mt-1 block w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-sm text-foreground shadow-sm placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Mínimo 8 caracteres"
                  disabled={submitting}
                  required
                  minLength={8}
                />
              </div>

              <Button type="submit" disabled={submitting} variant="default">
                {submitting ? "Cambiando contraseña..." : "Cambiar contraseña"}
              </Button>
            </form>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
