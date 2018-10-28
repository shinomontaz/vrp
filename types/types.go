package types

import (
	"math"
)

type LatLng struct {
	Lat float64
	Lng float64
}

type Courier struct {
	Weight float64
	ID     int
}

type Order struct {
	Coords *LatLng
	ID     int
	Weight float64
}

func (from *LatLng) Distance(to *LatLng) float64 {
	dLat := (from.Lat - to.Lat) * math.Pi / 180.0
	dLng := (from.Lng - to.Lng) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(to.Lat*math.Pi/180.0)*math.Cos(from.Lat*math.Pi/180.0)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := 3958.75 * c
	return dist * 1609
}

func (p *LatLng) Angle(center *LatLng) float64 {
	//	degrees := math.Acos(((p.Lat - center.Lat) + (p.Lng - center.Lng)) / (math.Sqrt(math.Pow((p.Lat-center.Lat), 2)) + math.Pow((p.Lng-center.Lng), 2)))

	x := (p.Lng - center.Lng)
	degrees := math.Asin(x / math.Sqrt(math.Pow((p.Lat-center.Lat), 2)+math.Pow((p.Lng-center.Lng), 2)))
	/*
		if degrees > 180 {
			degrees = 360 + degrees
		}
	*/
	return degrees * math.Pi / 180
}
