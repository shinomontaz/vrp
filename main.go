package main

import (
	"./generator"
	"github.com/shinomontaz/ga"
	//	"github.com/shinomontaz/ga"
)

func main() {
	// Сортируем точки по значению полярного угла
	// Начало генетики:
	// Особь: Начинаем набирать машины в порядке отсортированых точек. Последнюю точку включаем или не включаем
	// Перед вычислением фитнеса - 2-opt

	points := generator.CreatePoints(2000)

}

type Object struct {
	core map[int]int
}

type ScheduleFactory struct {
}

func (sf *ScheduleFactory) Create() ga.Individual {
	//

	return newSchedule
}
