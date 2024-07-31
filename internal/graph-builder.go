package internal

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ArshiaDadras/Ariadne/models"
)

func ParseCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.Comma = '\t'

	data, err := reader.ReadAll()
	return data[1:], err
}

func parsePoints(pointStr string) []models.Point {
	points := make([]models.Point, 0)
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

		points = append(points, models.Point{
			Longitude: longitude,
			Latitude:  latitude,
		})
	}
	return points
}

func getOrCreateNode(graph *models.Graph, nodeID string, point models.Point) (*models.Node, error) {
	node, err := graph.GetNode(nodeID)
	if err != nil {
		node, err = graph.AddNode(nodeID, point)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

func BuildRoadNetwork(graph *models.Graph, path string) error {
	data, err := ParseCSV(path)
	if err != nil {
		return err
	}

	for _, row := range data {
		points := parsePoints(row[6])

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
		_, err = graph.AddEdge(row[0], start, end, speed, points)
		if err != nil {
			return err
		}

		if row[3] == "1" {
			_, err = graph.AddEdge(row[0]+"_reverse", end, start, speed, points)
			if err != nil {
				return err
			}
		}
	}

	graph.Preprocess()

	return nil
}
