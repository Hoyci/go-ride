const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8081/api/v1"

export async function apiRequest(endpoint: string, method: string, body?: any) {
  let token = localStorage.getItem("access_token");
  
  const makeRequest = (t: string | null) => fetch(`${API_URL}${endpoint}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(t ? { Authorization: `Bearer ${t}` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  let res = await makeRequest(token);

  if (res.status === 401) {
    const refreshToken = localStorage.getItem("refresh_token");
    if (refreshToken) {
      const refreshRes = await fetch(`${API_URL}/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (refreshRes.ok) {
        const { access_token, refresh_token } = (await refreshRes.json()).data;
        localStorage.setItem("access_token", access_token);
        localStorage.setItem("refresh_token", refresh_token);
        
        res = await makeRequest(access_token);
      } else {
        localStorage.clear();
        window.location.href = "/";
      }
    }
  }

  const json = await res.json();
  if (!res.ok) throw new Error(json.error?.message || "Erro na requisição");
  return json.data;
}

export async function searchAddresses(query: string): Promise<LocationResult[]> {
  if (query.length < 3) return [];
  
  const response = await fetch(
    `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&limit=5&addressdetails=1`
  );
  
  if (!response.ok) return [];
  
  const data: NominatimResponse[] = await response.json();
  
  return data.map((item) => ({
    name: item.name,
    address: item.display_name,
    lat: parseFloat(item.lat),
    lon: parseFloat(item.lon)
  }));
}

export interface LocationResult {
  name: string;
  address: string;
  lat: number;
  lon: number;
}

export interface NominatimResponse {
  name: string;
  display_name: string;
  lat: string;
  lon: string;
}

