package path

import (
	"encoding/json"
	"net/http"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
	"github.com/OptimusCrime/oslomarka/backend/internal/render"
	"github.com/OptimusCrime/oslomarka/backend/internal/resterr"
)

type pathRequest struct {
	From geo.Coord `json:"from"`
	To   geo.Coord `json:"to"`
}

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) ShortestPath(w http.ResponseWriter, r *http.Request) {
	var req pathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, resterr.New("invalid JSON body", http.StatusBadRequest))
		return
	}

	result, err := h.s.ShortestPath(req.From, req.To)
	if err != nil {
		render.JSON(w, r, resterr.FromErr(err, http.StatusNotFound))
		return
	}

	render.JSON(w, r, result)
}
