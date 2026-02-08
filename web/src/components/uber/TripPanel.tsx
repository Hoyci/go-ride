import { User, Star } from "lucide-react";
import { Button } from "@/components/ui/button";

interface TripPanelProps {
  onFinish: () => void;
}

const TripPanel = ({ onFinish }: TripPanelProps) => {
  return (
    <div className="w-full bg-card rounded-t-3xl shadow-2xl p-6 pb-8">
      <div className="flex items-center justify-between mb-4 border-b border-border pb-4">
        <div className="text-sm font-bold bg-primary text-primary-foreground px-2 py-1 rounded">
          Em andamento
        </div>
        <div className="text-sm text-muted-foreground">Previsão: 12 min</div>
      </div>

      <div className="flex items-center gap-4 mb-6">
        <div className="w-16 h-16 bg-secondary rounded-full overflow-hidden flex items-center justify-center">
          <User size={28} className="text-muted-foreground" />
        </div>
        <div className="flex-1">
          <h3 className="font-bold text-lg">Carlos Silva</h3>
          <p className="text-muted-foreground text-sm">Fiat Argo • ABC-1234</p>
          <div className="flex items-center gap-1 text-sm mt-1">
            <Star size={14} className="text-[hsl(var(--uber-yellow))] fill-[hsl(var(--uber-yellow))]" /> 4.9
          </div>
        </div>
        <img
          src="https://www.uber-assets.com/image/upload/f_auto,q_auto:eco,c_fill,w_956,h_537/v1568070387/assets/b5/0a5191-836e-42bf-8ea1-b8f219582d0d/original/UberX.png"
          className="w-16 object-contain"
          alt="Carro"
        />
      </div>

      <div className="flex gap-3">
        <Button variant="outline" className="flex-1 py-3 font-semibold">
          Mensagem
        </Button>
        <Button onClick={onFinish} className="flex-1 py-3 font-semibold">
          Finalizar (Sim.)
        </Button>
      </div>
    </div>
  );
};

export default TripPanel;
