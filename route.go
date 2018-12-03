package main

import "github.com/shinomontaz/vrp/types"

type Route struct {
	Courier    *types.Courier
	List       []*types.Order
	Wareheouse *types.LatLng
	estimator  *Estimator
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

func (r *Route) SolveTsp() {
	currCost := r.estimator.Cost(r)
	for i := 0; i < len(r.List); i++ { // keep first item in place and do not swap it with any
		for j := 0; j < len(r.List); j++ {
			if i == j {
				continue
			}
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
