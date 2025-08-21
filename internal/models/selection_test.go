// Package models provides tests for the models package.
package models

import (
	"cmp"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// Test suite for TournamentSelector
//

// TestNewTournamentSelector tests that the constructor validates its input arguments correctly.
func TestNewTournamentSelector(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		tournamentSize int
		numElites      int
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "Valid parameters",
			tournamentSize: 3,
			numElites:      2,
			expectError:    false,
		},
		{
			name:           "Zero tournament size",
			tournamentSize: 0,
			numElites:      2,
			expectError:    true,
		},
		{
			name:           "Negative tournament size",
			tournamentSize: -1,
			numElites:      2,
			expectError:    true,
		},
		{
			name:           "Negative number of elites",
			tournamentSize: 3,
			numElites:      -1,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			selector, err := NewTournamentSelector[int](tc.tournamentSize, tc.numElites)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, selector)
				var se *SelectionError
				assert.ErrorAs(t, err, &se)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, selector)
			}
		})
	}

	t.Run("Returns error when number of elites is greater than population size", func(t *testing.T) {
		t.Parallel()
		population := &Population[int]{
			Individuals: []Solution[int]{
				{Chromosome: []int{1}, Fitness: 10},
			},
		}
		selector, err := NewTournamentSelector[int](3, 2)
		require.NoError(t, err)

		_, err = selector.Select(population)
		var se *SelectionError
		assert.ErrorAs(t, err, &se)
	})
}

// TestTournamentSelector_Select tests the Select method of TournamentSelector across different scenarios.
func TestTournamentSelector_Select(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		population  *Population[int]
		selector    *TournamentSelector[int]
		assertError func(t *testing.T, err error)
		checkFunc   func(t *testing.T, original, selected *Population[int])
	}{
		{
			name:       "Empty population returns error",
			population: &Population[int]{Individuals: []Solution[int]{}},
			selector:   newSelector[int](t, 3, 0),
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrPopulationEmpty)
			},
		},
		{
			name:       "Nil population returns error",
			population: nil,
			selector:   newSelector[int](t, 3, 0),
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrPopulationEmpty)
			},
		},
		{
			name: "All individuals having zero fitness pass",
			population: &Population[int]{
				Individuals: []Solution[int]{
					{Chromosome: []int{1}, Fitness: 0},
					{Chromosome: []int{2}, Fitness: 0},
					{Chromosome: []int{3}, Fitness: 0},
				},
			},
			selector: newSelector[int](t, 3, 0),
			checkFunc: func(t *testing.T, original, selected *Population[int]) {
				assert.Equal(t, len(original.Individuals), len(selected.Individuals))
			},
		},
		{
			name: "All individuals being identical pass",
			population: &Population[int]{
				Individuals: []Solution[int]{
					{Chromosome: []int{1}, Fitness: 10},
					{Chromosome: []int{1}, Fitness: 10},
					{Chromosome: []int{1}, Fitness: 10},
				},
			},
			selector: newSelector[int](t, 3, 0),
			checkFunc: func(t *testing.T, original, selected *Population[int]) {
				assert.Equal(t, len(original.Individuals), len(selected.Individuals))
				for _, ind := range selected.Individuals {
					assert.Equal(t, original.Individuals[0], ind)
				}
			},
		},
		{
			name: "Various individuals pass",
			population: &Population[int]{
				Individuals: []Solution[int]{
					{Chromosome: []int{1}, Fitness: 1},
					{Chromosome: []int{1}, Fitness: 2},
					{Chromosome: []int{1}, Fitness: 3},
				},
			},
			selector: newSelector[int](t, 3, 0),
			checkFunc: func(t *testing.T, original, selected *Population[int]) {
				assert.Equal(t, len(original.Individuals), len(selected.Individuals))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			selectedPopulation, err := tc.selector.Select(tc.population)
			if tc.assertError != nil {
				assert.Error(t, err)
				tc.assertError(t, err)
				assert.Nil(t, selectedPopulation)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, selectedPopulation)
				tc.checkFunc(t, tc.population, selectedPopulation)
			}
		})
	}
}

// TestTournamentSelector_Select_WithElitism tests that elitism correctly preserves the best individuals.
func TestTournamentSelector_Select_WithElitism(t *testing.T) {
	t.Parallel()
	population := &Population[int]{
		Individuals: []Solution[int]{
			{Chromosome: []int{1}, Fitness: 10},
			{Chromosome: []int{2}, Fitness: 20},
			{Chromosome: []int{3}, Fitness: 30},
			{Chromosome: []int{4}, Fitness: 40},
			{Chromosome: []int{5}, Fitness: 50},
		},
	}
	selector := newSelector[int](t, 3, 2)

	selectedPopulation, err := selector.Select(population)

	require.NoError(t, err)
	require.NotNil(t, selectedPopulation)
	assert.Equal(t, len(population.Individuals), len(selectedPopulation.Individuals))

	// Elites are the two best individuals
	elite1 := Solution[int]{Chromosome: []int{5}, Fitness: 50}
	elite2 := Solution[int]{Chromosome: []int{4}, Fitness: 40}
	assert.Contains(t, selectedPopulation.Individuals, elite1)
	assert.Contains(t, selectedPopulation.Individuals, elite2)
}

// newSelector is a helper function to create a TournamentSelector, failing the test on error.
func newSelector[T cmp.Ordered](t *testing.T, tournamentSize, numElites int) *TournamentSelector[T] {
	t.Helper()
	selector, err := NewTournamentSelector[T](tournamentSize, numElites)
	require.NoError(t, err)
	return selector
}

// createBenchmarkPopulation is a helper function to create a population for benchmarking.
func createBenchmarkPopulation(size, chromosomeLength int) *Population[int] {
	individuals := make([]Solution[int], size)
	for i := 0; i < size; i++ {
		chromosome := make([]int, chromosomeLength)
		// math/rand is sufficient for benchmark data
		for j := 0; j < chromosomeLength; j++ {
			chromosome[j] = rand.Intn(100)
		}
		individuals[i] = Solution[int]{
			Chromosome: chromosome,
			Fitness:    rand.Float64() * 100,
		}
	}
	return &Population[int]{Individuals: individuals}
}

// BenchmarkTournamentSelector_Select benchmarks the Select method over various population sizes.
func BenchmarkTournamentSelector_Select(b *testing.B) {
	populationSize := 1000
	chromosomeLength := 50
	tournamentSize := 3

	population := createBenchmarkPopulation(populationSize, chromosomeLength)
	selector, err := NewTournamentSelector[int](tournamentSize, 0)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = selector.Select(population)
	}
}
