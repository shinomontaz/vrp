package main

import (
	"fmt"
	"math/rand"
	"sort"

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
		fleet:  fleet,
		orders: orders, // уже отсортированные
	}

	var ga = ga.Ga{
		NewIndividual:  Ifactory.Create,
		PopSize:        20,
		KeepRate:       0.3,
		TournamentSize: 10,
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
	return 0
}

func (rs *RouteSet) UnFitness() float64 {
	return 0
}

func (rs *RouteSet) Mutate() ga.Individual {
	return rs
}

type ScheduleFactory struct {
	fleet  []*types.Courier
	orders []*types.Order
}

func (sf *ScheduleFactory) Create() ga.Individual {
	rs := RouteSet{List: make(map[int]*Route, len(sf.fleet)), Code2: make([]int, 0, len(sf.orders))}

	currCourier := 0
	start := rand.Intn(len(sf.orders))
	rs.List[currCourier] = &Route{Courier: currCourier, List: make([]*types.Order, 10)}
	for i := start; i < len(sf.orders)+start; i++ {
		j := i % len(sf.orders)
		order := sf.orders[j]
		if !rs.List[currCourier].IsValid() {
			currCourier++
			rs.List[currCourier] = &Route{Courier: currCourier, List: make([]*types.Order, 10)}
		}
		rs.List[currCourier].List = append(rs.List[currCourier].List, order)
	}

	fmt.Println(rs)
	panic("!")
	return &rs
}

type Route struct {
	Courier int
	List    []*types.Order
}

func (r *Route) IsValid() bool {
	return false
}
