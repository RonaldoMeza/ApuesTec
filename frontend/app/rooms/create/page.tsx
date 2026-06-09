"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { AppLayout } from "@/shared/components/AppLayout";
import { RoomForm } from "@/features/rooms/components/RoomForm";
import { roomService } from "@/features/rooms/services/room.service";
import type { CreateRoomRequest } from "@/features/rooms/types/room.types";

export default function CreateRoomPage() {
  return (
    <ProtectedRoute>
      <CreateRoomContent />
    </ProtectedRoute>
  );
}

function CreateRoomContent() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  async function handleSubmit(data: CreateRoomRequest) {
    setLoading(true);
    try {
      const room = await roomService.create(data);
      toast.success(`Sala "${room.name}" creada correctamente`);
      router.push(`/rooms/${room.id}`);
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      toast.error(apiErr?.message || "Error al crear la sala");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AppLayout>
      <div className="mx-auto max-w-2xl px-4 pt-8">
        <div className="mb-6">
          <Link
            href="/rooms"
            className="text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            &larr; Volver a mis salas
          </Link>
        </div>

        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Crear Sala</h1>
          <p className="mt-1 text-muted-foreground">
            Crea una sala privada o pública para competir con tus amigos
          </p>
        </div>

        <div className="rounded-xl border border-border bg-surface p-6">
          <RoomForm onSubmit={handleSubmit} isLoading={loading} />
        </div>
      </div>
    </AppLayout>
  );
}
