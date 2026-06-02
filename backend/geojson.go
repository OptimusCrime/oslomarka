package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Route holds a parsed route segment and its waypoints in (lat, lng) order.
type Route struct {
	Name        string
	RouteType   string // "Fotrute", "Sykkelrute", "AnnenRute"
	SurfaceCode string // "ST" (trail), "BV" (gravel road), "GS" (cycle path), …
	Marked      string // "JA" or "NEI"
	Coords      []Coord
}

// The structs below mirror the GeoJSON structure so encoding/json can decode
// the file directly. They are unexported because nothing outside this file
// needs them — callers work with []Route.

type featureCollection struct {
	Features []feature `json:"features"`
}

type feature struct {
	Geometry   geometry   `json:"geometry"`
	Properties properties `json:"properties"`
}

type geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"` // each entry is [lng, lat]
}

type properties struct {
	RouteType   string `json:"route_type"`
	Name        string `json:"name"`
	SurfaceCode string `json:"surface_code"`
	Marked      string `json:"marked"`
}

// WritePathGeoJSON encodes a resolved path as a GeoJSON FeatureCollection
// containing a single LineString. The result can be dragged into geojson.io,
// loaded with Leaflet's L.geoJSON(), or added to Google Maps via map.data.
func WritePathGeoJSON(w io.Writer, coords []Coord, metres float64) error {
	// GeoJSON coordinates are [lng, lat] — the opposite of our internal order.
	geoCoords := make([][]float64, len(coords))
	for i, c := range coords {
		geoCoords[i] = []float64{c.Lng, c.Lat}
	}
	fc := map[string]any{
		"type": "FeatureCollection",
		"features": []any{
			map[string]any{
				"type": "Feature",
				"geometry": map[string]any{
					"type":        "LineString",
					"coordinates": geoCoords,
				},
				"properties": map[string]any{
					"distance_km": metres / 1000,
					"waypoints":   len(coords),
					"polyline":    EncodePolyline(coords),
				},
			},
		},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(fc)
}

// LoadGeoJSON reads a GeoJSON file and returns all LineString features as Routes.
func LoadGeoJSON(path string) ([]Route, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var fc featureCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("parse GeoJSON: %w", err)
	}

	routes := make([]Route, 0, len(fc.Features))
	for _, f := range fc.Features {
		if f.Geometry.Type != "LineString" {
			continue
		}
		coords := make([]Coord, len(f.Geometry.Coordinates))
		for i, c := range f.Geometry.Coordinates {
			// GeoJSON stores coordinates as [lng, lat]; we flip to (lat, lng).
			coords[i] = Coord{Lat: c[1], Lng: c[0]}
		}
		routes = append(routes, Route{
			Name:        f.Properties.Name,
			RouteType:   f.Properties.RouteType,
			SurfaceCode: f.Properties.SurfaceCode,
			Marked:      f.Properties.Marked,
			Coords:      coords,
		})
	}
	return routes, nil
}
