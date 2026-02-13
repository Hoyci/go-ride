import { useState, useCallback } from "react";
import L from "leaflet";
import SearchPanel from "./SearchPanel";
import LocationModal, { LocationResult } from "./LocationModal";
import RideSelectionPanel from "./RideSelectionPanel";
import SearchingPanel from "./SearchingPanel";
import TripPanel from "./TripPanel";
import { toast } from "sonner";
import { apiRequest } from "@/lib/api";
import { useUser } from "@/contexts/UserContext";
import { useMutation } from "@tanstack/react-query";

type PassengerStep = "search" | "selecting" | "searching" | "trip";

interface PassengerUIProps {
  map: L.Map | null;
  userCoords: [number, number] | null;
}

const PassengerUI = ({ map, userCoords }: PassengerUIProps) => {
  const { user } = useUser();
  const [step, setStep] = useState<PassengerStep>("search");
  const [showModal, setShowModal] = useState(false);
  const [rideFares, setRideFares] = useState<any[]>([])
  const [selectedRideFareID, setSelectedRideFareID] = useState(null);

  const [pickupMarker, setPickupMarker] = useState<L.Marker | null>(null);
  const [destMarker, setDestMarker] = useState<L.Marker | null>(null);
  const [routeLine, setRouteLine] = useState<L.Polyline | null>(null);

  const pickupIcon = L.divIcon({
    className: "",
    html: `<div style="background-color: white; width: 12px; height: 12px; border: 3px solid black; border-radius: 50%;"></div>`,
    iconSize: [12, 12],
    iconAnchor: [6, 6],
  });

  const destinationIcon = L.divIcon({
    className: "",
    html: `<div style="background-color: black; width: 12px; height: 12px; border: 2px solid white;"></div>`,
    iconSize: [12, 12],
    iconAnchor: [6, 6],
  });

  const clearMap = useCallback(() => {
    if (map) {
      if (pickupMarker) map.removeLayer(pickupMarker);
      if (destMarker) map.removeLayer(destMarker);
      if (routeLine) map.removeLayer(routeLine);
    }
    setPickupMarker(null);
    setDestMarker(null);
    setRouteLine(null);
  }, [pickupMarker, destMarker, routeLine, map]);

  const handleConfirmSelection = async (pickup: LocationResult, destination: LocationResult) => {
    setShowModal(false);
    if (!map) return;

    try {
      const response = await apiRequest("/trip-preview", "POST", {
        passenger_id: user.id,
        origin: {
          latitude: pickup.lat,
          longitude: pickup.lon
        },
        destination: {
          latitude: destination.lat,
          longitude: destination.lon
        }
      });

      const { route, rideFares } = response;
      setRideFares(rideFares);

      const routePoints: [number, number][] = route.geometry[0].coordinates.map((coord: any) => [
        coord.longitude,
        coord.latitude
      ]);

      clearMap();

      const pMarker = L.marker([pickup.lat, pickup.lon], { icon: pickupIcon }).addTo(map);
      const dMarker = L.marker([destination.lat, destination.lon], { icon: destinationIcon }).addTo(map);

      const line = L.polyline(routePoints, {
        color: "black",
        weight: 5,
        opacity: 0.8,
        lineJoin: 'round'
      }).addTo(map);

      setPickupMarker(pMarker);
      setDestMarker(dMarker);
      setRouteLine(line);

      map.fitBounds(line.getBounds(), { padding: [50, 50] });

      setStep("selecting");
    } catch (error) {
      toast.error("Erro ao calcular rota. Tente novamente.");
      console.error(error);
    }
  };

    const createTripMutation = useMutation({
      mutationFn: (data: { ride_fare_id: string, user_id: string}) => apiRequest("/trip", "POST", data),
      onSuccess: () => {
        setStep("searching");
      },
      onError: (error: Error) => {
        toast.error("Erro ao criar corrida, tente novamente.");
      }
    });

    const handleConfirmRide = (rideFareId: string) => {
      createTripMutation.mutate({ride_fare_id: rideFareId, user_id: user.id})
    }

    // Quando receber o sinal de motorista encontrado via websocket, disparar o toast abaixo
    //     setTimeout(() => {
    //   toast.success("Motorista encontrado!");
    //   setStep("trip");
    // }, 3000);

  const handleCancelRide = () => {
    setStep("selecting");
  };

  const handleFinishTrip = () => {
    toast.success("Viagem finalizada!");
    clearMap();
    if (map && userCoords) {
      map.setView(userCoords, 16);
    }
    setStep("search");
  };

  return (
    <>
      {showModal && (
        <LocationModal
          onClose={() => setShowModal(false)}
          onConfirmSelection={handleConfirmSelection}
          initialPickup={userCoords ? {
            name: "Minha localização",
            address: "Localização atual",
            lat: userCoords[0],
            lon: userCoords[1]
          } : null}
        />
      )}

      <div className="absolute bottom-0 w-full z-20 flex flex-col items-center">
        {step === "search" && (
          <SearchPanel onOpenModal={() => setShowModal(true)} />
        )}

        {step === "selecting" && (
          <RideSelectionPanel
            fares={rideFares}
            selectedFareID={selectedRideFareID}
            setSelectedFareID={setSelectedRideFareID}
            onConfirm={handleConfirmRide}
            onBack={() => {
              clearMap();
              setStep("search");
            }}
          />
        )}

        {step === "searching" && (
          <SearchingPanel onCancel={handleCancelRide} />
        )}

        {step === "trip" && (
          <TripPanel onFinish={handleFinishTrip} />
        )}
      </div>
    </>
  );
};

export default PassengerUI;