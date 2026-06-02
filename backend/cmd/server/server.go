package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
	"github.com/OptimusCrime/oslomarka/backend/internal/geojson"
	"github.com/OptimusCrime/oslomarka/backend/internal/graph"
	"github.com/OptimusCrime/oslomarka/backend/internal/path"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	geojsonPath := os.Getenv("GEOJSON_PATH")
	if geojsonPath == "" {
		geojsonPath = "../_data/routes.geojson"
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("loading routes", "path", geojsonPath)
	routes, err := geojson.LoadGeoJSON(geojsonPath)
	if err != nil {
		slog.Error("failed to load GeoJSON", "err", err)
		os.Exit(1)
	}

	coordSlices := make([][]geo.Coord, len(routes))
	for i, r := range routes {
		coordSlices[i] = r.Coords
	}

	slog.Info("building graph", "routes", len(routes))
	g := graph.BuildGraph(coordSlices)
	slog.Info("graph ready", "nodes", g.NodeCount(), "edges", g.EdgeCount())

	pathHandler := path.NewHandler(path.NewService(g))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	r.Post("/shortest-path", pathHandler.ShortestPath)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}
