import ky from 'ky';
import type { Coord, PathResult } from '../types';

const api = ky.create({ retry: 0 });

export async function fetchShortestPath(from: Coord, to: Coord): Promise<PathResult> {
  return api.post('/shortest-path', { json: { from, to } }).json<PathResult>();
}
