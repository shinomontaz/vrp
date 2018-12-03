package main

import (
	"math/rand"
	"sort"
	"strconv"

	"github.com/shinomontaz/ga"
	"github.com/shinomontaz/vrp/generator"
	"github.com/shinomontaz/vrp/types"
)

func main() {
	// Сортируем точки по значению полярного угла
	// Начало генетики:
	// Особь: Начинаем набирать машины в порядке отсортированых точек. Последнюю точку включаем или не включаем
	// Перед вычислением фитнеса - 2-opt

	orders := generator.CreateOrders(20)
	fleet := generator.CreateFleet(10)
	wh := types.LatLng{
		Lat: 50,
		Lng: 50,
	}

	// сортировка заказов по полярному углу
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Coords.Angle(&wh) > orders[j].Coords.Angle(&wh)
	})
	// for _, ord := range orders {
	// 	fmt.Println("angle for order", ord.ID, ord.Coords.Angle(&wh), " - ", *ord.Coords)
	// }

	//	drawOrders("result.png", orders, &wh)

	Ifactory := &ScheduleFactory{
		fleet:     fleet,
		orders:    orders, // уже отсортированные
		estimator: &Estimator{},
		warehouse: &wh,
	}

	var ga = ga.Ga{
		NewIndividual:  Ifactory.Create,
		PopSize:        50,
		KeepRate:       0,
		TournamentSize: 2,
	}

	ga.Initialize()

	bestEver := ga.Population[0].Clone()
	bestFitness := bestEver.Fitness()

	for i := 0; i < 10000; i++ {
		currBest := ga.Record()
		if bestFitness < currBest.Fitness() {
			bestEver = currBest
			bestFitness = bestEver.Fitness()
		}
		ga.Evolve()
	}

}

type ScheduleFactory struct {
	fleet     []*types.Courier
	orders    []*types.Order
	warehouse *types.LatLng
	estimator *Estimator
}

func (sf *ScheduleFactory) Create() ga.Individual {
	rs := RouteSet{List: make(map[int]*Route, len(sf.fleet)), Code2: make([]int, 0, len(sf.orders)), Wareheouse: sf.warehouse, fleet: sf.fleet}
	currCourier := 0
	start := rand.Intn(len(sf.orders))
	rs.List[currCourier] = &Route{Courier: sf.fleet[currCourier], List: make([]*types.Order, 0, 10), estimator: sf.estimator, Wareheouse: sf.warehouse}
	for i := start; i < len(sf.orders)+start; i++ {
		j := i % len(sf.orders)
		order := sf.orders[j]
		if !rs.List[currCourier].IsValid() {
			currCourier++
			rs.List[currCourier] = &Route{Courier: sf.fleet[currCourier], List: make([]*types.Order, 0, 10), estimator: sf.estimator, Wareheouse: sf.warehouse}
		}
		rs.List[currCourier].List = append(rs.List[currCourier].List, order)
		rs.Code2 = append(rs.Code2, currCourier)
		rs.Code += strconv.Itoa(currCourier)
	}

	// for _, route := range rs.List {
	// 	fmt.Println(route.List, route.Total(), route.Courier.Weight, route.Courier.ID)
	// }

	return &rs
}
