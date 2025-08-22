// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import "cmp"

// Population represents a collection of solutions (individuals) in a genetic algorithm.
// It contains multiple Solution instances that form the current generation.
type Population[T cmp.Ordered] struct {
	// Individuals is a slice of Solution instances that make up the population.
	// Each individual represents a potential solution to the optimization problem.
	Individuals []Solution[T]
}

// BestSolution finds and returns the individual with the highest fitness in the population.
// If the population is empty, it returns an error.
func (p *Population[T]) BestSolution() (*Solution[T], error) {
	if p == nil || len(p.Individuals) == 0 {
		return nil, ErrPopulationEmpty
	}

	best := p.Individuals[0]
	for i := 1; i < len(p.Individuals); i++ {
		if p.Individuals[i].Fitness > best.Fitness {
			best = p.Individuals[i]
		}
	}

	return &best, nil
}

// BestFitness returns the fitness of the best individual in the population.
// If the population is empty, it returns an error.
func (p *Population[T]) BestFitness() (float64, error) {
	bestSolution, err := p.BestSolution()
	if err != nil {
		return 0, err
	}
	return bestSolution.Fitness, nil
}

// PopulationFactory provides factory methods for creating Population instances.
// It supports generic types that implement the cmp.Ordered interface.
type PopulationFactory[T cmp.Ordered] struct{}

// NewPopulationFactory creates and returns a new PopulationFactory instance.
// The factory can create populations of the specified generic type T.
func NewPopulationFactory[T cmp.Ordered]() *PopulationFactory[T] {
	return &PopulationFactory[T]{}
}

// CreatePopulation creates a new Population with the provided individuals.
// The population is initialized with the given slice of solutions.
//
// Parameters:
//   - individuals: A slice of Solution instances to populate the population
//
// Returns:
//   - A pointer to the newly created Population
func (pf *PopulationFactory[T]) CreatePopulation(individuals []Solution[T]) *Population[T] {
	return &Population[T]{
		Individuals: individuals,
	}
}

// CreateEmptyPopulation creates a new Population with no individuals.
// This is useful for initializing populations that will be populated later.
//
// Returns:
//   - A pointer to the newly created Population with no individuals
func (pf *PopulationFactory[T]) CreateEmptyPopulation() *Population[T] {
	return &Population[T]{
		Individuals: []Solution[T]{},
	}
}

// CreateRandomPopulation creates a new Population with randomly generated individuals.
// Each individual in the population is created using the provided solution factory
// and random generator function. The fitness of each individual is calculated
// using the provided fitness function.
//
// Parameters:
//   - size: The number of individuals to create in the population
//   - chromosomeLength: The length of each individual's chromosome
//   - solutionFactory: Factory for creating solution instances
//   - randomGen: Function that generates random values of type T
//   - fitnessFn: Function to calculate fitness for each individual's chromosome
//
// Returns:
//   - A pointer to the newly created Population with random individuals
func (pf *PopulationFactory[T]) CreateRandomPopulation(size, chromosomeLength int, solutionFactory *SolutionFactory[T], randomGen func() T, fitnessFn func(ch []T) (float64, error)) *Population[T] {
	individuals := make([]Solution[T], size)
	for i := 0; i < size; i++ {
		individuals[i] = *solutionFactory.CreateRandomSolution(chromosomeLength, randomGen)
	}
	return pf.CreatePopulation(individuals)
}
