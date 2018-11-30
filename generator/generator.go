package generator

import (
	"math/rand"

	"github.com/shinomontaz/vrp/types"
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

func CreateOrders(n int) []*types.Order {
	res := make([]*types.Order, 0)
	for i := 0; i < n; i++ {
		res = append(res, &types.Order{
			Coords: &types.LatLng{
				Lat: rand.Float64() * 100,
				Lng: rand.Float64() * 100,
			},
			ID:     i,
			Weight: float64(rand.Intn(100)),
		})
	}
	return res
}

func CreateFleet(n int) []*types.Courier {
	res := make([]*types.Courier, 0)
	for i := 0; i < n; i++ {
		res = append(res, &types.Courier{
			ID:     i,
			Weight: 10.0 * float64(rand.Intn(100)),
		})
	}
	return res
}
