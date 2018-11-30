package main

type Estimator struct {
}

func (e *Estimator) Cost(r *Route) float64 {
	cost := 0.0
	if len(r.List) == 0 {
		return cost
	}
	cost += r.Wareheouse.Distance(r.List[0].Coords)
	for i := 1; i < len(r.List); i++ {
		cost += r.List[i-1].Coords.Distance(r.List[i].Coords)
	}
	cost += r.List[len(r.List)-1].Coords.Distance(r.List[0].Coords)

	return cost
}
