/**
 * Map colors shared between the Leaflet layers (Map.tsx) and the legend /
 * coordinate readout (Panel.tsx). Kept in one place so the map and its legend
 * can never drift apart.
 *
 * These live in TS rather than the Tailwind theme because they are consumed by
 * Leaflet (path options, marker icons) and inline styles, not utility classes.
 */
export const MAP_COLORS = {
  start: '#16a34a',
  end: '#dc2626',
  fotrute: '#e63946',
  shortestPath: '#000000',
  unknownRoute: '#888888',
} as const;
