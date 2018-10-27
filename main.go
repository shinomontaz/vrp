package main

import (
	"fmt"
	"math/rand"
	"sort"

	"./generator"
	"./types"
	"github.com/shinomontaz/ga"
	//	"github.com/shinomontaz/ga"
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
		return orders[i].Coords.Angle(&wh) < orders[j].Coords.Angle(&wh)
	})

	fmt.Println(orders)

	for _, ord := range orders {
		fmt.Println("angle for order", ord.Coords.Angle(&wh), " - ", *ord.Coords)
	}

	panic("!")

	Ifactory := &ScheduleFactory{
		fleet:  fleet,
		orders: orders,
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

type RouteSet map[*types.Order]*types.Courier

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
	rs := make(map[*types.Order]*types.Courier, len(sf.orders))

	for _, order := range sf.orders {
		randFleet := rand.Intn(len(sf.fleet))
		rs[order] = sf.fleet[randFleet]
	}

	res := RouteSet(rs)

	return &res
}
