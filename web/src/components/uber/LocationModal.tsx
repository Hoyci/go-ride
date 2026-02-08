import { useState } from "react";
import { ArrowLeft, ArrowRight, MapPin, Plane, Trees } from "lucide-react";
import { Input } from "@/components/ui/input";

interface LocationModalProps {
  onClose: () => void;
  onSelectDestination: (name: string) => void;
}

const suggestions = [
  { name: "Shopping Center", address: "Av. Principal, 500", icon: MapPin },
  { name: "Aeroporto Internacional", address: "Terminal 2", icon: Plane },
  { name: "Parque da Cidade", address: "Portão 3", icon: Trees },
];

const LocationModal = ({ onClose, onSelectDestination }: LocationModalProps) => {
  const [destination, setDestination] = useState("");

  const handleConfirm = () => {
    if (destination) onSelectDestination(destination);
  };

  return (
    <div className="fixed inset-0 bg-card z-50 flex flex-col">
      <div className="p-4 shadow-sm bg-card">
        <div className="relative">
          <button onClick={onClose} className="absolute left-0 top-3 p-2 cursor-pointer">
            <ArrowLeft size={22} />
          </button>
          <h2 className="text-center text-lg font-semibold py-4">Planejar viagem</h2>
        </div>

        <div className="flex gap-3 mt-2">
          <div className="flex flex-col items-center pt-3 gap-1">
            <div className="w-2 h-2 bg-muted-foreground rounded-full" />
            <div className="w-0.5 h-8 bg-border" />
            <div className="w-2 h-2 bg-foreground" />
          </div>
          <div className="flex-1 flex flex-col gap-3">
            <div className="bg-secondary p-2 rounded-lg">
              <Input
                className="bg-transparent border-none shadow-none text-sm font-medium focus-visible:ring-0 p-0 h-auto"
                defaultValue="Localização atual"
                placeholder="Local de partida"
              />
            </div>
            <div className="bg-secondary p-2 rounded-lg flex items-center gap-2">
              <Input
                className="bg-transparent border-none shadow-none text-sm font-medium focus-visible:ring-0 p-0 h-auto"
                placeholder="Para onde?"
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                autoFocus
              />
              <button
                onClick={handleConfirm}
                className="bg-primary text-primary-foreground rounded-full p-2 w-8 h-8 flex items-center justify-center shrink-0 cursor-pointer"
              >
                <ArrowRight size={14} />
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4">
        <div className="text-xs font-bold text-muted-foreground mb-2">SUGESTÕES</div>
        {suggestions.map((s) => (
          <div
            key={s.name}
            onClick={() => onSelectDestination(s.name)}
            className="flex items-center gap-4 py-4 border-b border-border cursor-pointer hover:bg-secondary"
          >
            <div className="bg-secondary p-2 rounded-full">
              <s.icon size={16} className="text-muted-foreground" />
            </div>
            <div>
              <div className="font-bold">{s.name}</div>
              <div className="text-sm text-muted-foreground">{s.address}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default LocationModal;
