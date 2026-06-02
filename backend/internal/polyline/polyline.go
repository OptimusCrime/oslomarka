package polyline

import (
	"math"
	"strings"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
)

// Encode encodes a slice of coordinates using Google's Encoded Polyline Algorithm.
func Encode(coords []geo.Coord) string {
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
	v := delta << 1
	if v < 0 {
		v = ^v
	}
	for v >= 0x20 {
		buf.WriteByte(byte((0x20|(v&0x1f)) + 63))
		v >>= 5
	}
	buf.WriteByte(byte(v + 63))
}
