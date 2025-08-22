// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import (
	"fmt"
	"math/rand"
)

// ICrossover defines the interface for chromosome crossover in genetic algorithms.
type ICrossover[T any] interface {
	// Crossover performs crossover on two parent chromosomes to produce offspring.
	//
	// Parameters:
	//   - parent1: The first parent chromosome.
	//   - parent2: The second parent chromosome.
	//
	// Returns:
	//   - []T: The first offspring chromosome.
	//   - []T: The second offspring chromosome.
	//   - error: Any error that occurred during crossover.
	Crossover(parent1, parent2 []T) ([]T, []T, error)
}

// SinglePointCrossover implements the single-point crossover method.
// In single-point crossover, a point on both parent chromosome strings is
// picked randomly and designated a 'crossover point'. Genes to the right
// of that point are swapped between the two parent chromosomes.
// This results in two offspring, each carrying some genetic material
// from both parents.
type SinglePointCrossover[T any] struct{}

// NewSinglePointCrossover creates and returns a new SinglePointCrossover instance.
func NewSinglePointCrossover[T any]() *SinglePointCrossover[T] {
	return &SinglePointCrossover[T]{}
}

// Crossover performs a single-point crossover on two parent chromosomes.
func (s SinglePointCrossover[T]) Crossover(parent1, parent2 []T) ([]T, []T, error) {
	if parent1 == nil || len(parent1) == 0 || parent2 == nil || len(parent2) == 0 {
		return nil, nil, NewCrossoverError("cannot perform crossover", NewInvalidChromosomeError("parent chromosomes cannot be empty", nil))
	}

	if len(parent1) != len(parent2) {
		return nil, nil, NewCrossoverError("cannot perform crossover", NewInvalidChromosomeError("parent chromosomes must be of the same length", nil))
	}

	parent1Len := len(parent1)
	if parent1Len == 1 {
		offspring1 := make([]T, parent1Len)
		offspring2 := make([]T, parent1Len)
		copy(offspring1, parent1)
		copy(offspring2, parent2)
		return offspring1, offspring2, nil
	}

	// Crossover_point is between 1 and parent1Len-1 inclusive.
	crossoverPoint := rand.Intn(parent1Len-1) + 1

	offspring1 := make([]T, parent1Len)
	copy(offspring1[:crossoverPoint], parent1[:crossoverPoint])
	copy(offspring1[crossoverPoint:], parent2[crossoverPoint:])

	offspring2 := make([]T, parent1Len)
	copy(offspring2[:crossoverPoint], parent2[:crossoverPoint])
	copy(offspring2[crossoverPoint:], parent1[crossoverPoint:])

	return offspring1, offspring2, nil
}

// CrossoverError represents an error that occurs during a crossover process.
// Message provides a summary of the error, while Wrapped contains the underlying cause, if present.
type CrossoverError struct {
	// Message describes the error at a high level.
	Message string
	// Wrapped holds the underlying error that triggered this error. Can be nil.
	Wrapped error
}

// Error implements the error interface.
func (e *CrossoverError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
}

// Unwrap enables errors.Is and errors.As to traverse the error chain.
func (e *CrossoverError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Wrapped
}

// NewCrossoverError constructs a *CrossoverError with the provided message and wrapped error.
func NewCrossoverError(message string, wrapped error) *CrossoverError {
	return &CrossoverError{
		Message: message,
		Wrapped: wrapped,
	}
}
