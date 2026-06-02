import type { Coord, PathResult } from '../types';
import { Button } from './ui/button';

interface PanelProps {
  from: Coord | null;
  to: Coord | null;
  result: PathResult | null;
  isLoading: boolean;
  error: string | null;
  onFindPath: () => void;
  onClear: () => void;
}

function CoordRow({ label, coord, color }: { label: string; coord: Coord | null; color: string }) {
  return (
    <div className="flex items-start gap-2.5">
      <span
        className="mt-0.5 h-3 w-3 shrink-0 rounded-full border-2 border-white shadow-sm"
        style={{ background: color }}
      />
      <div className="min-w-0">
        <p className="text-xs font-medium text-ink-secondary">{label}</p>
        {coord ? (
          <p className="font-mono text-xs text-ink">
            {coord.lat.toFixed(5)}, {coord.lng.toFixed(5)}
          </p>
        ) : (
          <p className="text-xs text-ink-muted italic">Click on the map…</p>
        )}
      </div>
    </div>
  );
}

export function Panel({ from, to, result, isLoading, error, onFindPath, onClear }: PanelProps) {
  const canFind = from !== null && to !== null;

  return (
    <div
      className="absolute left-4 top-4 z-[1000] w-64 rounded-xl bg-surface shadow-lg border border-border
                 flex flex-col gap-4 p-4 pointer-events-auto"
    >
      {/* Instruction */}
      <p className="text-xs text-ink-secondary leading-snug">
        {!from
          ? 'Click on the map to set a start point.'
          : !to
            ? 'Now click to set an end point.'
            : 'Both points set. Find the shortest path along the routes.'}
      </p>

      {/* Coordinate display */}
      <div className="flex flex-col gap-3">
        <CoordRow label="Start" coord={from} color="#16a34a" />
        <CoordRow label="End"   coord={to}   color="#dc2626" />
      </div>

      {/* Actions */}
      <div className="flex flex-col gap-2">
        <Button onClick={onFindPath} disabled={!canFind} isLoading={isLoading} className="w-full">
          Find shortest path
        </Button>
        {(from || to || result) && (
          <Button variant="secondary" onClick={onClear} className="w-full">
            Clear
          </Button>
        )}
      </div>

      {/* Error */}
      {error && <p className="text-xs text-red-600">{error}</p>}

      {/* Result */}
      {result && !error && (
        <div className="rounded-lg bg-surface-hover border border-border px-3 py-2.5 flex flex-col gap-0.5">
          <p className="text-sm font-semibold text-ink">
            {result.distance_km.toFixed(2)} km
          </p>
          <p className="text-xs text-ink-secondary">{result.waypoints} waypoints along registered routes</p>
        </div>
      )}

      {/* Legend */}
      <div className="border-t border-border pt-3 flex flex-col gap-1.5">
        <p className="text-[11px] font-medium text-ink-secondary uppercase tracking-wide">Routes</p>
        <LegendItem color="#e63946" label="Fotrute (trail)" />
        <LegendItem color="#000000" label="Shortest path" thick />
      </div>
    </div>
  );
}

function LegendItem({ color, label, thick }: { color: string; label: string; thick?: boolean }) {
  return (
    <div className="flex items-center gap-2">
      <div
        className="shrink-0 rounded-full"
        style={{ width: 20, height: thick ? 4 : 3, background: color }}
      />
      <span className="text-xs text-ink-secondary">{label}</span>
    </div>
  );
}
