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
    lat += decodeChunk(encoded, index);
    index += chunkLength(encoded, index);

    lng += decodeChunk(encoded, index);
    index += chunkLength(encoded, index);

    result.push([lat / 1e5, lng / 1e5]);
  }

  return result;
}

function decodeChunk(encoded: string, start: number): number {
  let shift = 0;
  let value = 0;
  let i = start;
  let b: number;
  do {
    b = encoded.charCodeAt(i++) - 63;
    value |= (b & 0x1f) << shift;
    shift += 5;
  } while (b >= 0x20);
  return value & 1 ? ~(value >> 1) : value >> 1;
}

function chunkLength(encoded: string, start: number): number {
  let i = start;
  while (encoded.charCodeAt(i++) - 63 >= 0x20) {
    /* scan until continuation bit is unset */
  }
  return i - start;
}
