package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type pathRequest struct {
	From Coord `json:"from"` // { "lat": 59.92, "lng": 10.69 }
	To   Coord `json:"to"`
}

type pathResponse struct {
	Polyline   string  `json:"polyline"`
	DistanceKm float64 `json:"distance_km"`
	Waypoints  int     `json:"waypoints"`
}

func runServer(addr string, g *Graph) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shortest-path", func(w http.ResponseWriter, r *http.Request) {
		var req pathRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		startID := g.NearestNode(req.From)
		endID := g.NearestNode(req.To)

		path, metres := g.ShortestPath(startID, endID)
		if path == nil {
			http.Error(w, "no path found", http.StatusNotFound)
			return
		}

		coords := g.PathCoords(path)
		resp := pathResponse{
			Polyline:   EncodePolyline(coords),
			DistanceKm: metres / 1000,
			Waypoints:  len(coords),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // allow the static frontend to call this
		json.NewEncoder(w).Encode(resp)
	})

	// Preflight for browsers that send OPTIONS before POST
	mux.HandleFunc("OPTIONS /shortest-path", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	})

	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func (c Coord) String() string {
	return fmt.Sprintf("%.5f,%.5f", c.Lat, c.Lng)
}
