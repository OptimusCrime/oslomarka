package graph

import (
	"container/heap"
	"math"
	"slices"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
)

type Graph struct {
	adjacency map[geo.Coord][]edge
	// positions maps each snapped node key back to an original (un-snapped)
	// coordinate, so paths are returned at full input precision.
	positions map[geo.Coord]geo.Coord
	numEdges  int
}

type edge struct {
	to   geo.Coord
	dist float64
}

func newGraph() *Graph {
	return &Graph{
		adjacency: make(map[geo.Coord][]edge),
		positions: make(map[geo.Coord]geo.Coord),
	}
}

func (g *Graph) addEdge(a, b geo.Coord) {
	idA := geo.Snap(a)
	idB := geo.Snap(b)
	dist := geo.Haversine(a, b)

	g.positions[idA] = a
	g.positions[idB] = b
	g.adjacency[idA] = append(g.adjacency[idA], edge{to: idB, dist: dist})
	g.adjacency[idB] = append(g.adjacency[idB], edge{to: idA, dist: dist})
	g.numEdges++
}

// BuildGraph constructs a graph from route coordinate slices.
func BuildGraph(routeCoords [][]geo.Coord) *Graph {
	g := newGraph()
	for _, coords := range routeCoords {
		for i := range len(coords) - 1 {
			g.addEdge(coords[i], coords[i+1])
		}
	}
	return g
}

func (g *Graph) NodeCount() int { return len(g.positions) }
func (g *Graph) EdgeCount() int { return g.numEdges }

func (g *Graph) NearestNode(target geo.Coord) geo.Coord {
	var bestID geo.Coord
	bestDist := math.MaxFloat64
	for id, c := range g.positions {
		if d := geo.Haversine(target, c); d < bestDist {
			bestDist = d
			bestID = id
		}
	}
	return bestID
}

func (g *Graph) PathCoords(path []geo.Coord) []geo.Coord {
	coords := make([]geo.Coord, len(path))
	for i, id := range path {
		coords[i] = g.positions[id]
	}
	return coords
}

func (g *Graph) ShortestPath(startID, endID geo.Coord) ([]geo.Coord, float64) {
	dist := map[geo.Coord]float64{startID: 0}
	prev := map[geo.Coord]geo.Coord{}

	pq := priorityQueue{{id: startID, cost: 0}}
	heap.Init(&pq)

	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(*pqItem)

		if cur.id == endID {
			break
		}
		if cur.cost > dist[cur.id] {
			continue
		}

		for _, e := range g.adjacency[cur.id] {
			newCost := dist[cur.id] + e.dist
			if d, seen := dist[e.to]; !seen || newCost < d {
				dist[e.to] = newCost
				prev[e.to] = cur.id
				heap.Push(&pq, &pqItem{id: e.to, cost: newCost})
			}
		}
	}

	total, reached := dist[endID]
	if !reached {
		return nil, 0
	}

	// Walk prev pointers from endID back to startID (the only reached node
	// without a predecessor), then reverse into start-to-end order.
	var path []geo.Coord
	for n := endID; ; n = prev[n] {
		path = append(path, n)
		if n == startID {
			break
		}
	}
	slices.Reverse(path)
	return path, total
}

type pqItem struct {
	id   geo.Coord
	cost float64
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x any) { *pq = append(*pq, x.(*pqItem)) }
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
