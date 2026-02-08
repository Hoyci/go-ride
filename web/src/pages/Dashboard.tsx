import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useUser } from "@/contexts/UserContext";
import { Menu, UserCircle, LogOut, Star, Clock, Settings, HelpCircle } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import MapView from "@/components/uber/MapView";
import PassengerUI from "@/components/uber/PassengerUI";
import DriverUI from "@/components/uber/DriverUI";
import L from "leaflet";

const Dashboard = () => {
  const { user, logout } = useUser();
  const navigate = useNavigate();
  const [map, setMap] = useState<L.Map | null>(null);

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  useEffect(() => {
    if (!user) navigate("/");
  }, [user, navigate]);

  if (!user) return null;

  return (
    <div className="relative h-screen w-screen">
      <MapView onMapReady={setMap} />

      {/* Header */}
      <div className="absolute top-4 left-4 z-20">
        <button className="bg-card p-3 rounded-full shadow-lg hover:bg-secondary cursor-pointer">
          <Menu size={22} />
        </button>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button className="absolute top-4 right-4 z-20 bg-card px-4 py-2 rounded-full shadow-lg font-semibold flex items-center gap-2 cursor-pointer hover:bg-secondary transition">
            <UserCircle size={22} className="text-muted-foreground" />
            <span>{user.name.split(" ")[0]}</span>
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56 z-50">
          <DropdownMenuLabel>
            <div className="font-bold">{user.name}</div>
            <div className="text-xs text-muted-foreground font-normal">{user.email}</div>
            <div className="text-xs text-muted-foreground font-normal capitalize mt-0.5">
              {user.type === "passenger" ? "Passageiro" : "Motorista"}
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem className="cursor-pointer">
            <Clock size={16} className="mr-2" /> Minhas viagens
          </DropdownMenuItem>
          <DropdownMenuItem className="cursor-pointer">
            <Star size={16} className="mr-2" /> Avaliações
          </DropdownMenuItem>
          <DropdownMenuItem className="cursor-pointer">
            <Settings size={16} className="mr-2" /> Configurações
          </DropdownMenuItem>
          <DropdownMenuItem className="cursor-pointer">
            <HelpCircle size={16} className="mr-2" /> Ajuda
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={handleLogout} className="cursor-pointer text-destructive focus:text-destructive">
            <LogOut size={16} className="mr-2" /> Sair
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Conditional UI */}
      {user.type === "passenger" ? <PassengerUI map={map} /> : <DriverUI />}
    </div>
  );
};

export default Dashboard;
