// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import (
	"cmp"
)

// ISolution defines the interface for genetic algorithm solutions.
// Implementations must provide a method to refresh fitness values.
type ISolution[T any] interface {
	// RefreshFitness updates the fitness value of the solution.
	// This method should recalculate the fitness based on the current chromosome.
	RefreshFitness()

	// DeepCopy creates a deep copy of the solution.
	DeepCopy() *ISolution[T]
}

// Solution represents a single solution in a genetic algorithm.
// It contains a chromosome (genetic material) and its associated fitness value.
type Solution[T cmp.Ordered] struct {
	// Chromosome represents the genetic material of the solution.
	// It's a slice of ordered values that define the solution's characteristics.
	Chromosome []T
	// Fitness represents the quality or performance of the solution.
	// Higher values typically indicate better solutions.
	Fitness float64
}

// DeepCopy creates a deep copy of the solution.
// The returned solution contains a deep copy of the chromosome and its fitness value.
//
// Returns:
//   - A pointer to the newly created Solution
func (s *Solution[T]) DeepCopy() *Solution[T] {
	return &Solution[T]{
		Chromosome: append([]T{}, s.Chromosome...),
		Fitness:    s.Fitness,
	}
}

// SolutionFactory provides factory methods for creating Solution instances.
// It supports generic types that implement the cmp.Ordered interface.
type SolutionFactory[T cmp.Ordered] struct{}

// NewSolutionFactory creates and returns a new SolutionFactory instance.
// The factory can create solutions of the specified generic type T.
func NewSolutionFactory[T cmp.Ordered]() *SolutionFactory[T] {
	return &SolutionFactory[T]{}
}

// CreateSolution creates a new Solution with the provided chromosome.
// The fitness is initialized to 0.0 and should be calculated separately.
//
// Parameters:
//   - chromosome: The genetic material for the new solution
//
// Returns:
//   - A pointer to the newly created Solution
func (sf *SolutionFactory[T]) CreateSolution(chromosome []T) *Solution[T] {
	return &Solution[T]{
		Chromosome: chromosome,
		Fitness:    0.0,
	}
}

// CreateRandomSolution creates a new Solution with a randomly generated chromosome.
// The chromosome is populated using the provided random generator function.
//
// Parameters:
//   - length: The desired length of the chromosome
//   - randomGen: A function that generates random values of type T
//
// Returns:
//   - A pointer to the newly created Solution with random chromosome
func (sf *SolutionFactory[T]) CreateRandomSolution(length int, randomGen func() T) *Solution[T] {
	chromosome := make([]T, length)
	for i := 0; i < length; i++ {
		chromosome[i] = randomGen()
	}
	return sf.CreateSolution(chromosome)
}

// CreateEmptySolution creates a new Solution with an empty chromosome.
// This is useful for initializing solutions that will be populated later.
//
// Returns:
//   - A pointer to the newly created Solution with empty chromosome
func (sf *SolutionFactory[T]) CreateEmptySolution() *Solution[T] {
	return &Solution[T]{
		Chromosome: []T{},
		Fitness:    0.0,
	}
}
