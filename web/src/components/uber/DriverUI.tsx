import { useState, useEffect, useRef } from "react";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { toast } from "sonner";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8081/api/v1";
const WS_BASE_URL = API_URL.replace(/^http/, "ws");

interface DriverUIProps {
  userCoords: [number, number] | null;
}

const DriverUI = ({ userCoords }: DriverUIProps) => {
  const [online, setOnline] = useState(() => {
    return localStorage.getItem("driver_online_status") === "true";
  });
  const [showRequests, setShowRequests] = useState(false);

  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    localStorage.setItem("driver_online_status", String(online));
  }, [online]);

  useEffect(() => {
    if (online) {
      const token = localStorage.getItem("access_token");
      if (!token) {
        toast.error("Erro de autenticação. Faça login novamente.");
        setOnline(false);
        return;
      }

      const wsUrl = `${WS_BASE_URL}/driver/stream?token=${token}`;
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        toast.success("Você está online e visível no mapa!");
        // setTimeout(() => setShowRequests(true), 2000);
      };

      ws.onclose = () => {
        console.log("WebSocket desconectado");
        setOnline(false);
        setShowRequests(false);
        wsRef.current = null;
      };

      ws.onerror = (error) => {
        console.error("Erro no WebSocket:", error);
        toast.error("Erro de conexão em tempo real.");
        setOnline(false);
      };

    } else {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      setShowRequests(false);
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [online]);

  // 3. Efeito do Tracking (Reativo ao movimento + Heartbeat de segurança)
  useEffect(() => {
    if (!online || !userCoords || !wsRef.current) return;

    const sendLocation = () => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        const payload = {
          latitude: userCoords[0],
          longitude: userCoords[1],
        };
        wsRef.current.send(JSON.stringify(payload));
        console.log("📍 Dado enviado via WS:", payload);
      }
    };

    sendLocation();

    const heartbeatInterval = setInterval(() => {
      sendLocation();
    }, 10000);

    return () => {
      clearInterval(heartbeatInterval);
    };
  }, [userCoords, online]);

  const handleToggle = (checked: boolean) => {
    if (checked && !userCoords) {
      toast.error("Aguardando localização do GPS...");
      return;
    }
    setOnline(checked);
  };

  return (
    <div className="absolute bottom-0 w-full z-20">
      <div className="bg-card rounded-t-3xl shadow-2xl p-6 pb-12">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h2
              className={cn(
                "text-2xl font-bold",
                online && "text-[hsl(var(--uber-green))]"
              )}
            >
              {online ? "Online" : "Offline"}
            </h2>

            <p className="text-sm text-muted-foreground">
              {online
                ? "Buscando corridas próximas"
                : "Fique online para aceitar corridas"}
            </p>
          </div>

          <Switch
            checked={online}
            onCheckedChange={handleToggle}
          />
        </div>

        {showRequests && (
          <div className="border-t border-border pt-4">
            <h3 className="font-bold text-muted-foreground text-sm mb-4">
              SOLICITAÇÕES PRÓXIMAS
            </h3>

            <div className="bg-secondary border border-border p-4 rounded-xl flex justify-between items-center animate-pulse">
              <div>
                <div className="flex items-center gap-2 mb-1">
                  <span className="bg-primary text-primary-foreground text-xs px-2 py-0.5 rounded">
                    UberX
                  </span>
                  <span className="font-bold">R$ 24,50</span>
                </div>
                <div className="text-sm text-muted-foreground">
                  2.5 km • 8 min de distância
                </div>
              </div>

              <Button
                size="sm"
                className="font-bold"
                onClick={() =>
                  alert("Funcionalidade de aceitar corrida simulada.")
                }
              >
                Aceitar
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DriverUI;