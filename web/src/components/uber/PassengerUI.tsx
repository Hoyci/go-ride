import { useState, useEffect, useCallback } from "react";
import L from "leaflet";
import SearchPanel from "./SearchPanel";
import LocationModal from "./LocationModal";
import RideSelectionPanel from "./RideSelectionPanel";
import SearchingPanel from "./SearchingPanel";
import TripPanel from "./TripPanel";
import { START_COORDS } from "./MapView";
import { toast } from "sonner";

type PassengerStep = "search" | "selecting" | "searching" | "trip";

interface PassengerUIProps {
  map: L.Map | null;
}

const PassengerUI = ({ map }: PassengerUIProps) => {
  const [step, setStep] = useState<PassengerStep>("search");
  const [showModal, setShowModal] = useState(false);
  const [destMarker, setDestMarker] = useState<L.Marker | null>(null);
  const [routeLine, setRouteLine] = useState<L.Polyline | null>(null);

  const clearMap = useCallback(() => {
    if (destMarker && map) map.removeLayer(destMarker);
    if (routeLine && map) map.removeLayer(routeLine);
    setDestMarker(null);
    setRouteLine(null);
  }, [destMarker, routeLine, map]);

  const handleSelectDestination = (name: string) => {
    setShowModal(false);
    if (!map) return;

    const center = map.getCenter();
    const latOff = (Math.random() - 0.5) * 0.02;
    const lngOff = (Math.random() - 0.5) * 0.02;
    const destCoords: [number, number] = [center.lat + latOff, center.lng + lngOff];

    // Clear previous
    clearMap();

    // Destination marker
    const destIcon = L.divIcon({
      className: "",
      html: `<div style="background-color: hsl(0,0%,0%); width: 14px; height: 14px;"></div><div style="width: 2px; height: 10px; background: black; margin: 0 auto;"></div>`,
      iconSize: [20, 20],
      iconAnchor: [10, 20],
    });

    const marker = L.marker(destCoords, { icon: destIcon }).addTo(map);
    setDestMarker(marker);

    // Route line
    const line = L.polyline([START_COORDS, destCoords], {
      color: "black",
      weight: 4,
      opacity: 0.8,
      dashArray: "10, 10",
    }).addTo(map);
    setRouteLine(line);

    map.fitBounds(L.latLngBounds([START_COORDS, destCoords]), { padding: [50, 50] });

    setStep("selecting");
  };

  const handleConfirmRide = () => {
    setStep("searching");
    setTimeout(() => setStep("trip"), 3000);
  };

  const handleCancelRide = () => {
    setStep("selecting");
  };

  const handleFinishTrip = () => {
    toast.success("Viagem finalizada!");
    clearMap();
    map?.setView(START_COORDS, 15);
    setStep("search");
  };

  return (
    <>
      {showModal && (
        <LocationModal
          onClose={() => setShowModal(false)}
          onSelectDestination={handleSelectDestination}
        />
      )}

      <div className="absolute bottom-0 w-full z-20 flex flex-col items-center">
        {step === "search" && <SearchPanel onOpenModal={() => setShowModal(true)} />}
        {step === "selecting" && (
          <RideSelectionPanel onConfirm={handleConfirmRide} onBack={() => setStep("search")} />
        )}
        {step === "searching" && <SearchingPanel onCancel={handleCancelRide} />}
        {step === "trip" && <TripPanel onFinish={handleFinishTrip} />}
      </div>
    </>
  );
};

export default PassengerUI;
