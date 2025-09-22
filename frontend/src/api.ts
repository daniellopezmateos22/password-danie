// cliente HTTP sencillo para la API; añade el token automáticamente.
const BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export function setToken(token: string | null) {
  if (token) localStorage.setItem("token", token);
  else localStorage.removeItem("token");
}

export function getToken(): string | null {
  return localStorage.getItem("token");
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: HeadersInit = { "Content-Type": "application/json", ...(options.headers || {}) };
  const token = getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });
  if (!res.ok) {
    let msg = `HTTP ${res.status}`;
    try { const j = await res.json(); if (j.error) msg = j.error; } catch {}
    throw new Error(msg);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export const api = {
  // auth
  register: (email: string, password: string) =>
    request<{ id: number; email: string }>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
  login: (email: string, password: string) =>
    request<{ access_token: string; user: { id: number; email: string } }>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
  me: () => request<{ claims: Record<string, unknown> }>("/api/v1/users/me"),

  // vault
  createSecret: (payload: {
    username: string;
    password_plain: string;
    url?: string;
    notes?: string;
    icon?: string;
    title?: string;
  }) => request<{ id: number }>("/api/v1/vault/entries", { method: "POST", body: JSON.stringify(payload) }),

  listSecrets: (q = "", limit = 20, offset = 0) =>
    request<{ items: any[]; total: number }>(`/api/v1/vault/entries?q=${encodeURIComponent(q)}&limit=${limit}&offset=${offset}`),

  getSecret: (id: number) => request<any>(`/api/v1/vault/entries/${id}`),

  updateSecret: (id: number, payload: Partial<{
    username: string;
    password_plain: string;
    url: string;
    notes: string;
    icon: string;
    title: string;
  }>) => request(`/api/v1/vault/entries/${id}`, { method: "PUT", body: JSON.stringify(payload) }),

  deleteSecret: (id: number) => request(`/api/v1/vault/entries/${id}`, { method: "DELETE" }),
};
