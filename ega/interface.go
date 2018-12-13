package ega

type Individual interface {
	Fitness() float64
	Unfitness() float64
	Mutate() Individual
	Crossover(parner Individual) (child Individual)
	Clone() Individual
	Educate()
}
