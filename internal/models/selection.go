package models

import (
	"cmp"
	"fmt"
	"math/rand"
	"sort"
)

// SelectionError represents an error that occurs during a selection process.
type SelectionError struct {
	Message string
	Wrapped error
}

// Error implements the error interface.
func (e *SelectionError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
	}
	return e.Message
}

// Unwrap enables errors.Is and errors.As to traverse the error chain.
func (e *SelectionError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Wrapped
}

// NewSelectionError constructs a *SelectionError with the provided message and wrapped error.
func NewSelectionError(message string, wrapped error) *SelectionError {
	return &SelectionError{
		Message: message,
		Wrapped: wrapped,
	}
}

// ISelector defines the interface for selection operators in genetic algorithms.
type ISelector[T cmp.Ordered] interface {
	Select(population *Population[T]) (*Population[T], error)
}

// TournamentSelector performs selection using a tournament method.
// It includes support for elitism, where the best individuals from the
// current generation are carried over to the next.
type TournamentSelector[T cmp.Ordered] struct {
	TournamentSize int
	NumElites      int
}

// NewTournamentSelector creates a new TournamentSelector with the specified
// tournament size and number of elites.
func NewTournamentSelector[T cmp.Ordered](tournamentSize int, numElites int) (*TournamentSelector[T], error) {
	if tournamentSize <= 0 {
		return nil, NewSelectionError("invalid tournament size", fmt.Errorf("tournament size must be positive, but was %d", tournamentSize))
	}
	if numElites < 0 {
		return nil, NewSelectionError("invalid number of elites", fmt.Errorf("number of elites cannot be negative, but was %d", numElites))
	}
	return &TournamentSelector[T]{
		TournamentSize: tournamentSize,
		NumElites:      numElites,
	}, nil
}

// Select performs tournament selection on a population. It creates a new
// population of the same size, composed of individuals selected through
// a series of tournaments. If elitism is enabled, the fittest individuals
// are preserved and passed directly to the next generation.
func (ts *TournamentSelector[T]) Select(population *Population[T]) (*Population[T], error) {
	if population == nil || len(population.Individuals) == 0 {
		return nil, NewSelectionError("cannot perform selection on nil or empty population", ErrPopulationEmpty)
	}

	populationSize := len(population.Individuals)
	if ts.NumElites >= populationSize {
		return nil, NewSelectionError(
			fmt.Sprintf("number of elites (%d) is greater than or equal to population size (%d)", ts.NumElites, populationSize), nil)
	}

	offspring := make([]Solution[T], 0, populationSize)
	var selectionPool []Solution[T]

	if ts.NumElites > 0 {
		sortedIndividuals := make([]Solution[T], populationSize)
		copy(sortedIndividuals, population.Individuals)
		sort.Slice(sortedIndividuals, func(i, j int) bool {
			return sortedIndividuals[i].Fitness > sortedIndividuals[j].Fitness
		})

		for i := 0; i < ts.NumElites; i++ {
			offspring = append(offspring, *sortedIndividuals[i].DeepCopy())
		}
		selectionPool = sortedIndividuals[ts.NumElites:]
	} else {
		selectionPool = population.Individuals
	}

	numToSelect := populationSize - ts.NumElites
	selectionPoolSize := len(selectionPool)

	for i := 0; i < numToSelect; i++ {
		winnerIndex := rand.Intn(selectionPoolSize)
		for j := 1; j < ts.TournamentSize; j++ {
			competitorIndex := rand.Intn(selectionPoolSize)
			if selectionPool[competitorIndex].Fitness > selectionPool[winnerIndex].Fitness {
				winnerIndex = competitorIndex
			}
		}
		offspring = append(offspring, *selectionPool[winnerIndex].DeepCopy())
	}

	return &Population[T]{Individuals: offspring}, nil
}
