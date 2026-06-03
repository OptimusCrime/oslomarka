import ky from 'ky';
import type { Coord, PathResult } from '../types';

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '';

export async function fetchShortestPath(from: Coord, to: Coord): Promise<PathResult> {
  return ky.post(`${API_BASE}/shortest-path`, { retry: 0, json: { from, to } }).json<PathResult>();
}
