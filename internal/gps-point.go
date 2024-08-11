package internal

import (
	"slices"
	"time"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

const (
	MaxNearby = 2 * Sigma
)

type GPSPoint struct {
	Location pkg.Point
	Time     time.Time
}

func (p *GPSPoint) Distance(other GPSPoint) float64 {
	return p.Location.Distance(other.Location)
}

func (p *GPSPoint) TimeDifference(other GPSPoint) float64 {
	return p.Time.Sub(other.Time).Seconds()
}

func MapMatch(graph *pkg.Graph, points []GPSPoint) (match []*pkg.Edge, err error) {
	match, err = BestMatch(graph, points)
	slices.Reverse(match)
	return
}

func RemoveNearbyPoints(points []GPSPoint) (result []GPSPoint) {
	for i := 0; i < len(points); i++ {
		if i == 0 || points[i].Distance(points[i-1]) >= MaxNearby {
			result = append(result, points[i])
		}
	}
	return
}
