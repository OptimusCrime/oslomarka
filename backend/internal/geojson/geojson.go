package geojson

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
)

type Route struct {
	Name        string
	RouteType   string
	SurfaceCode string
	Marked      string
	Coords      []geo.Coord
}

type featureCollection struct {
	Features []feature `json:"features"`
}

type feature struct {
	Geometry   geometry   `json:"geometry"`
	Properties properties `json:"properties"`
}

type geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type properties struct {
	RouteType   string `json:"route_type"`
	Name        string `json:"name"`
	SurfaceCode string `json:"surface_code"`
	Marked      string `json:"marked"`
}

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
		coords := make([]geo.Coord, len(f.Geometry.Coordinates))
		for i, c := range f.Geometry.Coordinates {
			// GeoJSON stores coordinates as [lng, lat]; we flip to (lat, lng).
			coords[i] = geo.Coord{Lat: c[1], Lng: c[0]}
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
