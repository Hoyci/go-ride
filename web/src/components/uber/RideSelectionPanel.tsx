import { useState } from "react";
import { User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface RideSelectionPanelProps {
  onConfirm: () => void;
  onBack: () => void;
}

const rides = [
  {
    id: "uber-x",
    name: "UberX",
    seats: 4,
    eta: "14:25",
    price: "R$ 19,90",
    img: "https://www.uber-assets.com/image/upload/f_auto,q_auto:eco,c_fill,w_956,h_537/v1568070387/assets/b5/0a5191-836e-42bf-8ea1-b8f219582d0d/original/UberX.png",
  },
  {
    id: "uber-black",
    name: "Black",
    seats: 4,
    eta: "14:28",
    price: "R$ 32,50",
    img: "https://www.uber-assets.com/image/upload/f_auto,q_auto:eco,c_fill,w_956,h_537/v1568134115/assets/6d/354919-18b0-45d0-a151-501ab4c4b114/original/UberBlack.png",
  },
];

const RideSelectionPanel = ({ onConfirm }: RideSelectionPanelProps) => {
  const [selected, setSelected] = useState("uber-x");

  const selectedRide = rides.find((r) => r.id === selected);

  return (
    <div className="w-full bg-card rounded-t-3xl shadow-2xl p-6 pb-8">
      <div className="w-12 h-1 bg-border rounded-full mx-auto mb-4" />
      <h3 className="text-center font-bold mb-4 text-muted-foreground text-sm">ESCOLHA UMA VIAGEM</h3>

      <div className="space-y-2 mb-6 max-h-60 overflow-y-auto">
        {rides.map((ride) => (
          <div
            key={ride.id}
            onClick={() => setSelected(ride.id)}
            className={cn(
              "p-3 rounded-xl flex items-center justify-between cursor-pointer border-2 transition",
              selected === ride.id
                ? "border-primary bg-secondary"
                : "border-transparent hover:bg-secondary"
            )}
          >
            <div className="flex items-center gap-3">
              <img src={ride.img} className="w-16 h-10 object-contain" alt={ride.name} />
              <div>
                <div className="font-bold flex items-center gap-1">
                  {ride.name} <User size={12} /> {ride.seats}
                </div>
                <div className="text-xs text-muted-foreground">{ride.eta} de chegada</div>
              </div>
            </div>
            <div className="font-bold text-lg">{ride.price}</div>
          </div>
        ))}
      </div>

      <Button onClick={onConfirm} className="w-full py-4 text-lg font-bold rounded-lg h-auto">
        Confirmar {selectedRide?.name}
      </Button>
    </div>
  );
};

export default RideSelectionPanel;
