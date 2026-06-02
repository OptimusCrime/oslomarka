import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON, Polyline, Marker, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import type { Coord, PathResult } from '../types';
import { decodePolyline } from '../lib/polyline';

// Oslo city centre
const OSLO_CENTER: [number, number] = [59.92, 10.75];
const DEFAULT_ZOOM = 11;

const ROUTE_COLORS: Record<string, string> = {
  Fotrute: '#e63946',
};

// Avoids the broken-image issue with Leaflet's default marker icons in Vite.
function makePin(color: string): L.DivIcon {
  return L.divIcon({
    className: '',
    html: `<div style="
      width:18px;height:18px;border-radius:50%;
      background:${color};border:3px solid white;
      box-shadow:0 2px 8px rgba(0,0,0,0.35);
    "></div>`,
    iconSize: [18, 18],
    iconAnchor: [9, 9],
  });
}

const START_PIN = makePin('#16a34a');
const END_PIN   = makePin('#dc2626');

interface MapClickHandlerProps {
  onClick: (coord: Coord) => void;
}

function MapClickHandler({ onClick }: MapClickHandlerProps) {
  useMapEvents({
    click(e) {
      onClick({ lat: e.latlng.lat, lng: e.latlng.lng });
    },
  });
  return null;
}

interface MapProps {
  from: Coord | null;
  to: Coord | null;
  result: PathResult | null;
  onMapClick: (coord: Coord) => void;
  onFromMove: (coord: Coord) => void;
  onToMove: (coord: Coord) => void;
}

export function Map({ from, to, result, onMapClick, onFromMove, onToMove }: MapProps) {
  const [routeData, setRouteData] = useState<GeoJSON.FeatureCollection | null>(null);
  const [routeError, setRouteError] = useState(false);

  useEffect(() => {
    fetch('/routes.geojson')
      .then((r) => r.json())
      .then(setRouteData)
      .catch(() => setRouteError(true));
  }, []);

  const pathCoords = result ? decodePolyline(result.polyline) : null;

  return (
    <MapContainer
      center={OSLO_CENTER}
      zoom={DEFAULT_ZOOM}
      className="h-full w-full"
      // Prevent the panel's scroll from propagating to the map
      scrollWheelZoom
    >
      <TileLayer
        url="https://tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
        maxZoom={19}
      />

      {/* All registered outdoor routes */}
      {routeData && (
        <GeoJSON
          key="routes"
          data={routeData}
          style={(feature) => ({
            color: ROUTE_COLORS[feature?.properties?.route_type as string] ?? '#888',
            weight: 3,
            opacity: 0.75,
          })}
          onEachFeature={(feature, layer) => {
            const name = feature.properties?.name as string | null;
            if (name) layer.bindTooltip(name, { sticky: true });
          }}
        />
      )}

      {/* Route not yet loaded fallback */}
      {routeError && null /* silent — user can still use path finding */}

      {/* Shortest path overlay */}
      {pathCoords && (
        <Polyline
          positions={pathCoords}
          pathOptions={{ color: '#000000', weight: 5, opacity: 0.9 }}
        />
      )}

      {/* Start / end markers */}
      {from && (
        <Marker
          position={[from.lat, from.lng]}
          icon={START_PIN}
          draggable
          eventHandlers={{
            dragend(e) {
              const { lat, lng } = e.target.getLatLng();
              onFromMove({ lat, lng });
            },
          }}
        />
      )}
      {to && (
        <Marker
          position={[to.lat, to.lng]}
          icon={END_PIN}
          draggable
          eventHandlers={{
            dragend(e) {
              const { lat, lng } = e.target.getLatLng();
              onToMove({ lat, lng });
            },
          }}
        />
      )}

      <MapClickHandler onClick={onMapClick} />
    </MapContainer>
  );
}
