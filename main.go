package main

import (
	"fmt"
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

	orders := generator.CreateOrders(10)
	fleet := generator.CreateFleet(100)
	wh := types.LatLng{
		Lat: 50,
		Lng: 50,
	}
	fmt.Println(orders)

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

type RouteSet struct {
	List  map[int]*Route
	Code  string
	Code2 []int
}

func (rs *RouteSet) Clone() ga.Individual {
	return rs
}

func (rs *RouteSet) Crossover(parent ga.Individual) ga.Individual {
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

type ScheduleFactory struct {
	fleet     []*types.Courier
	orders    []*types.Order
	estimator *Estimator
}

func (sf *ScheduleFactory) Create() ga.Individual {
	rs := RouteSet{List: make(map[int]*Route, len(sf.fleet)), Code2: make([]int, 0, len(sf.orders))}

	currCourier := 0
	start := rand.Intn(len(sf.orders))
	rs.List[currCourier] = &Route{Courier: sf.fleet[currCourier], List: make([]*types.Order, 0, 10), estimator: sf.estimator}
	for i := start; i < len(sf.orders)+start; i++ {
		j := i % len(sf.orders)
		order := sf.orders[j]
		if !rs.List[currCourier].IsValid() {
			currCourier++
			rs.List[currCourier] = &Route{Courier: sf.fleet[currCourier], List: make([]*types.Order, 0, 10), estimator: sf.estimator}
		}
		rs.List[currCourier].List = append(rs.List[currCourier].List, order)
		rs.Code2 = append(rs.Code2, currCourier)
		rs.Code += strconv.Itoa(currCourier)
	}

	for _, route := range rs.List {
		fmt.Println(route.List, route.Total(), route.Courier.Weight, route.Courier.ID)
	}

	return &rs
}

type Route struct {
	Courier   *types.Courier
	List      []*types.Order
	estimator *Estimator
}

func (r *Route) IsValid() bool {
	return r.Total() <= r.Courier.Weight
}

func (r *Route) Total() float64 {
	sum := 0.0
	for _, ord := range r.List {
		sum += ord.Weight
	}
	return sum
}

func (r *Route) Cost() float64 {
	r.SolveTsp()
	return 1.0 / (r.estimator.Cost(r) + 1.0)
}

func (r *Route) Violation() float64 {
	// unfitness
	if !r.IsValid() {
		return ((r.Total() - r.Courier.Weight) / r.Courier.Weight)
	}

	return 0
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

func (r *Route) SolveTsp() {
	currCost := r.estimator.Cost(r)
	for i := 1; i < len(r.List); i++ { // keep first item in place and do not swap it with any
		for j := 1; j < len(r.List); j++ {
			r.List[i], r.List[j] = r.List[j], r.List[i]
			newCost := r.estimator.Cost(r)
			if newCost < currCost {
				currCost = newCost
			} else { // revert transposition here
				r.List[i], r.List[j] = r.List[j], r.List[i]
			}
		}
	}
}
