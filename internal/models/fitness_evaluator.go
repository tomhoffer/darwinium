// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import (
	"cmp"
	"context"
	"fmt"

	"github.com/tomhoffer/darwinium/internal/utils"
)

// IFitnessEvaluator defines the interface for fitness evaluation in genetic algorithms.
// Implementations must provide a method to evaluate the fitness of a chromosome.
type IFitnessEvaluator[T cmp.Ordered] interface {
	// Evaluate calculates the fitness value of a given chromosome.
	// Higher fitness values typically indicate better solutions.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - chromosome: The genetic material to evaluate
	//
	// Returns:
	//   - fitness: The calculated fitness value
	//   - error: Any error that occurred during evaluation
	Evaluate(ctx context.Context, chromosome *[]T) (float64, error)
}

// SimpleSumFitnessEvaluator implements a basic fitness evaluator that calculates
// fitness as the sum of all values in the chromosome. This is a simple example
// implementation that can be used for testing or as a baseline.
type SimpleSumFitnessEvaluator[T cmp.Ordered] struct{}

// NewSimpleSumFitnessEvaluator creates and returns a new SimpleSumFitnessEvaluator instance.
func NewSimpleSumFitnessEvaluator[T cmp.Ordered]() *SimpleSumFitnessEvaluator[T] {
	return &SimpleSumFitnessEvaluator[T]{}
}

// Evaluate calculates the fitness of a chromosome by summing all its values.
// Each value in the chromosome is converted to float64 before summing.
// This method implements the IFitnessEvaluator interface.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - chromosome: The genetic material to evaluate
//
// Returns:
//   - float64: The sum of all chromosome values as the fitness
//   - error: Any error that occurred during conversion or calculation
func (s SimpleSumFitnessEvaluator[T]) Evaluate(ctx context.Context, chromosome *[]T) (float64, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return 0.0, NewFitnessEvaluationError("context cancelled", ctx.Err())
	default:
	}

	if len(*chromosome) == 0 {
		return 0.0, NewFitnessEvaluationError("cannot calculate fitness", NewInvalidChromosomeError("empty chromosome found", nil))
	}

	var sum float64
	for _, val := range *chromosome {
		// Check for context cancellation in the loop
		select {
		case <-ctx.Done():
			return 0.0, NewFitnessEvaluationError("context cancelled", ctx.Err())
		default:
		}

		converted, err := utils.ConvertToFloat64(val)
		if err != nil {
			errMsg := fmt.Sprintf("invalid chromosome value %v", val)
			wrappedErr := NewInvalidChromosomeError("unable to convert chromosome value to float64", err)
			return 0, NewFitnessEvaluationError(errMsg, wrappedErr)
		}
		sum += converted
	}
	return sum, nil
}

// FitnessEvaluationError represents an error that occurs during a fitness computation or evaluation process.
// Message provides a summary of the error, while Wrapped contains the underlying cause, if present.
type FitnessEvaluationError struct {
	// Message describes the error at a high level.
	Message string
	// Wrapped holds the underlying error that triggered this error. Can be nil.
	Wrapped error
}

// Error implements the error interface.
func (e *FitnessEvaluationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
}

// Unwrap enables errors.Is and errors.As to traverse the error chain.
func (e *FitnessEvaluationError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Wrapped
}

// NewFitnessEvaluationError constructs a *FitnessEvaluationError with the provided
// message and wrapped error. Returning the concrete type makes the constructor's
// intent explicit while callers can still use it as an error (the value satisfies
// the error interface).
func NewFitnessEvaluationError(message string, wrapped error) *FitnessEvaluationError {
	return &FitnessEvaluationError{
		Message: message,
		Wrapped: wrapped,
	}
}
