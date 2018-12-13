package ega

import (
	"math/rand"
	"sort"
)

type IndividualFactory func() Individual

type Ega struct {
	NewIndividual IndividualFactory
	PopSize       int

	Population []Individual
	Best       Individual

	CreateRate     float64
	KeepRate       float64
	TournamentSize int
}

func (g *Ega) Initialize() {
	// create initial population

	if g.TournamentSize == 0 {
		g.TournamentSize = 3
	}

	g.Population = make([]Individual, 0, g.PopSize)
	for i := 0; i < g.PopSize; i++ {
		g.Population = append(g.Population, g.NewIndividual())
	}

	g.Best = g.Population[0]
}

func (g *Ega) pick() Individual {
	return g.tournamentSelection()
}

func (g *Ega) tournamentSelection() Individual {
	var best Individual
	for i := 0; i < g.TournamentSize; i++ {
		inst := g.Population[rand.Intn(len(g.Population))]
		if best == nil || inst.Fitness() > best.Fitness() {
			best = inst
		}
	}
	return best
}

func (g *Ega) Evolve() {
	parent1 := g.pick()
	parent2 := g.pick()
	child := parent1.Crossover(parent2).Mutate()
	child.Educate()

	// выбрать члена популяции, которого надо заменить
	set1, set2, set3 := make(map[int]Individual, g.PopSize), make(map[int]Individual, g.PopSize), make(map[int]Individual, g.PopSize)

	to_replace1, to_replace2, to_replace3 := struct {
		idx       int
		unfitness float64
	}{}, struct {
		idx       int
		unfitness float64
	}{}, struct {
		idx       int
		unfitness float64
	}{}

	for idx, ind := range g.Population {
		if ind.Fitness() <= child.Fitness() && ind.Unfitness() >= child.Unfitness() {
			set1[idx] = ind
			if to_replace1.unfitness < ind.Unfitness() {
				to_replace1.idx = idx
				to_replace1.unfitness = ind.Unfitness()
			}
		}
		if ind.Fitness() > child.Fitness() && ind.Unfitness() >= child.Unfitness() {
			set2[idx] = ind
			if to_replace2.unfitness < ind.Unfitness() {
				to_replace2.idx = idx
				to_replace2.unfitness = ind.Unfitness()
			}
		}
		if ind.Fitness() <= child.Fitness() && ind.Unfitness() < child.Unfitness() {
			set3[idx] = ind
			if to_replace3.unfitness < ind.Unfitness() {
				to_replace3.idx = idx
				to_replace3.unfitness = ind.Unfitness()
			}
		}
	}

	if len(set1) >= 0 {
		g.Population[to_replace1.idx] = child
	}
	if len(set2) >= 0 {
		g.Population[to_replace2.idx] = child
	}
	if len(set3) >= 0 {
		g.Population[to_replace3.idx] = child
	}

}

func (g *Ega) Record() Individual {
	sort.Slice(g.Population, func(i, j int) bool {
		return g.Population[i].Fitness() > g.Population[j].Fitness() // it is a "less" function, so we need bigger first
	})

	return g.Population[0]
}
