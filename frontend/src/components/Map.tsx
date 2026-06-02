import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON, Polyline, Marker, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import type { Coord, PathResult } from '../types';
import { decodePolyline } from '../lib/polyline';
import { MAP_COLORS } from '../lib/colors';

// Oslo city centre
const OSLO_CENTER: [number, number] = [59.92, 10.75];
const DEFAULT_ZOOM = 11;

const ROUTE_COLORS: Record<string, string> = {
  Fotrute: MAP_COLORS.fotrute,
};

// Color each route by its type, falling back to a neutral grey for unknown types.
function routeStyle(feature?: GeoJSON.Feature): L.PathOptions {
  const type = feature?.properties?.route_type as string;
  return { color: ROUTE_COLORS[type] ?? MAP_COLORS.unknownRoute, weight: 3, opacity: 0.75 };
}

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

const START_PIN = makePin(MAP_COLORS.start);
const END_PIN = makePin(MAP_COLORS.end);

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

  useEffect(() => {
    fetch('/routes.geojson')
      .then((r) => r.json())
      .then(setRouteData)
      // Non-fatal: the map and path finding still work without the route overlay.
      .catch((err) => console.error('Failed to load route overlay', err));
  }, []);

  const pathCoords = result ? decodePolyline(result.polyline) : null;

  return (
    <MapContainer
      center={OSLO_CENTER}
      zoom={DEFAULT_ZOOM}
      className="h-full w-full"
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
          style={routeStyle}
          onEachFeature={(feature, layer) => {
            const name = feature.properties?.name as string | null;
            if (name) layer.bindTooltip(name, { sticky: true });
          }}
        />
      )}

      {/* Shortest path overlay */}
      {pathCoords && (
        <Polyline
          positions={pathCoords}
          pathOptions={{ color: MAP_COLORS.shortestPath, weight: 5, opacity: 0.9 }}
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
