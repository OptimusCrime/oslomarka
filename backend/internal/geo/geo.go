package geo

import (
	"fmt"
	"math"
)

// Coord is a geographic point in WGS84 (latitude/longitude in degrees).
type Coord struct {
	Lat, Lng float64
}

// NodeID returns a stable string key for a coordinate, snapped to 5 decimal
// places (~1 metre precision).
func NodeID(lat, lng float64) string {
	return fmt.Sprintf("%.5f,%.5f", snap(lat), snap(lng))
}

func snap(x float64) float64 {
	return math.Round(x*1e5) / 1e5
}

// Haversine returns the great-circle distance between two points in metres.
func Haversine(a, b Coord) float64 {
	const earthRadius = 6_371_000.0
	lat1 := toRad(a.Lat)
	lat2 := toRad(b.Lat)
	dlat := toRad(b.Lat - a.Lat)
	dlng := toRad(b.Lng - a.Lng)
	h := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlng/2)*math.Sin(dlng/2)
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

func toRad(deg float64) float64 { return deg * math.Pi / 180 }
