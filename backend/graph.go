package main

import (
	"container/heap"
	"math"
)

// Graph is an undirected weighted graph. Nodes are identified by NodeID strings
// (snapped lat/lng), and edge weights are distances in metres.
type Graph struct {
	adjacency map[string][]edge
	positions map[string]Coord // NodeID → coordinate for reverse-lookup
	numEdges  int
}

type edge struct {
	to   string
	dist float64 // metres
}

func NewGraph() *Graph {
	return &Graph{
		adjacency: make(map[string][]edge),
		positions:  make(map[string]Coord),
	}
}

func (g *Graph) addEdge(a, b Coord) {
	idA := NodeID(a.Lat, a.Lng)
	idB := NodeID(b.Lat, b.Lng)
	dist := Haversine(a, b)

	g.positions[idA] = a
	g.positions[idB] = b
	g.adjacency[idA] = append(g.adjacency[idA], edge{to: idB, dist: dist})
	g.adjacency[idB] = append(g.adjacency[idB], edge{to: idA, dist: dist})
	g.numEdges++
}

// BuildGraph constructs a graph from route segments. Consecutive waypoints
// on the same route become connected edges; routes that share an endpoint
// are automatically joined at that node.
func BuildGraph(routes []Route) *Graph {
	g := NewGraph()
	for _, route := range routes {
		for i := range len(route.Coords) - 1 {
			g.addEdge(route.Coords[i], route.Coords[i+1])
		}
	}
	return g
}

// NearestNode returns the NodeID of the graph node closest to target.
// For our dataset (~65k nodes) a linear scan completes in milliseconds.
func (g *Graph) NearestNode(target Coord) string {
	bestID := ""
	bestDist := math.MaxFloat64
	for id, c := range g.positions {
		if d := Haversine(target, c); d < bestDist {
			bestDist = d
			bestID = id
		}
	}
	return bestID
}

// PathCoords resolves a slice of NodeIDs (as returned by ShortestPath) into
// their actual Coord values, ready to be written to a map or GeoJSON file.
func (g *Graph) PathCoords(path []string) []Coord {
	coords := make([]Coord, len(path))
	for i, id := range path {
		coords[i] = g.positions[id]
	}
	return coords
}

// ShortestPath runs Dijkstra's algorithm from startID to endID.
// It returns the path as an ordered slice of NodeIDs and the total distance
// in metres. If no path exists both return values are nil/0.
func (g *Graph) ShortestPath(startID, endID string) (path []string, totalMetres float64) {
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
			continue // stale entry; a cheaper path was already processed
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

	// Reconstruct path by walking prev pointers from end to start, then reverse.
	for n := endID; n != ""; n = prev[n] {
		path = append(path, n)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, total
}

// --- Min-heap used by Dijkstra ---

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
