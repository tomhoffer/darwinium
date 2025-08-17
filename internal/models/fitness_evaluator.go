// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import (
	"cmp"

	"github.com/tomhoffer/darwinium/internal/utils"
)

// IFitnessEvaluator defines the interface for fitness evaluation in genetic algorithms.
// Implementations must provide a method to evaluate the fitness of a chromosome.
type IFitnessEvaluator[T cmp.Ordered] interface {
	// Evaluate calculates the fitness value of a given chromosome.
	// Higher fitness values typically indicate better solutions.
	//
	// Parameters:
	//   - chromosome: The genetic material to evaluate
	//
	// Returns:
	//   - fitness: The calculated fitness value
	//   - error: Any error that occurred during evaluation
	Evaluate(chromosome []T) (float64, error)
}

// SimpleSumFitnessEvaluator implements a basic fitness evaluator that calculates
// fitness as the sum of all values in the chromosome. This is a simple example
// implementation that can be used for testing or as a baseline.
type SimpleSumFitnessEvaluator[T cmp.Ordered] struct{}

// NewSimpleSumFitnessEvaluator creates and returns a new SimpleSumFitnessEvaluator instance.
// This evaluator can work with any generic type T that implements cmp.Ordered.
func NewSimpleSumFitnessEvaluator[T cmp.Ordered]() *SimpleSumFitnessEvaluator[T] {
	return &SimpleSumFitnessEvaluator[T]{}
}

// Evaluate calculates the fitness of a chromosome by summing all its values.
// Each value in the chromosome is converted to float64 before summing.
// This method implements the IFitnessEvaluator interface.
//
// Parameters:
//   - chromosome: The genetic material to evaluate
//
// Returns:
//   - float64: The sum of all chromosome values as the fitness
//   - error: Any error that occurred during conversion or calculation
func (s SimpleSumFitnessEvaluator[T]) Evaluate(chromosome []T) (float64, error) {
	var sum float64
	for _, val := range chromosome {
		converted, err := utils.ConvertToFloat64(val)
		if err != nil {
			return 0, err
		}
		sum += converted
	}
	return sum, nil
}
