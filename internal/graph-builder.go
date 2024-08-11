package internal

import (
	"slices"
	"strconv"
	"strings"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

func parsePoints(pointStr string) (points []pkg.Point, err error) {
	for _, point := range strings.Split(pointStr[11:len(pointStr)-1], ", ") {
		coordinates := strings.Split(point, " ")
		longitude, err := strconv.ParseFloat(coordinates[0], 64)
		if err != nil {
			return nil, err
		}
		latitude, err := strconv.ParseFloat(coordinates[1], 64)
		if err != nil {
			return nil, err
		}

		points = append(points, pkg.Point{
			Longitude: longitude,
			Latitude:  latitude,
		})
	}
	return
}

func getOrCreateNode(graph *pkg.Graph, nodeID string, point pkg.Point, mp map[pkg.Point]string) (node *pkg.Node, err error) {
	if mp != nil {
		if id, ok := mp[point]; ok {
			nodeID = id
		} else {
			mp[point] = nodeID
		}
	}

	node, err = graph.GetNode(nodeID)
	if err != nil {
		node, err = graph.AddNode(nodeID, point)
		if err != nil {
			return nil, err
		}
	}
	return
}

func parseRow(row []string, graph *pkg.Graph, mp map[pkg.Point]string) (start, end *pkg.Node, speed float64, points []pkg.Point, err error) {
	points, err = parsePoints(row[6])
	if err != nil {
		return
	}

	start, err = getOrCreateNode(graph, row[1], points[0], mp)
	if err != nil {
		return
	}

	end, err = getOrCreateNode(graph, row[2], points[len(points)-1], mp)
	if err != nil {
		return
	}

	speed, err = strconv.ParseFloat(row[4], 64)
	if err != nil {
		return
	}
	speed *= 1000.0 / 3600.0

	return
}

func BuildRoadNetwork(graph *pkg.Graph, path string, removeDuplicates bool) error {
	data, err := ParseCSV(path)
	if err != nil {
		return err
	}

	var mp map[pkg.Point]string = nil
	if removeDuplicates {
		mp = make(map[pkg.Point]string)
	}

	for _, row := range data {
		start, end, speed, points, err := parseRow(row, graph, mp)
		if err != nil {
			return err
		}

		_, err = graph.AddEdge(row[0], start, end, speed, points)
		if err != nil {
			return err
		}

		if row[3] == "1" {
			slices.Reverse(points)
			_, err = graph.AddEdge(row[0]+"_reverse", end, start, speed, points)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Preprocess(graph *pkg.Graph) {
	maxLength := 2 * MaxCandidateDistance
	segmentNodes := make([]*pkg.SegmentNode, 0)
	for _, edge := range graph.Edges {
		for i := 1; i < len(edge.Poly); i++ {
			for start, end := edge.Poly[i-1], edge.Poly[i]; start.Distance(end) > maxLength; {
				start = start.MoveTowards(end, maxLength)
				segmentNodes = append(segmentNodes, &pkg.SegmentNode{
					Point: start,
					Edge:  edge,
				})
			}
		}

		for _, point := range edge.Poly {
			segmentNodes = append(segmentNodes, &pkg.SegmentNode{
				Point: point,
				Edge:  edge,
			})
		}
	}

	graph.Seg = pkg.NewSegment2D(segmentNodes)
}
