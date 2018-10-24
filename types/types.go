package types

type LatLng struct {
	Lat float64
	Lng float64
}

func (p *LatLng) Distance(to *LatLng) float64 {
	return 0.0
}
