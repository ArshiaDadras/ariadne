package main

import (
	"log"

	"github.com/ArshiaDadras/Ariadne/internal"
	"github.com/ArshiaDadras/Ariadne/pkg"
)

func main() {
	graph := pkg.NewGraph()
	if err := internal.BuildRoadNetwork(graph, "data/road_network.csv", true); err != nil {
		log.Fatalf("Error building road network: %v", err)
	}
	log.Println("Graph created successfully")

	internal.Preprocess(graph)
	log.Println("Graph preprocessed successfully")

	points, err := internal.ParseGPSData("data/gps_data.csv")
	if err != nil {
		log.Fatalf("Error parsing GPS data: %v", err)
	}
	log.Println("GPS data parsed successfully")

	points = internal.RemoveNearbyPoints(points)
	log.Println("Nearby points removed successfully")

	edges, err := internal.MapMatch(graph, points)
	if err != nil {
		log.Fatalf("Error map matching: %v", err)
	}
	log.Println("Map matching completed successfully")

	if err := internal.SaveObject(points, "data/points.json"); err != nil {
		log.Fatalf("Error writing points: %v", err)
	}
	log.Println("Points written successfully")
	if err := internal.SaveObject(edges, "data/edges.json"); err != nil {
		log.Fatalf("Error writing edges: %v", err)
	}
	log.Println("Edges written successfully")
	if err := internal.SaveObject(graph, "data/graph.json"); err != nil {
		log.Fatalf("Error writing graph: %v", err)
	}
	log.Println("Graph written successfully")

	log.Println("All tasks completed successfully")
}
