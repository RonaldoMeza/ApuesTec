"use client";

import { useRouter } from "next/navigation";
import { useAuth } from "@/features/auth/context/AuthContext";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LayoutDashboard, User, LogOut, Shield, Trophy, FileText, BarChart3 } from "lucide-react";

export function UserNav() {
  const { user, logout } = useAuth();
  const router = useRouter();

  if (!user) return null;

  const initials = user.fullName
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button className="flex items-center gap-2 rounded-full outline-none transition-opacity hover:opacity-80">
          <Avatar size="sm">
            <AvatarFallback>{initials}</AvatarFallback>
          </Avatar>
          <span className="hidden text-sm text-foreground sm:inline">
            {user.fullName}
          </span>
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel>
          <div className="flex flex-col">
            <span className="text-sm font-medium text-foreground">{user.fullName}</span>
            <span className="text-xs text-muted-foreground">{user.email}</span>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => router.push("/dashboard")}>
          <LayoutDashboard className="mr-2 h-4 w-4" />
          Dashboard
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => router.push("/matches")}>
          <Trophy className="mr-2 h-4 w-4" />
          Partidos
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => router.push("/predictions")}>
          <FileText className="mr-2 h-4 w-4" />
          Mis predicciones
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => router.push("/stats")}>
          <BarChart3 className="mr-2 h-4 w-4" />
          Mis estadísticas
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => router.push("/profile")}>
          <User className="mr-2 h-4 w-4" />
          Perfil
        </DropdownMenuItem>
        {user.roles.some((r) => r === "ADMIN" || r === "SUPER_ADMIN") && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => router.push("/admin")}>
              <Shield className="mr-2 h-4 w-4" />
              Admin
            </DropdownMenuItem>
          </>
        )}
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={logout} variant="destructive">
          <LogOut className="mr-2 h-4 w-4" />
          Cerrar sesión
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
