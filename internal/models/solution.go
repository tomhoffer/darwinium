package models

import (
	"cmp"
)

type ISolution[T any] interface {
	RefreshFitness()
}

type Solution[T cmp.Ordered] struct {
	Chromosome []T
	Fitness    float64
}

type SolutionFactory[T cmp.Ordered] struct{}

func NewSolutionFactory[T cmp.Ordered]() *SolutionFactory[T] {
	return &SolutionFactory[T]{}
}

func (sf *SolutionFactory[T]) CreateSolution(chromosome []T) *Solution[T] {
	return &Solution[T]{
		Chromosome: chromosome,
		Fitness:    0.0,
	}
}

func (sf *SolutionFactory[T]) CreateRandomSolution(length int, randomGen func() T) *Solution[T] {
	chromosome := make([]T, length)
	for i := 0; i < length; i++ {
		chromosome[i] = randomGen()
	}
	return sf.CreateSolution(chromosome)
}

func (sf *SolutionFactory[T]) CreateEmptySolution() *Solution[T] {
	return &Solution[T]{
		Chromosome: []T{},
		Fitness:    0.0,
	}
}
