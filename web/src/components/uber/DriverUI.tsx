import { useState, useEffect } from "react";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const DriverUI = () => {
  const [online, setOnline] = useState(false);
  const [showRequests, setShowRequests] = useState(false);

  useEffect(() => {
    if (online) {
      const timer = setTimeout(() => setShowRequests(true), 2000);
      return () => clearTimeout(timer);
    } else {
      setShowRequests(false);
    }
  }, [online]);

  return (
    <div className="absolute bottom-0 w-full z-20">
      <div className="bg-card rounded-t-3xl shadow-2xl p-6 pb-12">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h2 className={cn("text-2xl font-bold", online && "text-[hsl(var(--uber-green))]")}>
              {online ? "Online" : "Offline"}
            </h2>
            <p className="text-sm text-muted-foreground">
              {online ? "Buscando corridas próximas" : "Fique online para aceitar corridas"}
            </p>
          </div>
          <Switch checked={online} onCheckedChange={setOnline} />
        </div>

        {showRequests && (
          <div className="border-t border-border pt-4">
            <h3 className="font-bold text-muted-foreground text-sm mb-4">SOLICITAÇÕES PRÓXIMAS</h3>
            <div className="bg-secondary border border-border p-4 rounded-xl flex justify-between items-center animate-pulse">
              <div>
                <div className="flex items-center gap-2 mb-1">
                  <span className="bg-primary text-primary-foreground text-xs px-2 py-0.5 rounded">UberX</span>
                  <span className="font-bold">R$ 24,50</span>
                </div>
                <div className="text-sm text-muted-foreground">2.5 km • 8 min de distância</div>
              </div>
              <Button
                size="sm"
                onClick={() => alert("Funcionalidade de aceitar corrida simulada.")}
                className="font-bold"
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
