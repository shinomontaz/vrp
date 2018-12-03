package main

import (
	"sort"

	"github.com/shinomontaz/ga"
	"github.com/shinomontaz/vrp/types"
)

type RouteSet struct {
	List       map[int]*Route
	Wareheouse *types.LatLng
	Code       string
	Code2      []int
	fleet      []*types.Courier
	orders     []*types.Order
}

func (rs *RouteSet) Clone() ga.Individual {
	return rs
}

func (rs *RouteSet) Crossover(parent ga.Individual) ga.Individual {
	// перенумеровать курьеров!
	// вычсислить

	// по каждому маршруту вычислить среднюю координату всего множества заказов
	// получить по этой координате средний угол
	// после полного вычисления всех углов пройтись по углам в порядке увеличения
	// и присвоить каждой новому маршруту по порядку из курьеров

	listCenters := make(map[int]*types.LatLng, len(rs.List))

	for idx, route := range rs.List {
		listCenters[idx] = &types.LatLng{}
		for _, order := range route.List {
			listCenters[idx].Lat += order.Coords.Lat
			listCenters[idx].Lng += order.Coords.Lng
		}

		listCenters[idx].Lat /= float64(len(route.List))
		listCenters[idx].Lng /= float64(len(route.List))
	}

	listCentersKeys := make([]int, 0, len(rs.List))
	for idx := range listCenters {
		listCentersKeys = append(listCentersKeys, idx)
	}

	sort.Slice(listCentersKeys, func(i, j int) bool {
		return listCenters[i].Angle(rs.Wareheouse) > listCenters[j].Angle(rs.Wareheouse)
	})

	// новый порядок
	currCourier := 0
	for _, idx := range listCentersKeys {
		rs.List[idx].Courier = rs.fleet[currCourier]

		currCourier++
	}

	return rs
}

func (rs *RouteSet) Educate() {
}

func (rs *RouteSet) Fitness() float64 {
	fitness := 0.0
	for _, route := range rs.List {
		fitness += route.Cost()
	}

	return fitness
}

func (rs *RouteSet) UnFitness() float64 {
	unfitness := 0.0
	for _, route := range rs.List {
		unfitness += route.Violation()
	}

	return unfitness
}

func (rs *RouteSet) Mutate() ga.Individual {
	return rs
}

/*
func (rs *RouteSet) decodeCode() []*Route {
	result := make([]*Route, 0, 10)
	currCourier := -1
	var currRoute Route
	for i, courier := range rs.Code2 {
		if currCourier < 0 || currCourier != courier {
			currCourier = courier
			currRoute = Route{}
		}

		currRoute.List = append(currRoute.List)

	}
}
*/
