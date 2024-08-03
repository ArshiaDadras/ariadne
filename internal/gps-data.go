package internal

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

func ParseGPSData(path string) ([]pkg.Point, error) {
	data, err := ParseCSV(path)
	if err != nil {
		return nil, err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][0]+data[i][1] < data[j][0]+data[j][1]
	})

	points := make([]pkg.Point, 0, len(data))
	for _, row := range data {
		latitude, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}

		longitude, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, err
		}

		points = append(points, pkg.Point{
			Longitude: longitude,
			Latitude:  latitude,
		})
	}

	return points, nil
}

func MapMatch(graph *pkg.Graph, points []pkg.Point) ([]pkg.Edge, error) {
	nodes := graph.GetClosestNodes(points[0], 500)

	// m := make(map[pkg.Point]bool)
	// for _, node := range nodes {
	// 	m[node.Position] = true
	// }
	// others := make([]*pkg.Node, 0)
	// for _, node := range graph.Nodes {
	// 	if _, ok := m[node.Position]; !ok && node.Position.Distance(points[0]) < 1000 {
	// 		others = append(others, node)
	// 		m[node.Position] = true
	// 	}
	// }

	fmt.Println(points[0])
	// fmt.Println(len(nodes), len(others))
	SaveObject(nodes, "data/nodes.json")
	// SaveObject(others, "data/nodes.json")

	return nil, nil
}
