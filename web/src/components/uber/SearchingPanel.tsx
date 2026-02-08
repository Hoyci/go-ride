import { Car } from "lucide-react";

interface SearchingPanelProps {
  onCancel: () => void;
}

const SearchingPanel = ({ onCancel }: SearchingPanelProps) => {
  return (
    <div className="w-full bg-card rounded-t-3xl shadow-2xl p-8 text-center pb-12">
      <div className="relative w-20 h-20 mx-auto mb-6 pulse-animation flex items-center justify-center bg-primary rounded-full">
        <Car className="text-primary-foreground" size={30} />
      </div>
      <h3 className="text-xl font-bold mb-2">Procurando motoristas próximos...</h3>
      <p className="text-muted-foreground mb-6">Aguarde um momento</p>
      <div className="w-full h-2 bg-secondary rounded-full overflow-hidden">
        <div className="h-full bg-primary w-1/3 loading-bar" />
      </div>
      <button
        onClick={onCancel}
        className="mt-6 text-destructive font-bold text-sm uppercase tracking-wider cursor-pointer"
      >
        Cancelar
      </button>
    </div>
  );
};

export default SearchingPanel;
