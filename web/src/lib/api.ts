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

