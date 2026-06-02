package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	geojsonPath := flag.String("geojson", "../routes.geojson", "path to the GeoJSON file")
	fromFlag := flag.String("from", "", "start coordinate as lat,lng  e.g. 59.9264,10.6986")
	toFlag := flag.String("to", "", "end coordinate as lat,lng    e.g. 59.9750,10.8100")
	outputPath := flag.String("output", "", "write the resolved path as GeoJSON to this file")
	serveAddr := flag.String("serve", "", "run as HTTP server on this address  e.g. :8080")
	flag.Parse()

	fmt.Printf("Loading %s …\n", *geojsonPath)
	routes, err := LoadGeoJSON(*geojsonPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d routes\n", len(routes))

	fmt.Println("Building graph …")
	g := BuildGraph(routes)
	fmt.Printf("Graph: %d nodes, %d edges\n\n", len(g.positions), g.numEdges)

	if *serveAddr != "" {
		runServer(*serveAddr, g)
		return
	}

	if *fromFlag == "" || *toFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	from, err := parseCoord(*fromFlag)
	if err != nil {
		log.Fatalf("-from: %v", err)
	}
	to, err := parseCoord(*toFlag)
	if err != nil {
		log.Fatalf("-to: %v", err)
	}

	startID := g.NearestNode(from)
	endID := g.NearestNode(to)
	fmt.Printf("Start node: %s\n", startID)
	fmt.Printf("End node:   %s\n\n", endID)

	fmt.Println("Running Dijkstra …")
	path, metres := g.ShortestPath(startID, endID)
	if path == nil {
		fmt.Println("No path found between those two points.")
		return
	}

	coords := g.PathCoords(path)
	polyline := EncodePolyline(coords)

	fmt.Printf("\nShortest path: %.2f km over %d waypoints\n", metres/1000, len(path))
	fmt.Printf("Encoded polyline (%d chars): %s…\n", len(polyline), polyline[:min(60, len(polyline))])
	fmt.Println("\nFirst 5 waypoints:")
	for i, id := range path {
		if i == 5 {
			fmt.Printf("  … and %d more\n", len(path)-5)
			break
		}
		c := g.positions[id]
		fmt.Printf("  [%d] lat=%.5f  lng=%.5f\n", i+1, c.Lat, c.Lng)
	}

	if *outputPath != "" {
		f, err := os.Create(*outputPath)
		if err != nil {
			log.Fatalf("create output file: %v", err)
		}
		defer f.Close()
		if err := WritePathGeoJSON(f, coords, metres); err != nil {
			log.Fatalf("write GeoJSON: %v", err)
		}
		fmt.Printf("\nPath written to %s\n", *outputPath)
	}
}

func parseCoord(s string) (Coord, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return Coord{}, fmt.Errorf("expected lat,lng but got %q", s)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Coord{}, fmt.Errorf("invalid latitude: %w", err)
	}
	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Coord{}, fmt.Errorf("invalid longitude: %w", err)
	}
	return Coord{Lat: lat, Lng: lng}, nil
}
