const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081/api/v1";

interface RequestOptions {
  method?: string;
  body?: unknown;
  skipAuth?: boolean;
}

interface ApiResponseWrapper<T> {
  success: boolean;
  data: T | null;
  message?: string;
  error?: { code: string; message: string };
}

async function getAccessToken(): Promise<string | null> {
  const { getAccessToken } = await import("@/features/auth/utils/token-storage");
  return getAccessToken();
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, skipAuth = false } = options;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (!skipAuth) {
    const token = await getAccessToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
  }

  const response = await fetch(`${API_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const json: ApiResponseWrapper<T> = await response.json();

  if (!json.success) {
    const errMsg = json.error?.message || json.message || "Error del servidor";
    const errCode = json.error?.code || "UNKNOWN";
    if (response.status === 401 && !skipAuth) {
      const { clearTokens } = await import("@/features/auth/utils/token-storage");
      clearTokens();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    }
    throw { status: response.status, code: errCode, message: errMsg } as import("@/features/auth/types/auth.types").ApiError;
  }

  return json.data as T;
}

export async function apiRequestRaw<T>(path: string, options: RequestOptions = {}): Promise<ApiResponseWrapper<T>> {
  const { method = "GET", body, skipAuth = false } = options;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (!skipAuth) {
    const token = await getAccessToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
  }

  const response = await fetch(`${API_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const json: ApiResponseWrapper<T> = await response.json();

  if (!json.success) {
    if (response.status === 401 && !skipAuth) {
      const { clearTokens } = await import("@/features/auth/utils/token-storage");
      clearTokens();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    }
    throw new Error(json.error?.message || json.message || "Error del servidor");
  }

  return json;
}
