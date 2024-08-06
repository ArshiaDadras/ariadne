package pkg

import (
	"math"
)

const (
	EarthRadius = 6378137
	Epsilon     = 1e-9
)

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (p *Point) Distance(other Point) float64 {
	lat1 := p.Latitude * math.Pi / 180
	long1 := p.Longitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	long2 := other.Longitude * math.Pi / 180

	dlat := lat2 - lat1
	dlong := long2 - long1

	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * EarthRadius
}

func (p *Point) Move(dx, dy float64) Point {
	return Point{
		Longitude: p.Longitude + (180 * dx / (math.Pi * EarthRadius * math.Cos(p.Latitude*math.Pi/180))),
		Latitude:  p.Latitude + (180 * dy / (math.Pi * EarthRadius)),
	}
}

func (p *Point) ClosestPointOnSegment(a, b Point) Point {
	ap := Point{Longitude: p.Longitude - a.Longitude, Latitude: p.Latitude - a.Latitude}
	ab := Point{Longitude: b.Longitude - a.Longitude, Latitude: b.Latitude - a.Latitude}
	t := (ap.Longitude*ab.Longitude + ap.Latitude*ab.Latitude) / (ab.Longitude*ab.Longitude + ab.Latitude*ab.Latitude)

	if t < 0 {
		return a
	} else if t > 1 {
		return b
	} else {
		return Point{
			Longitude: a.Longitude + ab.Longitude*t,
			Latitude:  a.Latitude + ab.Latitude*t,
		}
	}
}

func (p *Point) IsOnSegment(a, b Point) bool {
	return math.Abs(a.Distance(b)-(p.Distance(a)+p.Distance(b))) < Epsilon
}

func (p *Point) ClosestPointOnEdge(edge *Edge) Point {
	minDistance := math.Inf(1)
	var closestPoint Point

	for i := 0; i < len(edge.Poly)-1; i++ {
		closest := p.ClosestPointOnSegment(edge.Poly[i], edge.Poly[i+1])
		distance := p.Distance(closest)

		if distance < minDistance {
			minDistance = distance
			closestPoint = closest
		}
	}

	return closestPoint
}

func (p *Point) DistanceToEdge(edge *Edge) float64 {
	return p.Distance(p.ClosestPointOnEdge(edge))
}
