package pkg

import (
	"math"
)

const (
	EarthRadius = 6378137
	Epsilon     = 1e-6
)

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func toDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func (p *Point) Distance(other Point) float64 {
	lat1, long1 := toRadians(p.Latitude), toRadians(p.Longitude)
	lat2, long2 := toRadians(other.Latitude), toRadians(other.Longitude)

	dlat, dlong := lat2-lat1, long2-long1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * EarthRadius
}

func (p *Point) Move(dx, dy float64) Point {
	return Point{
		Longitude: p.Longitude + toDegrees(dx/(EarthRadius*math.Cos(toRadians(p.Latitude)))),
		Latitude:  p.Latitude + toDegrees(dy/EarthRadius),
	}
}

func (p *Point) ClosestPointOnSegment(a, b Point) Point {
	latP, longP := toRadians(p.Latitude), toRadians(p.Longitude)
	latA, longA := toRadians(a.Latitude), toRadians(a.Longitude)
	latB, longB := toRadians(b.Latitude), toRadians(b.Longitude)

	xP, yP, zP := math.Cos(latP)*math.Cos(longP), math.Cos(latP)*math.Sin(longP), math.Sin(latP)
	xA, yA, zA := math.Cos(latA)*math.Cos(longA), math.Cos(latA)*math.Sin(longA), math.Sin(latA)
	xB, yB, zB := math.Cos(latB)*math.Cos(longB), math.Cos(latB)*math.Sin(longB), math.Sin(latB)

	ABx, ABy, ABz := xB-xA, yB-yA, zB-zA
	APx, APy, APz := xP-xA, yP-yA, zP-zA

	abAb := ABx*ABx + ABy*ABy + ABz*ABz
	apAb := APx*ABx + APy*ABy + APz*ABz
	projFactor := apAb / abAb
	if projFactor < 0 {
		return a
	} else if projFactor > 1 {
		return b
	}

	x, y, z := xA+projFactor*ABx, yA+projFactor*ABy, zA+projFactor*ABz
	lat, long := math.Asin(z), math.Atan2(y, x)

	return Point{
		Latitude:  toDegrees(lat),
		Longitude: toDegrees(long),
	}
}

func (p *Point) IsOnSegment(a, b Point) bool {
	return math.Abs(a.Distance(b)-(p.Distance(a)+p.Distance(b))) < Epsilon
}

func (p *Point) ClosestPointOnEdge(edge *Edge) (closestPoint Point) {
	minDistance := math.Inf(1)
	for i := 0; i < len(edge.Poly)-1; i++ {
		closest := p.ClosestPointOnSegment(edge.Poly[i], edge.Poly[i+1])
		distance := p.Distance(closest)

		if distance < minDistance {
			minDistance = distance
			closestPoint = closest
		}
	}
	return
}

func (p *Point) DistanceToEdge(edge *Edge) float64 {
	return p.Distance(p.ClosestPointOnEdge(edge))
}

func (a *Point) MoveTowards(b Point, d float64) Point {
	latA, longA := toRadians(a.Latitude), toRadians(a.Longitude)
	latB, longB := toRadians(b.Latitude), toRadians(b.Longitude)

	latDiff, longDiff := latB-latA, longB-longA
	lat, long := latA+latDiff*d/a.Distance(b), longA+longDiff*d/a.Distance(b)

	return Point{
		Latitude:  toDegrees(lat),
		Longitude: toDegrees(long),
	}
}
