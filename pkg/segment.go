package pkg

import (
	"cmp"
	"slices"
	"sort"
)

type segment struct {
	Start  float64
	End    float64
	Left   *segment
	Right  *segment
	Values []*Edge
}

type SegmentNode struct {
	Point Point
	Edge  *Edge
}

func uniqueEdges(nodes []*SegmentNode) (edges []*Edge) {
	visited := make(map[*Edge]bool)
	for _, node := range nodes {
		if _, ok := visited[node.Edge]; !ok {
			visited[node.Edge] = true
			edges = append(edges, node.Edge)
		}
	}
	return
}

func merge(A, B []*Edge) (merged []*Edge) {
	for i, j := 0, 0; i < len(A) || j < len(B); {
		if j == len(B) || (i < len(A) && A[i].ID < B[j].ID) {
			if len(merged) == 0 || merged[len(merged)-1] != A[i] {
				merged = append(merged, A[i])
			}
			i++
		} else {
			if len(merged) == 0 || merged[len(merged)-1] != B[j] {
				merged = append(merged, B[j])
			}
			j++
		}
	}
	return
}

func (s *segment) getInterval(l, r float64) []*Edge {
	if r < s.Start || l > s.End {
		return nil
	}
	if l <= s.Start && r >= s.End {
		return s.Values
	}

	left, right := s.Left.getInterval(l, r), s.Right.getInterval(l, r)
	if left == nil {
		return right
	} else if right == nil {
		return left
	}
	return merge(left, right)
}

func build(sortedNodes []*SegmentNode, values []float64) (s *segment) {
	left, right := (*segment)(nil), (*segment)(nil)
	if len(values) > 1 {
		middle := len(values) >> 1
		lValues, rValues := values[:middle], values[middle:]
		middleIndex := sort.Search(len(sortedNodes), func(i int) bool {
			return sortedNodes[i].Point.Latitude >= values[middle]
		})
		lNodes, rNodes := sortedNodes[:middleIndex], sortedNodes[middleIndex:]
		left, right = build(lNodes, lValues), build(rNodes, rValues)
	}

	s = &segment{
		Start:  values[0],
		End:    values[len(values)-1],
		Values: uniqueEdges(sortedNodes),
		Left:   left,
		Right:  right,
	}
	return
}

func newSegment(nodes []*SegmentNode) *segment {
	values := make([]float64, 0)
	slices.SortFunc(nodes, func(a, b *SegmentNode) int {
		return cmp.Compare(a.Point.Latitude, b.Point.Latitude)
	})
	for _, node := range nodes {
		if len(values) == 0 || values[len(values)-1] != node.Point.Latitude {
			values = append(values, node.Point.Latitude)
		}
	}
	return build(nodes, values)
}

type Segment2D struct {
	Start float64
	End   float64
	Seg   *segment
	Left  *Segment2D
	Right *Segment2D
}

func (s *Segment2D) GetInterval(l1, r1, l2, r2 float64) []*Edge {
	if r1 < s.Start || l1 > s.End {
		return nil
	}
	if l1 <= s.Start && r1 >= s.End {
		return s.Seg.getInterval(l2, r2)
	}

	left, right := s.Left.GetInterval(l1, r1, l2, r2), s.Right.GetInterval(l1, r1, l2, r2)
	if left == nil {
		return right
	} else if right == nil {
		return left
	}
	return merge(left, right)
}

func build2D(sortedNodes []*SegmentNode, values []float64) (s *Segment2D) {
	left, right := (*Segment2D)(nil), (*Segment2D)(nil)
	if len(values) > 1 {
		middle := len(values) >> 1
		lValues, rValues := values[:middle], values[middle:]
		middleIndex := sort.Search(len(sortedNodes), func(i int) bool {
			return sortedNodes[i].Point.Longitude >= values[middle]
		})
		left, right = build2D(sortedNodes[:middleIndex], lValues), build2D(sortedNodes[middleIndex:], rValues)
	}

	s = &Segment2D{
		Start: values[0],
		End:   values[len(values)-1],
		Seg:   newSegment(sortedNodes),
		Left:  left,
		Right: right,
	}
	return
}

func NewSegment2D(nodes []*SegmentNode) *Segment2D {
	sortedNodes := make([]*SegmentNode, len(nodes))
	values := make([]float64, 0)
	copy(sortedNodes, nodes)

	slices.SortFunc(sortedNodes, func(a, b *SegmentNode) int {
		return cmp.Compare(a.Point.Longitude, b.Point.Longitude)
	})
	for _, node := range sortedNodes {
		if len(values) == 0 || values[len(values)-1] != node.Point.Longitude {
			values = append(values, node.Point.Longitude)
		}
	}
	return build2D(sortedNodes, values)
}

func (s *Segment2D) Get(point Point, distance float64) []*Edge {
	buttomLeft, topRight := point.Move(-distance, -distance), point.Move(distance, distance)
	return s.GetInterval(buttomLeft.Longitude, topRight.Longitude, buttomLeft.Latitude, topRight.Latitude)
}
