import { useState } from 'react';
import { Map } from './components/Map';
import { Panel } from './components/Panel';
import { useShortestPath } from './hooks/useShortestPath';
import type { Coord } from './types';

export default function App() {
  const [from, setFrom] = useState<Coord | null>(null);
  const [to, setTo]     = useState<Coord | null>(null);
  const { result, isLoading, error, find, clear } = useShortestPath();

  function handleMapClick(coord: Coord) {
    // If both points are already set, start over
    if (from && to) {
      clear();
      setFrom(coord);
      setTo(null);
      return;
    }
    if (!from) { setFrom(coord); return; }
    setTo(coord);
  }

  function handleFindPath() {
    if (from && to) find(from, to);
  }

  function handleClear() {
    setFrom(null);
    setTo(null);
    clear();
  }

  function handleFromMove(coord: Coord) {
    setFrom(coord);
    clear();
  }

  function handleToMove(coord: Coord) {
    setTo(coord);
    clear();
  }

  return (
    <div className="flex h-screen flex-col overflow-hidden">
      <header className="z-[1001] flex h-14 shrink-0 items-center border-b border-border bg-forest px-5">
        <div className="flex items-center gap-2.5">
          <span className="text-xl">🌲</span>
          <h1 className="font-display text-lg font-bold tracking-tight text-white">
            Oslomarkakart
          </h1>
        </div>
      </header>

      {/* Map fills remaining height; Panel floats over it */}
      <div className="relative flex-1">
        <Map
          from={from}
          to={to}
          result={result}
          onMapClick={handleMapClick}
          onFromMove={handleFromMove}
          onToMove={handleToMove}
        />
        <Panel
          from={from}
          to={to}
          result={result}
          isLoading={isLoading}
          error={error}
          onFindPath={handleFindPath}
          onClear={handleClear}
        />
      </div>
    </div>
  );
}
