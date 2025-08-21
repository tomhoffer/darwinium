// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import (
	"cmp"
	"fmt"
	"math/rand"
)

// IMutator defines the interface for chromosome mutation in genetic algorithms.
// Implementations should modify the provided chromosome in place.
type IMutator[T cmp.Ordered] interface {
	// Mutate applies a mutation to the given chromosome in place.
	//
	// Parameters:
	//   - chromosome: The genetic material to mutate
	//
	// Returns:
	//   - error: Any error that occurred during mutation
	Mutate(chromosome *[]T) error
}

// SimpleSwapMutator performs a simple mutation by swapping two distinct genes
// at randomly selected positions. This operation is valid for any ordered type.
type SimpleSwapMutator[T cmp.Ordered] struct{}

// NewSimpleSwapMutator creates and returns a new SimpleSwapMutator instance.
func NewSimpleSwapMutator[T cmp.Ordered]() *SimpleSwapMutator[T] {
	return &SimpleSwapMutator[T]{}
}

// Mutate swaps two distinct positions in the chromosome. The mutation is performed in place.
// Returns a MutationError wrapping an InvalidChromosomeError if the chromosome is empty or too short.
func (s SimpleSwapMutator[T]) Mutate(chromosome *[]T) error {
	// Validate chromosome
	if chromosome == nil || *chromosome == nil || len(*chromosome) == 0 {
		return NewMutationError("cannot mutate chromosome", NewInvalidChromosomeError("empty chromosome found", nil))
	}
	if len(*chromosome) < 2 {
		return NewMutationError("cannot mutate chromosome", NewInvalidChromosomeError("chromosome must contain at least 2 genes", nil))
	}

	n := len(*chromosome)
	firstPosition := rand.Intn(n)
	secondPosition := rand.Intn(n - 1)

	if firstPosition == secondPosition {
		secondPosition++
	}

	// Perform swap in place
	ch := *chromosome
	ch[firstPosition], ch[secondPosition] = ch[secondPosition], ch[firstPosition]
	return nil
}

// MutationError represents an error that occurs during a mutation process.
// Message provides a summary of the error, while Wrapped contains the underlying cause, if present.
type MutationError struct {
	// Message describes the error at a high level.
	Message string
	// Wrapped holds the underlying error that triggered this error. Can be nil.
	Wrapped error
}

// Error implements the error interface.
func (e *MutationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
}

// Unwrap enables errors.Is and errors.As to traverse the error chain.
func (e *MutationError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Wrapped
}

// NewMutationError constructs a *MutationError with the provided message and wrapped error.
func NewMutationError(message string, wrapped error) *MutationError {
	return &MutationError{
		Message: message,
		Wrapped: wrapped,
	}
}
