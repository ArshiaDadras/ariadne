package pkg

import (
	"math"
	"strconv"
)

const (
	EarthRadius = 6378137
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

func (p *Point) String() string {
	return "(" + strconv.FormatFloat(p.Longitude, 'f', -1, 64) + ", " + strconv.FormatFloat(p.Latitude, 'f', -1, 64) + ")"
}
