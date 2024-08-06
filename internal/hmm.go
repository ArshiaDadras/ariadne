package internal

import (
	"github.com/ArshiaDadras/Ariadne/pkg"
)

const (
	Sigma          = 5
	Beta           = 10
	MaxDistance    = 2000
	ConsiderSpeed  = true
	DistanceFactor = 3
)

func emmisionProbability(gpsPoint GPSPoint, edge *pkg.Edge) float64 {
	return EmmisionLogProbability(gpsPoint.Location.Distance(edge.Poly[0]), Sigma)
}

func transitionProbability(graph *pkg.Graph, edge1, edge2 *pkg.Edge, gpsPoint1, gpsPoint2 GPSPoint) float64 {
	point1 := gpsPoint1.Location
	point2 := gpsPoint2.Location
	distance1 := point1.Distance(point2)

	distance2, err := graph.Distance(edge1.End, edge2.Start, distance1+MaxDistance, ConsiderSpeed)
	if err != nil {
		return 0
	}
	distance2 += edge1.LengthFrom(point1.ClosestPointOnEdge(edge1)) + edge2.LengthTo(point2.ClosestPointOnEdge(edge2))
	if distance2 > DistanceFactor*gpsPoint1.TimeDifference(gpsPoint2) {
		return 0
	}

	return TransitionLogProbability(distance1, distance2, Beta)
}

func BestMatch(graph *pkg.Graph, points []GPSPoint) ([]pkg.Edge, error) {
	return nil, nil
}
