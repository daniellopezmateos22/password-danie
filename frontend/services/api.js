const API_BASE =
  (typeof window !== "undefined" && window.API_URL) ||
  import.meta.env.VITE_API_URL ||
  "http://localhost:8080";

export async function api(path, opts = {}) {
  const res = await fetch(`${API_BASE}${path}`, opts);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}
