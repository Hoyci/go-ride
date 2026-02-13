import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useUser } from "@/contexts/UserContext";
import { UserCircle, LogOut, Star, Clock, Settings, HelpCircle, MapPinOff } from "lucide-react";
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
import { useMutation } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api";
import { toast } from "sonner"

const Dashboard = () => {
  const [locationStatus, setLocationStatus] = useState<"loading" | "granted" | "denied">("loading");
  const [coords, setCoords] = useState<[number, number] | null>(null)
  const [map, setMap] = useState<L.Map | null>(null);
  const { user, logout } = useUser();
  const navigate = useNavigate();

  const logoutMutation = useMutation({
    mutationFn: () => apiRequest("/logout", "POST",),
    onSuccess: () => {
      logout();
      navigate("/");
    },
    onError: (error: Error) => {
      toast.error("Credenciais inválidas ou erro de conexão.");
    }
  });

  const handleLogout = (e: React.FormEvent) => {
    e.preventDefault();
    logoutMutation.mutate()
  };

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (!user && !token) {
      navigate("/");
      return;
    }

    const geo = navigator.geolocation;
    if (!geo) {
      setLocationStatus("denied");
      return;
    }

    const geoOptions = {
      enableHighAccuracy: true, // Força o uso de GPS/Hardware
      timeout: 10000,           // Tempo limite de 10 segundos
      maximumAge: 0             // Não aceita localização em cache (antiga)
    };

    const watchId = geo.watchPosition(
      (pos) => {
        setCoords([pos.coords.latitude, pos.coords.longitude]);
        setLocationStatus("granted");
      },
      (err) => {
        console.error(err);
        setLocationStatus("denied");
      },
      geoOptions
    );

    return () => geo.clearWatch(watchId);
  }, [user, navigate]);


  if (!user) return null;

  if (locationStatus === "denied") {
    return (
      <div className="h-screen w-screen flex flex-col items-center justify-center bg-background p-6 text-center">
        <div className="bg-destructive/10 p-4 rounded-full mb-4">
          <MapPinOff size={48} className="text-destructive" />
        </div>
        <h1 className="text-2xl font-bold mb-2">Acesso à localização necessário</h1>
        <p className="text-muted-foreground mb-6 max-w-sm">
          Para usar o GoRide, precisamos saber onde você está. Por favor, habilite a localização nas configurações do seu navegador e recarregue a página.
        </p>
        <button
          onClick={() => window.location.reload()}
          className="bg-primary text-primary-foreground px-6 py-2 rounded-lg font-bold"
        >
          Tentar novamente
        </button>
      </div>
    );
  }

  if (locationStatus === "loading") {
    return (
      <div className="h-screen w-screen flex items-center justify-center bg-background">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="relative h-screen w-screen">
      <MapView onMapReady={setMap} userCoords={coords} />
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
              {user.type === "PASSENGER" ? "Passageiro" : "Motorista"}
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

      {user.type === "PASSENGER" ? <PassengerUI map={map} userCoords={coords} /> : <DriverUI />}
    </div>
  );
};

export default Dashboard;
