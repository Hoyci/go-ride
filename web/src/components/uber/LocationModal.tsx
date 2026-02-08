import { useState, useEffect } from "react";
import { ArrowLeft, MapPin, Loader2, Circle, Square } from "lucide-react";
import { Input } from "@/components/ui/input";
import { searchAddresses } from "@/lib/api";

export interface LocationResult {
  name: string;
  address: string;
  lat: number;
  lon: number;
}

interface LocationModalProps {
  onClose: () => void;
  onConfirmSelection: (pickup: LocationResult, destination: LocationResult) => void;
  initialPickup?: LocationResult | null;
}

const LocationModal = ({ onClose, onConfirmSelection, initialPickup }: LocationModalProps) => {
  const [pickup, setPickup] = useState<LocationResult | null>(initialPickup || null);
  const [destination, setDestination] = useState<LocationResult | null>(null);

  const [activeField, setActiveField] = useState<"pickup" | "destination">("destination");
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<LocationResult[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const delayDebounce = setTimeout(async () => {
      if (query.length > 2) {
        setLoading(true);
        try {
          const data = await searchAddresses(query);
          setResults(data);
        } catch (error) {
          console.error("Erro ao buscar endereços:", error);
          setResults([]);
        } finally {
          setLoading(false);
        }
      } else {
        setResults([]);
      }
    }, 500);

    return () => clearTimeout(delayDebounce);
  }, [query]);

  const handleSelectLocation = (res: LocationResult) => {
    if (activeField === "pickup") {
      setPickup(res);
      setActiveField("destination");
      setQuery("");
      setResults([]);
    } else {
      setDestination(res);
      const finalPickup = pickup || initialPickup;
      
      if (finalPickup) {
        onConfirmSelection(finalPickup, res);
      }
    }
  };

  return (
    <div className="fixed inset-0 bg-card z-50 flex flex-col animate-in slide-in-from-bottom duration-300">
      {/* Header e Inputs */}
      <div className="p-4 shadow-sm bg-card border-b">
        <div className="flex items-center mb-4">
          <button 
            onClick={onClose} 
            className="p-2 -ml-2 hover:bg-secondary rounded-full transition-colors"
          >
            <ArrowLeft size={22} />
          </button>
          <h2 className="flex-1 text-center text-lg font-semibold pr-8">Planejar viagem</h2>
        </div>

        <div className="flex gap-3 relative">
          {/* Indicador Visual de Rota */}
          <div className="flex flex-col items-center pt-3 gap-1">
            <Circle size={8} className="fill-muted-foreground text-muted-foreground" />
            <div className="w-0.5 h-10 bg-border" />
            <Square size={8} className="fill-foreground text-foreground" />
          </div>

          <div className="flex-1 flex flex-col gap-3">
            {/* Input Ponto de Partida */}
            <div className={`bg-secondary p-2 rounded-lg flex items-center gap-2 border-2 transition-all ${activeField === "pickup" ? "border-primary" : "border-transparent"}`}>
              <Input
                className="bg-transparent border-none shadow-none text-sm font-medium focus-visible:ring-0 p-0 h-8"
                placeholder="Ponto de partida"
                value={activeField === "pickup" ? query : (pickup?.name || "Localização atual")}
                onFocus={() => {
                  setActiveField("pickup");
                  setQuery("");
                  setResults([]);
                }}
                onChange={(e) => setQuery(e.target.value)}
              />
              {loading && activeField === "pickup" && <Loader2 size={16} className="animate-spin text-muted-foreground mr-2" />}
            </div>

            {/* Input Destino */}
            <div className={`bg-secondary p-2 rounded-lg flex items-center gap-2 border-2 transition-all ${activeField === "destination" ? "border-primary" : "border-transparent"}`}>
              <Input
                className="bg-transparent border-none shadow-none text-sm font-medium focus-visible:ring-0 p-0 h-8"
                placeholder="Para onde?"
                value={activeField === "destination" ? query : (destination?.name || "")}
                onFocus={() => {
                  setActiveField("destination");
                  setQuery("");
                  setResults([]);
                }}
                onChange={(e) => setQuery(e.target.value)}
                autoFocus
              />
              {loading && activeField === "destination" && <Loader2 size={16} className="animate-spin text-muted-foreground mr-2" />}
            </div>
          </div>
        </div>
      </div>

      {/* Lista de Resultados */}
      <div className="flex-1 overflow-y-auto bg-background">
        {results.length > 0 ? (
          results.map((res, index) => (
            <div
              key={index}
              onClick={() => handleSelectLocation(res)}
              className="flex items-center gap-4 px-6 py-4 border-b border-border cursor-pointer hover:bg-secondary transition-colors"
            >
              <div className="bg-secondary p-2 rounded-full shrink-0">
                <MapPin size={18} className="text-muted-foreground" />
              </div>
              <div className="overflow-hidden">
                <div className="font-bold text-sm truncate">{res.name}</div>
                <div className="text-xs text-muted-foreground truncate">{res.address}</div>
              </div>
            </div>
          ))
        ) : query.length > 2 && !loading ? (
          <div className="p-8 text-center text-muted-foreground text-sm italic">
            Nenhum local encontrado para "{query}"
          </div>
        ) : (
          <div className="p-8 text-center text-muted-foreground text-xs uppercase tracking-widest opacity-50">
            Digite pelo menos 3 caracteres
          </div>
        )}
      </div>
    </div>
  );
};

export default LocationModal;