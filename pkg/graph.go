package pkg

import "errors"

var ErrNodeExists = errors.New("node already exists")
var ErrEdgeExists = errors.New("edge already exists")
var ErrNodeNotFound = errors.New("node not found")
var ErrEdgeNotFound = errors.New("edge not found")

type Node struct {
	ID       string   `json:"id"`
	Position Point    `json:"position"`
	InEdges  []string `json:"in_edges"`
	OutEdges []string `json:"out_edges"`
}

type Edge struct {
	ID    string  `json:"id"`
	Start string  `json:"start"`
	End   string  `json:"end"`
	Speed float64 `json:"speed"`
	Poly  []Point `json:"polygon"`
}

type Graph struct {
	Nodes map[string]*Node `json:"nodes"`
	Edges map[string]*Edge `json:"edges"`
	Seg   *Segment2D       `json:"-"`
}

func NewGraph() *Graph {
	graph := &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string]*Edge),
	}
	return graph
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

	edge := &Edge{
		ID:    id,
		Start: start.ID,
		End:   end.ID,
		Speed: speed,
		Poly:  poly,
	}
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
	// m := make(map[Point]string)
	// dupNodes, missEdges := 0, 0
	// for _, node := range g.Nodes {
	// 	if _, ok := m[node.Position]; ok {
	// 		dupNodes++
	// 	} else {
	// 		m[node.Position] = node.ID
	// 	}
	// }
	// for _, edge := range g.Edges {
	// 	for _, point := range edge.Poly {
	// 		if _, ok := m[point]; !ok {
	// 			m[point] = "dummy"
	// 			missEdges++
	// 		}
	// 	}
	// }
	// fmt.Println("Duplicate nodes:", dupNodes)
	// fmt.Println("Missing edges:", missEdges)

	nodes := make([]*Node, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}
	g.Seg = NewSegment2D(nodes)
}

func (g *Graph) GetClosestNodes(point Point, distance float64) []*Node {
	buttomLeft, topRight := point.Move(-distance, -distance), point.Move(distance, distance)
	candidates := g.Seg.GetInterval(buttomLeft.Longitude, topRight.Longitude, buttomLeft.Latitude, topRight.Latitude)

	result := make([]*Node, 0)
	for _, node := range candidates {
		if node.Position.Distance(point) <= distance {
			result = append(result, node)
		}
	}
	return result
}
