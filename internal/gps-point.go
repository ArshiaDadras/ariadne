package internal

import (
	"time"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

const (
	TimeFormat = "02-Jan-2006 15:04:05"
	MaxGap     = 1000
	MaxBreak   = 180
	MaxNearby  = 2 * Sigma
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

func MapMatch(graph *pkg.Graph, points []GPSPoint) ([]pkg.Edge, error) {
	for i := 0; i < len(points)-1; i++ {
		j := 0
		for i-j >= 0 && i+j+1 < len(points) && points[i-j].Distance(points[i+j+1]) > MaxGap && points[i-j].TimeDifference(points[i+j+1]) <= MaxBreak {
			j++
		}

		if i-j < 0 || i+j+1 >= len(points) || points[i-j].TimeDifference(points[i+j+1]) > MaxBreak {
			match1, err := MapMatch(graph, points[:i])
			if err != nil {
				return nil, err
			}

			match2, err := MapMatch(graph, points[i+1:])
			if err != nil {
				return nil, err
			}

			return append(match1, match2...), nil
		} else if j > 0 {
			points = append(points[:i-j], points[i+j+1:]...)
		}
	}

	return BestMatch(graph, points)
}

func RemoveNearbyPoints(points []GPSPoint) (result []GPSPoint) {
	for i := 0; i < len(points); i++ {
		if i == 0 || points[i].Distance(points[i-1]) >= MaxNearby {
			result = append(result, points[i])
		}
	}
	return
}
