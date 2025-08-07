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
	fitnessFn  func(ch []T) (float64, error)
}

func (s *Solution[T]) RefreshFitness() error {
	res, err := s.fitnessFn(s.Chromosome)
	s.Fitness = res
	return err
}

type SolutionFactory[T cmp.Ordered] struct{}

func NewSolutionFactory[T cmp.Ordered]() *SolutionFactory[T] {
	return &SolutionFactory[T]{}
}

func (sf *SolutionFactory[T]) CreateSolution(chromosome []T, fitnessFn func(ch []T) (float64, error)) *Solution[T] {
	return &Solution[T]{
		Chromosome: chromosome,
		Fitness:    0.0,
		fitnessFn:  fitnessFn,
	}
}

func (sf *SolutionFactory[T]) CreateRandomSolution(length int, randomGen func() T, fitnessFn func(ch []T) (float64, error)) *Solution[T] {
	chromosome := make([]T, length)
	for i := 0; i < length; i++ {
		chromosome[i] = randomGen()
	}
	return sf.CreateSolution(chromosome, fitnessFn)
}

func (sf *SolutionFactory[T]) CreateEmptySolution(fitnessFn func(ch []T) (float64, error)) *Solution[T] {
	return &Solution[T]{
		Chromosome: []T{},
		Fitness:    0.0,
		fitnessFn:  fitnessFn,
	}
}
