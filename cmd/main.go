package main

import (
	"log"

	"github.com/ArshiaDadras/Ariadne/internal"
	"github.com/ArshiaDadras/Ariadne/models"
)

func main() {
	graph := models.NewGraph()
	if err := internal.BuildRoadNetwork(graph, "data/road_network.csv"); err != nil {
		log.Fatalf("Error building road network: %v", err)
	}
	log.Println("Graph created successfully")

	graph.Preprocess()
	log.Println("Graph preprocessed successfully")

	points, err := internal.ParseGPSData("data/gps_data.csv")
	if err != nil {
		log.Fatalf("Error parsing GPS data: %v", err)
	}
	log.Println("GPS data parsed successfully")

	edges, err := internal.MapMatch(graph, points)
	if err != nil {
		log.Fatalf("Error map matching: %v", err)
	}
	log.Println("Map matching completed successfully")

	if err := internal.SaveObject(edges, "data/edges.json"); err != nil {
		log.Fatalf("Error writing edges: %v", err)
	}
	log.Println("Edges written successfully")

	if err := internal.SaveObject(graph, "data/graph.json"); err != nil {
		log.Fatalf("Error writing graph: %v", err)
	}
	log.Println("Graph written successfully")

	log.Println("All tasks completed successfully")

	// P1 := models.Point{Longitude: 51.372723607382994, Latitude: 35.792019008390966}
	// P2 := models.Point{Longitude: 51.371323391458816, Latitude: 35.79170002313525}
	// P3 := models.Point{Longitude: 51.37272618815106, Latitude: 35.79319153141927}

	// fmt.Println("Distance between P1 and P2:", P1.Distance(P2), "meters")
	// fmt.Println("Distance between P1 and P3:", P1.Distance(P3), "meters")
}
