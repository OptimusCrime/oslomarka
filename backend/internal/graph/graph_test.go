package graph

import (
	"math"
	"testing"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
)

// Layout: a straight line A-B-D running east, plus a long detour A-C-D that
// loops far north. The shortest route from A to D must be A-B-D.
func TestShortestPathPicksShorterRoute(t *testing.T) {
	a := geo.Coord{Lat: 59.90, Lng: 10.70}
	b := geo.Coord{Lat: 59.90, Lng: 10.71}
	d := geo.Coord{Lat: 59.90, Lng: 10.72}
	c := geo.Coord{Lat: 59.95, Lng: 10.71}

	g := BuildGraph([][]geo.Coord{
		{a, b, d},
		{a, c, d},
	})

	path, metres := g.ShortestPath(g.NearestNode(a), g.NearestNode(d))
	if path == nil {
		t.Fatal("expected a path, got none")
	}

	coords := g.PathCoords(path)
	if len(coords) != 3 {
		t.Fatalf("expected 3 waypoints (A, B, D), got %d: %v", len(coords), coords)
	}

	want := geo.Haversine(a, b) + geo.Haversine(b, d)
	if math.Abs(metres-want) > 1e-6 {
		t.Errorf("distance = %v m, want %v m", metres, want)
	}

	// "exact output" guarantee: endpoints are the original coords, not snapped.
	if coords[0] != a || coords[len(coords)-1] != d {
		t.Errorf("endpoints not preserved: got %v ... %v", coords[0], coords[len(coords)-1])
	}
}

// Two disconnected route fragments: no path should exist between them.
func TestShortestPathDisconnected(t *testing.T) {
	left := geo.Coord{Lat: 59.90, Lng: 10.70}
	leftB := geo.Coord{Lat: 59.90, Lng: 10.71}
	right := geo.Coord{Lat: 60.50, Lng: 11.50}
	rightB := geo.Coord{Lat: 60.50, Lng: 11.51}

	g := BuildGraph([][]geo.Coord{
		{left, leftB},
		{right, rightB},
	})

	if path, _ := g.ShortestPath(g.NearestNode(left), g.NearestNode(right)); path != nil {
		t.Errorf("expected no path between disconnected components, got %v", path)
	}
}
