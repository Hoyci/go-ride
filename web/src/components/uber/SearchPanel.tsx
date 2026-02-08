import { Search, Home, Briefcase } from "lucide-react";

interface SearchPanelProps {
  onOpenModal: () => void;
}

const SearchPanel = ({ onOpenModal }: SearchPanelProps) => {
  return (
    <div className="w-full bg-card rounded-t-3xl shadow-2xl p-6 pb-8">
      <h3 className="text-xl font-bold mb-4">Para onde vamos?</h3>

      <div
        onClick={onOpenModal}
        className="bg-secondary p-4 rounded-xl flex items-center gap-3 cursor-pointer hover:bg-accent transition"
      >
        <Search size={20} />
        <div className="flex-1">
          <div className="font-bold text-lg">Buscar destino</div>
        </div>
      </div>

      <div className="mt-6">
        <div className="flex items-center gap-4 py-3 border-b border-border">
          <div className="bg-secondary p-2 rounded-full">
            <Home size={16} className="text-muted-foreground" />
          </div>
          <div>
            <div className="font-bold">Casa</div>
            <div className="text-sm text-muted-foreground">Rua das Flores, 123</div>
          </div>
        </div>
        <div className="flex items-center gap-4 py-3">
          <div className="bg-secondary p-2 rounded-full">
            <Briefcase size={16} className="text-muted-foreground" />
          </div>
          <div>
            <div className="font-bold">Trabalho</div>
            <div className="text-sm text-muted-foreground">Av. Paulista, 1000</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SearchPanel;
