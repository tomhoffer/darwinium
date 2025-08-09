package models

import "cmp"

type Population[T cmp.Ordered] struct {
	Individuals []Solution[T]
}

type PopulationFactory[T cmp.Ordered] struct{}

func NewPopulationFactory[T cmp.Ordered]() *PopulationFactory[T] {
	return &PopulationFactory[T]{}
}

func (pf *PopulationFactory[T]) CreatePopulation(individuals []Solution[T]) *Population[T] {
	return &Population[T]{
		Individuals: individuals,
	}
}

func (pf *PopulationFactory[T]) CreateEmptyPopulation() *Population[T] {
	return &Population[T]{
		Individuals: []Solution[T]{},
	}
}

func (pf *PopulationFactory[T]) CreateRandomPopulation(size, chromosomeLength int, solutionFactory *SolutionFactory[T], randomGen func() T, fitnessFn func(ch []T) (float64, error)) *Population[T] {
	individuals := make([]Solution[T], size)
	for i := 0; i < size; i++ {
		individuals[i] = *solutionFactory.CreateRandomSolution(chromosomeLength, randomGen)
	}
	return pf.CreatePopulation(individuals)
}
