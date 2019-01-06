package main

import (
	"math/rand"
	"sort"

	"./ega"
	"github.com/shinomontaz/vrp/types"
)

type RouteSet struct {
	List       map[int]*Route
	Wareheouse *types.LatLng
	Code       string
	Code2      []int
	fleet      []*types.Courier
	orders     []*types.Order
	fitness    float64
	unfitness  float64
}

func (rs *RouteSet) Clone() ega.Individual {
	return rs
}

func (rs *RouteSet) Crossover(parent ega.Individual) ega.Individual {
	// перенумеровать курьеров!
	// вычсислить

	// по каждому маршруту вычислить среднюю координату всего множества заказов
	// получить по этой координате средний угол
	// после полного вычисления всех углов пройтись по углам в порядке увеличения
	// и присвоить каждой новому маршруту по порядку из курьеров

	rs.renumber()
	parent.(*RouteSet).renumber()

	// fmt.Println("parent 1", rs.Code2, rs.Fitness(), rs.Unfitness())
	// fmt.Println("parent 2", parent.(*RouteSet).Code2, parent.Fitness(), parent.Unfitness())
	// drawRouteSet(fmt.Sprintf("parent1%v.png", rs.Code2), *rs)
	// drawRouteSet(fmt.Sprintf("parent2%v.png", parent.(*RouteSet).Code2), *parent.(*RouteSet))

	// now crossover!

	point1 := rand.Intn(len(rs.orders))
	point2 := point1 + rand.Intn(len(rs.orders)-point1)

	//	fmt.Println(point1, point2)

	child := RouteSet{
		Wareheouse: rs.Wareheouse,
		orders:     rs.orders,
		fleet:      rs.fleet,
		Code2:      make([]int, 0, len(rs.orders)),
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
	// drawRouteSet("child.png", child)

	// fmt.Println("child", child.Code2, child.Fitness(), child.Unfitness())

	// panic("!")
	return &child
}

func (rs *RouteSet) Educate() {
}

func (rs *RouteSet) Fitness() float64 { // Fitness less is better
	if rs.fitness != 0 {
		return rs.fitness
	}
	for _, route := range rs.List {
		rs.fitness += route.Cost()
	}

	return rs.fitness
}

func (rs *RouteSet) Unfitness() float64 { // Unfitness less is better
	if rs.unfitness != 0 {
		return rs.fitness
	}
	for _, route := range rs.List {
		rs.unfitness += route.Violation()
	}

	return rs.unfitness
}

func (rs *RouteSet) Mutate() ega.Individual {
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

	listCenters := make(map[int]*types.LatLng, len(rs.List)) // список центров масс маршрутов

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

	codeIndexMap := make(map[int][]int, len(rs.Code2))
	for i, cur := range rs.Code2 {
		codeIndexMap[cur] = append(codeIndexMap[cur], i)
	}

	for i, idx := range listCentersKeys {
		rs.List[idx].Courier = rs.fleet[i] // это правильно, т.к. в ключах listCentersKeys лежит новый порядок, а в значениях старый
		for _, orderIdx := range codeIndexMap[idx] {
			rs.Code2[orderIdx] = i
		}
		/*
			for _, order := range rs.List[idx].List {
				rs.Code2[order.ID] = currCourier
			}
		*/
	}
}
