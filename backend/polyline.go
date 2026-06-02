package main

import (
	"math"
	"strings"
)

// EncodePolyline encodes a slice of coordinates using Google's Encoded Polyline
// Algorithm. The resulting string can be passed directly to the Google Maps
// JavaScript API (google.maps.geometry.encoding.decodePath) or the Static Maps
// API (&path=enc:...) without any further transformation.
//
// The algorithm works in three steps per coordinate:
//  1. Delta-encode: store the difference from the previous point, not the
//     absolute value. This keeps numbers small for nearby points.
//  2. Zigzag-encode: left-shift by 1, then invert all bits if negative.
//     This maps negative numbers to positive odd integers so a single
//     unsigned encoding handles both signs.
//  3. Chunk into 5-bit groups, set a continuation bit (0x20) on every group
//     except the last, then shift into printable ASCII by adding 63.
func EncodePolyline(coords []Coord) string {
	var buf strings.Builder
	prevLat, prevLng := 0, 0

	for _, c := range coords {
		lat := int(math.Round(c.Lat * 1e5))
		lng := int(math.Round(c.Lng * 1e5))

		encodeChunks(&buf, lat-prevLat)
		encodeChunks(&buf, lng-prevLng)

		prevLat, prevLng = lat, lng
	}
	return buf.String()
}

func encodeChunks(buf *strings.Builder, delta int) {
	// Zigzag: left-shift then invert if negative.
	v := delta << 1
	if v < 0 {
		v = ^v
	}
	// Emit 5-bit chunks, least-significant first. Set the continuation bit
	// (0x20) on every chunk that still has more data following it.
	for v >= 0x20 {
		buf.WriteByte(byte((0x20|(v&0x1f)) + 63))
		v >>= 5
	}
	buf.WriteByte(byte(v + 63))
}
