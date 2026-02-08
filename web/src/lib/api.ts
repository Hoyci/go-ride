const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8081/api/v1"

export async function apiRequest(endpoint: string, method: string, body?: any) {
  const token = localStorage.getItem("access_token");
  const res = await fetch(`${API_URL}${endpoint}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  const json = await res.json();
  if (!res.ok) throw new Error(json.error?.message || "Erro na requisição");
  return json.data;
}