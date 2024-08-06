package internal

import (
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

func parsePoints(pointStr string) (points []pkg.Point) {
	for _, point := range strings.Split(pointStr[11:len(pointStr)-1], ", ") {
		coordinates := strings.Split(point, " ")
		longitude, err := strconv.ParseFloat(coordinates[0], 64)
		if err != nil {
			log.Fatalf("Error parsing longitude: %v", err)
		}
		latitude, err := strconv.ParseFloat(coordinates[1], 64)
		if err != nil {
			log.Fatalf("Error parsing latitude: %v", err)
		}

		points = append(points, pkg.Point{
			Longitude: longitude,
			Latitude:  latitude,
		})
	}
	return
}

func getOrCreateNode(graph *pkg.Graph, nodeID string, point pkg.Point) (node *pkg.Node, err error) {
	node, err = graph.GetNode(nodeID)
	if err != nil {
		node, err = graph.AddNode(nodeID, point)
		if err != nil {
			return nil, err
		}
	}
	return
}

func BuildRoadNetwork(graph *pkg.Graph, path string, removeDuplicates bool) error {
	data, err := ParseCSV(path)
	if err != nil {
		return err
	}

	mp := make(map[pkg.Point]string, 0)
	for _, row := range data {
		points := parsePoints(row[6])
		if removeDuplicates {
			if _, ok := mp[points[0]]; ok {
				row[1] = mp[points[0]]
			} else {
				mp[points[0]] = row[1]
			}
			if _, ok := mp[points[len(points)-1]]; ok {
				row[2] = mp[points[len(points)-1]]
			} else {
				mp[points[len(points)-1]] = row[2]
			}
		}

		start, err := getOrCreateNode(graph, row[1], points[0])
		if err != nil {
			return err
		}

		end, err := getOrCreateNode(graph, row[2], points[len(points)-1])
		if err != nil {
			return err
		}

		speed, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return err
		}
		speed *= 1000.0 / 3600.0

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
