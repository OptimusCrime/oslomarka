import { useState } from 'react';
import type { Coord, PathResult } from '../types';
import { fetchShortestPath } from '../api/client';

interface UseShortestPathReturn {
  result: PathResult | null;
  isLoading: boolean;
  error: string | null;
  find: (from: Coord, to: Coord) => Promise<void>;
  clear: () => void;
}

export function useShortestPath(): UseShortestPathReturn {
  const [result, setResult] = useState<PathResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function find(from: Coord, to: Coord) {
    setIsLoading(true);
    setError(null);
    setResult(null);
    try {
      const data = await fetchShortestPath(from, to);
      setResult(data);
    } catch {
      setError('Could not find a path between these points.');
    } finally {
      setIsLoading(false);
    }
  }

  function clear() {
    setResult(null);
    setError(null);
  }

  return { result, isLoading, error, find, clear };
}
