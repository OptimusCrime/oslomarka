package graph

import (
	"container/heap"
	"math"

	"github.com/OptimusCrime/oslomarka/backend/internal/geo"
)

// Graph is an undirected weighted graph. Nodes are identified by NodeID strings
// (snapped lat/lng), and edge weights are distances in metres.
type Graph struct {
	adjacency map[string][]edge
	positions map[string]geo.Coord
	numEdges  int
}

type edge struct {
	to   string
	dist float64
}

func newGraph() *Graph {
	return &Graph{
		adjacency: make(map[string][]edge),
		positions: make(map[string]geo.Coord),
	}
}

func (g *Graph) addEdge(a, b geo.Coord) {
	idA := geo.NodeID(a.Lat, a.Lng)
	idB := geo.NodeID(b.Lat, b.Lng)
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

// NearestNode returns the NodeID of the graph node closest to target.
func (g *Graph) NearestNode(target geo.Coord) string {
	bestID := ""
	bestDist := math.MaxFloat64
	for id, c := range g.positions {
		if d := geo.Haversine(target, c); d < bestDist {
			bestDist = d
			bestID = id
		}
	}
	return bestID
}

// PathCoords resolves a slice of NodeIDs into their actual Coord values.
func (g *Graph) PathCoords(path []string) []geo.Coord {
	coords := make([]geo.Coord, len(path))
	for i, id := range path {
		coords[i] = g.positions[id]
	}
	return coords
}

// ShortestPath runs Dijkstra's algorithm from startID to endID.
// Returns the path as ordered NodeIDs and total distance in metres.
func (g *Graph) ShortestPath(startID, endID string) ([]string, float64) {
	dist := map[string]float64{startID: 0}
	prev := map[string]string{}

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

	var path []string
	for n := endID; n != ""; n = prev[n] {
		path = append(path, n)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, total
}

type pqItem struct {
	id   string
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
