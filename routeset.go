package main

import (
	"math/rand"
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

	rs.renumber()
	parent.(*RouteSet).renumber()

	// now crossover!

	point1 := rand.Intn(len(rs.orders))
	point2 := point1 + rand.Intn(len(rs.orders)-point1)

	child := RouteSet{
		Wareheouse: rs.Wareheouse,
		orders:     rs.orders,
		fleet:      rs.fleet,
		Code2:      make([]int, len(rs.orders)),
		List:       make(map[int]*Route, len(rs.fleet)),
	}
	for carIdx, courier := range child.fleet {
		child.List[carIdx] = &Route{estimator: rs.List[1].estimator, Wareheouse: rs.Wareheouse, Courier: courier}
	}

	for i := 0; i < len(rs.orders); i++ {
		if i < point1 || i >= point2 {
			child.Code2 = append(child.Code2, rs.Code2[i])
		} else {
			child.Code2 = append(child.Code2, parent.(*RouteSet).Code2[i])
		}
	}

	// encode Code2 into List!
	for ordIdx, carIdx := range child.Code2 {
		child.List[carIdx].List = append(child.List[carIdx].List, child.orders[ordIdx])
	}

	return &child
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

func (rs *RouteSet) renumber() {

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
	// TODO: renumber Code2 also! VITAL!
}
