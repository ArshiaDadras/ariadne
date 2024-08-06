package pkg

import "errors"

var (
	ErrNodeExists       = errors.New("node already exists")
	ErrEdgeExists       = errors.New("edge already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrEdgeNotFound     = errors.New("edge not found")
	ErrNodeNotReachable = errors.New("node not reachable")
)

const (
	AverageSpeed = 50.0
)

type Node struct {
	ID       string   `json:"id"`
	Position Point    `json:"position"`
	InEdges  []string `json:"in_edges"`
	OutEdges []string `json:"out_edges"`
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

	for i := 0; i < len(poly)-1; i++ {
		edge.Length += poly[i].Distance(poly[i+1])
	}
	return
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

	start.OutEdges = append(start.OutEdges, edge.ID)
	end.InEdges = append(end.InEdges, edge.ID)
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

func (g *Graph) Preprocess() {
	nodes := make([]*Node, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}
	g.Seg = NewSegment2D(nodes)
}

func (g *Graph) GetSquare(point Point, distance float64) []*Node {
	buttomLeft, topRight := point.Move(-distance, -distance), point.Move(distance, distance)
	return g.Seg.GetInterval(buttomLeft.Latitude, topRight.Longitude, buttomLeft.Latitude, topRight.Latitude)
}

func (g *Graph) GetCircle(point Point, distance float64) (result []*Node) {
	candidates := g.GetSquare(point, distance)
	for _, node := range candidates {
		if node.Position.Distance(point) <= distance {
			result = append(result, node)
		}
	}
	return
}

type heapNode struct {
	node     *Node
	distance float64
}

func (g *Graph) Distance(start, end string, maxDuration float64, considerSpeed bool) (float64, error) {
	startNode, err := g.GetNode(start)
	if err != nil {
		return -1, err
	}

	priorityQueue := NewHeap(func(i, j interface{}) bool {
		if i.(heapNode).distance == j.(heapNode).distance {
			return i.(heapNode).node.ID < j.(heapNode).node.ID
		} else {
			return i.(heapNode).distance < j.(heapNode).distance
		}
	})
	visited := make(map[string]bool)
	dist := make(map[string]float64)

	dist[startNode.ID] = 0
	priorityQueue.Push(heapNode{node: startNode, distance: 0})
	for priorityQueue.Length() > 0 {
		current := priorityQueue.Pop().(heapNode)
		if current.distance > maxDuration {
			break
		}
		if visited[current.node.ID] {
			continue
		}
		visited[current.node.ID] = true

		if current.node.ID == end {
			return current.distance, nil
		}

		g.updateDistances(current, priorityQueue, visited, dist, considerSpeed)
	}

	return -1, ErrNodeNotReachable
}

func (g *Graph) updateDistances(current heapNode, priorityQueue *Heap, visited map[string]bool, dist map[string]float64, considerSpeed bool) {
	for _, edgeID := range current.node.OutEdges {
		edge := g.Edges[edgeID]
		neighbour, _ := g.GetNode(edge.End)
		if !visited[neighbour.ID] {
			distance := calculateDistance(current.distance, edge, considerSpeed)

			if current_distance, ok := dist[neighbour.ID]; !ok || distance < current_distance {
				dist[neighbour.ID] = distance
				priorityQueue.Push(heapNode{node: neighbour, distance: distance})
			}
		}
	}
}

func calculateDistance(currentDistance float64, edge *Edge, considerSpeed bool) (distance float64) {
	distance = currentDistance
	if considerSpeed {
		if edge.Speed > 0 {
			distance += edge.Length / edge.Speed
		} else {
			distance += edge.Length / AverageSpeed
		}
	} else {
		distance += edge.Length
	}
	return
}

func (e *Edge) LengthTo(point Point) (length float64) {
	intersect := point.ClosestPointOnEdge(e)

	for i := 0; i < len(e.Poly)-1; i++ {
		if intersect.IsOnSegment(e.Poly[i], e.Poly[i+1]) {
			length += intersect.Distance(e.Poly[i])
			break
		}
		length += e.Poly[i].Distance(e.Poly[i+1])
	}

	return length
}

func (e *Edge) LengthFrom(point Point) float64 {
	return e.Length - e.LengthTo(point)
}
