package generator

import (
	"math/rand"

	"../types"
)

func CreatePoints(n int) []*types.LatLng {
	res := make([]*types.LatLng, 0)
	for i := 0; i < n; i++ {
		res = append(res, &types.LatLng{
			Lat: rand.Float64() * 100,
			Lng: rand.Float64() * 100,
		})
	}
	return res
}
