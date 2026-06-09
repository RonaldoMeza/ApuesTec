"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useRef,
  startTransition,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { authService } from "@/features/auth/services/auth.service";
import {
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
} from "@/features/auth/utils/token-storage";
import type {
  UserInfo,
  LoginRequest,
  RegisterRequest,
  ChangePasswordRequest,
} from "@/features/auth/types/auth.types";

interface AuthContextType {
  user: UserInfo | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  changePassword: (data: ChangePasswordRequest) => Promise<void>;
  getMe: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();
  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;

    const token = getAccessToken();
    if (!token) {
      startTransition(() => setIsLoading(false));
      return;
    }

    let cancelled = false;
    authService.me()
      .then((userInfo) => {
        startTransition(() => {
          if (!cancelled) setUser(userInfo);
        });
      })
      .catch(() => {
        if (!cancelled) clearTokens();
      })
      .finally(() => {
        startTransition(() => {
          if (!cancelled) setIsLoading(false);
        });
      });

    return () => { cancelled = true; };
  }, [router]);

  const login = useCallback(async (data: LoginRequest) => {
    const response = await authService.login(data);
    setTokens(response.accessToken, response.refreshToken);
    setUser(response.user);
    toast.success("Inicio de sesión exitoso", {
      description: `Bienvenido, ${response.user.fullName}`,
    });
    router.push("/dashboard");
  }, [router]);

  const register = useCallback(async (data: RegisterRequest) => {
    const response = await authService.register(data);
    setTokens(response.accessToken, response.refreshToken);
    setUser(response.user);
    toast.success("Cuenta creada exitosamente", {
      description: `Bienvenido, ${response.user.fullName}`,
    });
    router.push("/dashboard");
  }, [router]);

  const logout = useCallback(async () => {
    const refreshToken = getRefreshToken();
    if (refreshToken) {
      try {
        await authService.logout(refreshToken);
      } catch {
        // continue with local cleanup even if API call fails
      }
    }
    clearTokens();
    setUser(null);
    toast.info("Sesión cerrada", {
      description: "Has cerrado sesión correctamente",
    });
    router.push("/login");
  }, [router]);

  const changePassword = useCallback(async (data: ChangePasswordRequest) => {
    await authService.changePassword(data);
    const refreshToken = getRefreshToken();
    if (refreshToken) {
      try {
        await authService.logout(refreshToken);
      } catch {
        // continue with local cleanup
      }
    }
    clearTokens();
    setUser(null);
    toast.success("Contraseña cambiada", {
      description: "Tu contraseña se actualizó correctamente. Inicia sesión de nuevo.",
    });
    router.push("/login");
  }, [router]);

  const getMe = useCallback(async () => {
    const userInfo = await authService.me();
    startTransition(() => setUser(userInfo));
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
        changePassword,
        getMe,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
