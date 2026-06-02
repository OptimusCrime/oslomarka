package path

import (
	"errors"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
	"github.com/OptimusCrime/oslomarka/backend/internal/graph"
	"github.com/OptimusCrime/oslomarka/backend/internal/polyline"
)

type Result struct {
	Polyline   string  `json:"polyline"`
	DistanceKm float64 `json:"distance_km"`
	Waypoints  int     `json:"waypoints"`
}

type Service struct {
	g *graph.Graph
}

func NewService(g *graph.Graph) *Service {
	return &Service{g: g}
}

func (s *Service) ShortestPath(from, to geo.Coord) (*Result, error) {
	startID := s.g.NearestNode(from)
	endID := s.g.NearestNode(to)

	p, metres := s.g.ShortestPath(startID, endID)
	if p == nil {
		return nil, errors.New("no path found")
	}

	coords := s.g.PathCoords(p)
	return &Result{
		Polyline:   polyline.Encode(coords),
		DistanceKm: metres / 1000,
		Waypoints:  len(coords),
	}, nil
}
