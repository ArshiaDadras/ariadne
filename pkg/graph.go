package pkg

import (
	"errors"
)

var (
	ErrNodeExists       = errors.New("node already exists")
	ErrEdgeExists       = errors.New("edge already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrEdgeNotFound     = errors.New("edge not found")
	ErrNodeNotReachable = errors.New("node not reachable")
)

type dijkstraData struct {
	MaxDuration float64
	Distances   map[*Node]float64
	Parents     map[*Node]*Node
	Visited     map[*Node]bool
	Queue       *Heap
}

type Node struct {
	ID       string                 `json:"id"`
	Position Point                  `json:"position"`
	InEdges  map[*Node]*Edge        `json:"-"`
	OutEdges map[*Node]*Edge        `json:"-"`
	Data     map[bool]*dijkstraData `json:"-"`
}

type Edge struct {
	ID     string  `json:"id"`
	Start  string  `json:"start"`
	End    string  `json:"end"`
	Speed  float64 `json:"speed"`
	Poly   []Point `json:"polygon"`
	Length float64 `json:"length"`
}

func NewEdge(id string, start, end *Node, speed float64, poly []Point) (edge *Edge) {
	edge = &Edge{
		ID:     id,
		Start:  start.ID,
		End:    end.ID,
		Speed:  speed,
		Poly:   poly,
		Length: 0,
	}

	for i := 1; i < len(poly); i++ {
		edge.Length += poly[i].Distance(poly[i-1])
	}
	return
}

func (e *Edge) LengthTo(point Point) (length float64) {
	intersect := point.ClosestPointOnEdge(e)
	for i := 1; i < len(e.Poly); i++ {
		if intersect.IsOnSegment(e.Poly[i-1], e.Poly[i]) {
			length += intersect.Distance(e.Poly[i-1])
			break
		}
		length += e.Poly[i-1].Distance(e.Poly[i])
	}
	return
}

func (e *Edge) LengthFrom(point Point) float64 {
	return e.Length - e.LengthTo(point)
}

type Graph struct {
	Nodes map[string]*Node `json:"nodes"`
	Edges map[string]*Edge `json:"edges"`
	Seg   *Segment2D       `json:"-"`
}

func NewGraph() (graph *Graph) {
	graph = &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string]*Edge),
	}
	return
}

func (g *Graph) AddNode(id string, position Point) (*Node, error) {
	if _, ok := g.Nodes[id]; ok {
		return nil, ErrNodeExists
	}

	g.Nodes[id] = &Node{
		ID:       id,
		Position: position,
		InEdges:  make(map[*Node]*Edge),
		OutEdges: make(map[*Node]*Edge),
		Data:     make(map[bool]*dijkstraData),
	}
	return g.Nodes[id], nil
}

func (g *Graph) AddEdge(id string, start, end *Node, speed float64, poly []Point) (*Edge, error) {
	if _, ok := g.Edges[id]; ok {
		return nil, ErrEdgeExists
	}
	if _, ok := g.Nodes[start.ID]; !ok {
		return nil, ErrNodeNotFound
	}
	if _, ok := g.Nodes[end.ID]; !ok {
		return nil, ErrNodeNotFound
	}

	edge := NewEdge(id, start, end, speed, poly)
	g.Edges[id] = edge

	start.OutEdges[end] = edge
	end.InEdges[start] = edge
	return edge, nil
}

func (g *Graph) GetNode(id string) (*Node, error) {
	node, ok := g.Nodes[id]
	if !ok {
		return nil, ErrNodeNotFound
	}
	return node, nil
}

func (g *Graph) GetEdge(id string) (*Edge, error) {
	edge, ok := g.Edges[id]
	if !ok {
		return nil, ErrEdgeNotFound
	}
	return edge, nil
}

func (g *Graph) getData(node *Node, maxDuration float64, reverse bool) *dijkstraData {
	if data, ok := node.Data[reverse]; !ok {
		g.dijkstra(node, maxDuration, reverse)
	} else if data.MaxDuration < maxDuration {
		g.dijkstra(node, maxDuration, reverse)
	}
	return node.Data[reverse]
}

func (g *Graph) GetDistance(start, end *Node, maxDuration float64, reverse bool) (float64, error) {
	data := g.getData(start, maxDuration, reverse)
	if distance, ok := data.Distances[end]; ok {
		return distance, nil
	}
	return -1, ErrNodeNotReachable
}

func (g *Graph) GetBestPath(start, end *Node, maxDuration float64, reverse bool) ([]*Edge, error) {
	data := g.getData(start, maxDuration, reverse)
	if _, ok := data.Distances[end]; !ok {
		return nil, ErrNodeNotReachable
	}

	path := make([]*Edge, 0)
	for current := end; current != start; current = data.Parents[current] {
		if reverse {
			path = append(path, current.OutEdges[data.Parents[current]])
		} else {
			path = append(path, current.InEdges[data.Parents[current]])
		}
	}
	return path, nil
}

type heapNode struct {
	node     *Node
	distance float64
}

func (g *Graph) dijkstra(start *Node, maxDuration float64, reverse bool) {
	if data, ok := start.Data[reverse]; !ok {
		start.Data[reverse] = &dijkstraData{
			MaxDuration: maxDuration,
			Distances:   make(map[*Node]float64),
			Parents:     make(map[*Node]*Node),
			Visited:     make(map[*Node]bool),
			Queue: NewHeap(func(i, j interface{}) bool {
				if i.(heapNode).distance == j.(heapNode).distance {
					return i.(heapNode).node.ID < j.(heapNode).node.ID
				} else {
					return i.(heapNode).distance < j.(heapNode).distance
				}
			}),
		}

		start.Data[reverse].Distances[start] = 0
		start.Data[reverse].Queue.Push(heapNode{node: start, distance: 0})
	} else if data.MaxDuration < maxDuration {
		data.MaxDuration = maxDuration
	} else {
		return
	}

	priorityQueue := start.Data[reverse].Queue
	visited := start.Data[reverse].Visited
	dist := start.Data[reverse].Distances
	par := start.Data[reverse].Parents

	for priorityQueue.Length() > 0 {
		current := priorityQueue.Pop().(heapNode)
		if maxDuration > 0 && current.distance > maxDuration {
			break
		}
		if visited[current.node] {
			continue
		}
		visited[current.node] = true

		g.updateDistances(current, priorityQueue, visited, dist, par, reverse)
	}
}

func (g *Graph) updateDistances(current heapNode, priorityQueue *Heap, visited map[*Node]bool, dist map[*Node]float64, par map[*Node]*Node, reverse bool) {
	edges := current.node.OutEdges
	if reverse {
		edges = current.node.InEdges
	}

	for neighbour, edge := range edges {
		if !visited[neighbour] {
			distance := current.distance + edge.Length
			if current_distance, ok := dist[neighbour]; !ok || distance < current_distance {
				dist[neighbour], par[neighbour] = distance, current.node
				priorityQueue.Push(heapNode{node: neighbour, distance: distance})
			}
		}
	}
}
