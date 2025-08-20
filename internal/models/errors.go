// Package models provides data structures and interfaces for genetic algorithm solutions.
package models

import (
	"errors"
	"fmt"
)

// InvalidChromosomeError represents an error that occurs when a chromosome contains invalid data type or content.
type InvalidChromosomeError struct {
	Message string
	Wrapped error
}

// Error implements the error interface
func (e *InvalidChromosomeError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
	}
	return e.Message
}

// Unwrap enables errors.Is and errors.As to traverse the error chain.
func (e *InvalidChromosomeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Wrapped
}

// NewInvalidChromosomeError creates a new InvalidChromosomeError with the given message and wrapped error.
func NewInvalidChromosomeError(message string, wrapped error) *InvalidChromosomeError {
	return &InvalidChromosomeError{
		Message: message,
		Wrapped: wrapped,
	}
}

// Sentinel errors for the models package
var (

	// ErrFitnessEvaluationFailed indicates that fitness evaluation failed.
	// This error occurs when the fitness evaluator cannot compute a fitness value.
	ErrFitnessEvaluationFailed = errors.New("fitness evaluation failed")

	// ErrPopulationEmpty indicates that a population has no individuals.
	// This error occurs when trying to perform operations on an empty population.
	ErrPopulationEmpty = errors.New("population is empty")
)
