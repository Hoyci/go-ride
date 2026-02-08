import { useState } from "react";
import { Button } from "@/components/ui/button";

interface RideSelectionPanelProps {
  fares: any[];
  onConfirm: () => void;
  onBack: () => void;
}


const RideSelectionPanel = ({ fares, onConfirm }: RideSelectionPanelProps) => {
  const [selectedId, setSelectedId] = useState(fares[0]?.id);

  return (
    <div className="w-full bg-card rounded-t-3xl shadow-2xl p-6 pb-8">
      <h3 className="text-center font-bold mb-4 text-muted-foreground text-sm uppercase">Opções de Viagem</h3>
      <div className="space-y-2 mb-6">
        {fares.map((fare) => (
          <div
            key={fare.id}
            onClick={() => setSelectedId(fare.id)}
            className={`p-4 rounded-xl flex justify-between items-center border-2 transition cursor-pointer ${selectedId === fare.id ? "border-black bg-secondary" : "border-transparent"
              }`}
          >
            <div className="flex items-center gap-3">
              <img
                src={fare.packageSlug === 1 ? "/uberx.png" : "/black.png"}
                alt={fare.packageSlug === 1 ? "UberX" : "Black"}
                className="w-12"
              />
              <div>
                <p className="font-bold">{fare.packageSlug === 1 ? "UberX" : "Uber Black"}</p>
                <p className="text-xs text-muted-foreground">Melhor preço</p>
              </div>
            </div>
            <p className="font-bold text-lg">
              {(fare.totalPriceInCents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}
            </p>
          </div>
        ))}
      </div>
      <Button onClick={onConfirm} className="w-full py-4 font-bold">Confirmar Viagem</Button>
    </div>
  );
};

export default RideSelectionPanel;
