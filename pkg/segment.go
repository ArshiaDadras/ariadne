package pkg

import (
	"sort"
)

type segment struct {
	Start float64
	End   float64
	Nodes []*Node
	Left  *segment
	Right *segment
}

func (s *segment) getInterval(l, r float64) []*Node {
	if r < s.Start || l > s.End {
		return nil
	}
	if l <= s.Start && r >= s.End {
		return s.Nodes
	}

	left := s.Left.getInterval(l, r)
	right := s.Right.getInterval(l, r)
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if len(left) < len(right) {
		left, right = right, left
	}
	return append(left, right...)
}

func build(nodes []*Node, values []float64) (s *segment) {
	if len(nodes) == 0 {
		return nil
	}
	s = &segment{Nodes: nodes, Start: values[0], End: values[len(values)-1]}
	if len(values) == 1 {
		return
	}

	median := values[(len(values)-1)>>1]
	lNodes, rNodes := make([]*Node, 0), make([]*Node, 0)
	lValues, rValues := make([]float64, 0), make([]float64, 0)
	for _, node := range nodes {
		if node.Position.Latitude <= median {
			lNodes = append(lNodes, node)
			if len(lValues) == 0 || lValues[len(lValues)-1] != node.Position.Latitude {
				lValues = append(lValues, node.Position.Latitude)
			}
		} else {
			rNodes = append(rNodes, node)
			if len(rValues) == 0 || rValues[len(rValues)-1] != node.Position.Latitude {
				rValues = append(rValues, node.Position.Latitude)
			}
		}
	}

	s.Left, s.Right = build(lNodes, lValues), build(rNodes, rValues)
	return
}

func newSegment(nodes []*Node) *segment {
	sortedNodes := make([]*Node, len(nodes))
	values := make([]float64, 0)
	copy(sortedNodes, nodes)

	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].Position.Latitude < sortedNodes[j].Position.Latitude
	})
	for _, node := range sortedNodes {
		if len(values) == 0 || values[len(values)-1] != node.Position.Latitude {
			values = append(values, node.Position.Latitude)
		}
	}
	return build(sortedNodes, values)
}

type Segment2D struct {
	Start float64
	End   float64
	Seg   *segment
	Left  *Segment2D
	Right *Segment2D
}

func (s *Segment2D) GetInterval(l1, r1, l2, r2 float64) []*Node {
	if r1 < s.Start || l1 > s.End {
		return nil
	}
	if l1 <= s.Start && r1 >= s.End {
		return s.Seg.getInterval(l2, r2)
	}

	left := s.Left.GetInterval(l1, r1, l2, r2)
	right := s.Right.GetInterval(l1, r1, l2, r2)
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if len(left) < len(right) {
		left, right = right, left
	}
	return append(left, right...)
}

func Build2D(nodes []*Node, values []float64) (s *Segment2D) {
	if len(nodes) == 0 {
		return nil
	}
	s = &Segment2D{Start: values[0], End: values[len(values)-1], Seg: newSegment(nodes)}
	if len(values) == 1 {
		return
	}

	median := values[(len(values)-1)>>1]
	lNodes, rNodes := make([]*Node, 0), make([]*Node, 0)
	lValues, rValues := make([]float64, 0), make([]float64, 0)
	for _, node := range nodes {
		if node.Position.Longitude <= median {
			lNodes = append(lNodes, node)
			if len(lValues) == 0 || lValues[len(lValues)-1] != node.Position.Longitude {
				lValues = append(lValues, node.Position.Longitude)
			}
		} else {
			rNodes = append(rNodes, node)
			if len(rValues) == 0 || rValues[len(rValues)-1] != node.Position.Longitude {
				rValues = append(rValues, node.Position.Longitude)
			}
		}
	}

	s.Left, s.Right = Build2D(lNodes, lValues), Build2D(rNodes, rValues)
	return
}

func NewSegment2D(nodes []*Node) *Segment2D {
	sortedNodes := make([]*Node, len(nodes))
	values := make([]float64, 0)
	copy(sortedNodes, nodes)

	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].Position.Longitude < sortedNodes[j].Position.Longitude
	})
	for _, node := range sortedNodes {
		if len(values) == 0 || values[len(values)-1] != node.Position.Longitude {
			values = append(values, node.Position.Longitude)
		}
	}
	return Build2D(sortedNodes, values)
}
