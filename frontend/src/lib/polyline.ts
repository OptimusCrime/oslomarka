/**
 * Decodes a Google Encoded Polyline string into an array of [lat, lng] pairs.
 * This is the exact reverse of the encoder in the Go backend (shortestpath/polyline.go).
 *
 * react-leaflet's <Polyline> accepts LatLngTuple[] which is [lat, lng][],
 * so the output can be passed directly as the `positions` prop.
 */
export function decodePolyline(encoded: string): [number, number][] {
  const result: [number, number][] = [];
  let index = 0;
  let lat = 0;
  let lng = 0;

  while (index < encoded.length) {
    let dLat: number;
    let dLng: number;
    [dLat, index] = decodeValue(encoded, index);
    [dLng, index] = decodeValue(encoded, index);

    lat += dLat;
    lng += dLng;
    result.push([lat / 1e5, lng / 1e5]);
  }

  return result;
}

/**
 * Decodes a single signed value starting at `start` and returns it together
 * with the index of the next chunk, so each chunk is scanned exactly once.
 */
function decodeValue(encoded: string, start: number): [value: number, next: number] {
  let shift = 0;
  let value = 0;
  let index = start;
  let byte: number;
  do {
    byte = encoded.charCodeAt(index++) - 63;
    value |= (byte & 0x1f) << shift;
    shift += 5;
  } while (byte >= 0x20);

  const delta = value & 1 ? ~(value >> 1) : value >> 1;
  return [delta, index];
}
