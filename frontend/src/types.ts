export interface Coord {
  lat: number;
  lng: number;
}

export interface PathResult {
  polyline: string;
  distance_km: number;
  waypoints: number;
}
