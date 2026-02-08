import { useEffect, useRef } from "react";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

interface MapViewProps {
  onMapReady: (map: L.Map) => void;
  userCoords: [number, number] | null;
}

const MapView = ({ onMapReady, userCoords }: MapViewProps) => {
  const userMarkerRef = useRef<L.Marker | null>(null);
  const accuracyCircleRef = useRef<L.Circle | null>(null);
  const mapRef = useRef<L.Map | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const userIcon = L.divIcon({
    className: "",
    html: `<div style="background-color: hsl(217, 91%, 60%); width: 16px; height: 16px; border-radius: 50%; border: 3px solid white; box-shadow: 0 0 10px rgba(0,0,0,0.3);"></div>`,
    iconSize: [20, 20],
    iconAnchor: [10, 10],
  });

const updateUserMarker = (coords: [number, number]) => {
    if (!mapRef.current) return;

    const latLng = L.latLng(coords[0], coords[1]);

    if (userMarkerRef.current) {
      userMarkerRef.current.setLatLng(latLng);
    } else {
      userMarkerRef.current = L.marker(latLng, { icon: userIcon }).addTo(mapRef.current);
    }

    if (accuracyCircleRef.current) {
      accuracyCircleRef.current.setLatLng(latLng);
    } else {
      accuracyCircleRef.current = L.circle(latLng, {
        radius: 30, // metros
        fillColor: "hsl(217, 91%, 60%)",
        fillOpacity: 0.15,
        color: "transparent",
        interactive: false,
      }).addTo(mapRef.current);
    }
  };

  useEffect(() => {
    if (!mapRef.current && containerRef.current && userCoords) {
      const map = L.map(containerRef.current, { 
        zoomControl: false,
        fadeAnimation: true 
      }).setView(userCoords, 16);

      L.tileLayer("https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png", {
        attribution: '&copy; OpenStreetMap',
        subdomains: "abcd",
        maxZoom: 20,
      }).addTo(map);

      mapRef.current = map;
      onMapReady(map);
    }
  }, [userCoords]);

  useEffect(() => {
    if (userCoords && mapRef.current) {
      updateUserMarker(userCoords);
      
      mapRef.current.panTo(userCoords); 
    }
  }, [userCoords]);

  return <div ref={containerRef} className="h-screen w-full z-0" />;
};

export default MapView;
