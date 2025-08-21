package models

import (
	"cmp"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock fitness evaluator for testing
type mockFitnessEvaluator[T cmp.Ordered] struct {
	fitnessValues []float64
	errorOnIndex  int // -1 means no error, otherwise error on this index
	callCount     int
}

func (m *mockFitnessEvaluator[T]) Evaluate(chromosome *[]T) (float64, error) {
	if m.errorOnIndex >= 0 && m.callCount == m.errorOnIndex {
		m.callCount++
		return 0, ErrFitnessEvaluationFailed
	}

	if m.callCount >= len(m.fitnessValues) {
		m.callCount++
		return 0, nil
	}

	fitness := m.fitnessValues[m.callCount]
	m.callCount++
	return fitness, nil
}

func (m *mockFitnessEvaluator[T]) reset() {
	m.callCount = 0
}

// Helper to create test population using factories
func createTestPopulation[T cmp.Ordered](chromosomes [][]T) *Population[T] {
	solutionFactory := NewSolutionFactory[T]()
	populationFactory := NewPopulationFactory[T]()

	individuals := make([]Solution[T], len(chromosomes))
	for i, chromosome := range chromosomes {
		solution := solutionFactory.CreateSolution(chromosome)
		individuals[i] = *solution
	}

	return populationFactory.CreatePopulation(individuals)
}

func TestNewGeneticAlgorithmExecutor(t *testing.T) {
	t.Run("creates executor with valid parameters", func(t *testing.T) {
		population := createTestPopulation([][]int{{1, 2}, {3, 4}})
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{10.0, 20.0},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		require.NotNil(t, executor)
		assert.Equal(t, population, executor.population)
		assert.Equal(t, fitnessEvaluator, executor.fitnessEvaluator)
	})

	t.Run("creates executor with empty population", func(t *testing.T) {
		populationFactory := NewPopulationFactory[int]()
		population := populationFactory.CreateEmptyPopulation()
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		require.NotNil(t, executor)
		assert.Len(t, executor.population.Individuals, 0)
	})

	t.Run("creates executor with different types", func(t *testing.T) {
		// Test with float64
		floatPopulation := createTestPopulation([][]float64{{1.5, 2.5}, {3.5, 4.5}})
		floatEvaluator := &mockFitnessEvaluator[float64]{
			fitnessValues: []float64{15.0, 25.0},
			errorOnIndex:  -1,
		}
		floatMutator := &mockMutator[float64]{}
		floatSelector := &mockSelector[float64]{}

		floatExecutor := NewGeneticAlgorithmExecutor(floatPopulation, floatEvaluator, floatMutator, floatSelector)
		require.NotNil(t, floatExecutor)

		// Test with string
		stringPopulation := createTestPopulation([][]string{{"1", "2"}, {"3", "4"}})
		stringEvaluator := &mockFitnessEvaluator[string]{
			fitnessValues: []float64{100.0, 200.0},
			errorOnIndex:  -1,
		}
		stringMutator := &mockMutator[string]{}
		stringSelector := &mockSelector[string]{}

		stringExecutor := NewGeneticAlgorithmExecutor(stringPopulation, stringEvaluator, stringMutator, stringSelector)
		require.NotNil(t, stringExecutor)
	})
}

func TestGeneticAlgorithmExecutor_RefreshFitness(t *testing.T) {
	t.Run("successfully refreshes fitness for all individuals", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		})

		expectedFitness := []float64{10.5, 20.5, 30.5}
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: expectedFitness,
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		// Verify initial fitness values are 0
		for i, individual := range executor.population.Individuals {
			assert.Equal(t, 0.0, individual.Fitness, "Individual %d should start with 0 fitness", i)
		}

		err := executor.RefreshFitness()

		require.NoError(t, err)
		assert.Equal(t, len(expectedFitness), fitnessEvaluator.callCount)

		// Verify fitness values were updated
		for i, individual := range executor.population.Individuals {
			assert.Equal(t, expectedFitness[i], individual.Fitness, "Individual %d should have correct fitness", i)
		}
	})

	t.Run("returns error when fitness evaluator fails", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		})

		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{10.5, 20.5, 30.5},
			errorOnIndex:  1, // Error on second individual
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		assert.ErrorIs(t, err, ErrFitnessEvaluationFailed)

		// Should have called evaluator twice (first success, second failure)
		assert.Equal(t, 2, fitnessEvaluator.callCount)

		// First individual should have updated fitness, others should remain 0
		assert.Equal(t, 10.5, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 0.0, executor.population.Individuals[1].Fitness)
		assert.Equal(t, 0.0, executor.population.Individuals[2].Fitness)
	})

	t.Run("returns error when fitness evaluator fails on first individual", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
		})

		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{10.5, 20.5},
			errorOnIndex:  0, // Error on first individual
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		assert.ErrorIs(t, err, ErrFitnessEvaluationFailed)
		assert.Equal(t, 1, fitnessEvaluator.callCount)

		// All individuals should still have 0 fitness
		for i, individual := range executor.population.Individuals {
			assert.Equal(t, 0.0, individual.Fitness, "Individual %d should remain with 0 fitness", i)
		}
	})

	t.Run("containing nil population returns error", func(t *testing.T) {
		var population *Population[int] = nil
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()
		assert.ErrorIs(t, err, ErrPopulationEmpty)
		assert.Equal(t, 0, fitnessEvaluator.callCount)
	})

	t.Run("handles empty population", func(t *testing.T) {
		population := createTestPopulation([][]int{})
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()
		assert.ErrorIs(t, err, ErrPopulationEmpty)
		assert.Equal(t, 0, fitnessEvaluator.callCount)
	})

	t.Run("handles single individual population", func(t *testing.T) {
		population := createTestPopulation([][]int{{1, 2, 3}})
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{42.0},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.NoError(t, err)
		assert.Equal(t, 1, fitnessEvaluator.callCount)
		assert.Equal(t, 42.0, executor.population.Individuals[0].Fitness)
	})

	t.Run("can refresh fitness multiple times", func(t *testing.T) {
		population := createTestPopulation([][]int{{1, 2}, {3, 4}})
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{10.0, 20.0, 15.0, 25.0}, // Values for two refresh cycles
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		// First refresh
		err := executor.RefreshFitness()
		require.NoError(t, err)
		assert.Equal(t, 10.0, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 20.0, executor.population.Individuals[1].Fitness)

		// Second refresh - should update with new values
		err = executor.RefreshFitness()
		require.NoError(t, err)
		assert.Equal(t, 15.0, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 25.0, executor.population.Individuals[1].Fitness)
		assert.Equal(t, 4, fitnessEvaluator.callCount)
	})
}

func TestGeneticAlgorithmExecutor_RefreshFitness_WithDifferentTypes(t *testing.T) {
	t.Run("works with float64 chromosomes", func(t *testing.T) {
		population := createTestPopulation([][]float64{
			{1.1, 2.2, 3.3},
			{4.4, 5.5, 6.6},
		})

		fitnessEvaluator := &mockFitnessEvaluator[float64]{
			fitnessValues: []float64{11.11, 22.22},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[float64]{}
		selector := &mockSelector[float64]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.NoError(t, err)
		assert.Equal(t, 11.11, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 22.22, executor.population.Individuals[1].Fitness)
	})

	t.Run("works with string chromosomes", func(t *testing.T) {
		population := createTestPopulation([][]string{
			{"1", "2", "3"},
			{"4", "5", "6"},
		})

		fitnessEvaluator := &mockFitnessEvaluator[string]{
			fitnessValues: []float64{100.0, 200.0},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[string]{}
		selector := &mockSelector[string]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.NoError(t, err)
		assert.Equal(t, 100.0, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 200.0, executor.population.Individuals[1].Fitness)
	})

	t.Run("works with uint8 chromosomes", func(t *testing.T) {
		population := createTestPopulation([][]uint8{
			{1, 2, 3},
			{4, 5, 6},
		})

		fitnessEvaluator := &mockFitnessEvaluator[uint8]{
			fitnessValues: []float64{6.0, 15.0},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[uint8]{}
		selector := &mockSelector[uint8]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.NoError(t, err)
		assert.Equal(t, 6.0, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 15.0, executor.population.Individuals[1].Fitness)
	})
}

func TestGeneticAlgorithmExecutor_RefreshFitness_EdgeCases(t *testing.T) {
	t.Run("handles an individual with nil chromosome", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3}, // Normal chromosome
			nil,       // Nil chromosome
		})

		fitnessEvaluator := NewSimpleSumFitnessEvaluator[int]()
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.Error(t, err)

		var fe *FitnessEvaluationError
		require.True(t, errors.As(err, &fe), "expected FitnessEvaluationError")

		assert.Equal(t, 6.0, executor.population.Individuals[0].Fitness)
		// The second individual should not have been processed due to the error
		assert.Equal(t, 0.0, executor.population.Individuals[1].Fitness)
	})

	t.Run("handles an individual with empty chromosome", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3}, // Normal chromosome
			{},        // Empty chromosome
		})

		fitnessEvaluator := NewSimpleSumFitnessEvaluator[int]()
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.Error(t, err)

		var fe *FitnessEvaluationError
		require.True(t, errors.As(err, &fe), "expected FitnessEvaluationError")

		assert.Equal(t, 6.0, executor.population.Individuals[0].Fitness)
		// The second individual should not have been processed due to the error
		assert.Equal(t, 0.0, executor.population.Individuals[1].Fitness)
	})
}

// Integration test with real SimpleSumFitnessEvaluator
func TestGeneticAlgorithmExecutor_Integration(t *testing.T) {
	t.Run("integration with real components", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},   // Sum: 6
			{4, 5, 6},   // Sum: 15
			{-1, 0, 1},  // Sum: 0
			{10, -5, 2}, // Sum: 7
		})
		initialSums := []int{6, 15, 0, 7}

		// 1. Setup
		fitnessEvaluator := NewSimpleSumFitnessEvaluator[int]()
		mutator := NewSimpleSwapMutator[int]()
		selector, err := NewTournamentSelector[int](3, 1) // elitism of 1
		require.NoError(t, err)

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		// 2. Refresh fitness
		err = executor.RefreshFitness()
		require.NoError(t, err)
		assert.Equal(t, 6.0, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 15.0, executor.population.Individuals[1].Fitness)
		assert.Equal(t, 0.0, executor.population.Individuals[2].Fitness)
		assert.Equal(t, 7.0, executor.population.Individuals[3].Fitness)

		bestSolution := executor.population.Individuals[1] // fitness 15

		// 3. Perform selection
		newPopulation, err := executor.PerformSelection()
		require.NoError(t, err)
		require.NotNil(t, newPopulation)
		assert.Len(t, newPopulation.Individuals, 4)

		// Check that the best solution was carried over by elitism
		assert.Contains(t, newPopulation.Individuals, bestSolution)

		// 4. Perform mutation
		executor.population = newPopulation // Update executor's population
		err = executor.PerformMutation()
		require.NoError(t, err)

		// 5. Verify mutation integrity
		for i, individual := range executor.population.Individuals {
			currentSum := 0
			for _, gene := range individual.Chromosome {
				currentSum += gene
			}
			// SimpleSwapMutator should not change the sum of genes
			assert.Contains(t, initialSums, currentSum, "Sum of genes for individual %d changed after mutation", i)
		}
	})
}

func TestGeneticAlgorithmExecutor_CustomTypes(t *testing.T) {
	t.Run("works with custom ordered types", func(t *testing.T) {
		population := createTestPopulation([][]CustomInt{
			{CustomInt(1), CustomInt(2), CustomInt(3)},
			{CustomInt(4), CustomInt(5), CustomInt(6)},
		})

		fitnessEvaluator := &mockFitnessEvaluator[CustomInt]{
			fitnessValues: []float64{60.0, 150.0},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[CustomInt]{}
		selector := &mockSelector[CustomInt]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector)

		err := executor.RefreshFitness()

		require.NoError(t, err)
		assert.Equal(t, 60.0, executor.population.Individuals[0].Fitness)
		assert.Equal(t, 150.0, executor.population.Individuals[1].Fitness)
	})
}

// Mock mutator for testing PerformMutation
type mockMutator[T cmp.Ordered] struct {
	// errorOnIndex: -1 means no error; otherwise, return error on this call index
	errorOnIndex int
	callCount    int
}

var errMockMutation = errors.New("mock mutation error")

func (m *mockMutator[T]) Mutate(chromosome *[]T) error {
	if m.errorOnIndex >= 0 && m.callCount == m.errorOnIndex {
		m.callCount++
		return errMockMutation
	}
	// simple deterministic mutation: swap first two genes if possible
	if chromosome != nil && *chromosome != nil && len(*chromosome) >= 2 {
		ch := *chromosome
		ch[0], ch[1] = ch[1], ch[0]
	}
	m.callCount++
	return nil
}

// Mock selector for testing PerformSelection
type mockSelector[T cmp.Ordered] struct {
	// populationToReturn: the population to be returned on successful selection
	populationToReturn *Population[T]
	// errorOnIndex: -1 means no error; otherwise, return error on this call index
	errorOnIndex int
	callCount    int
}

var errMockSelection = errors.New("mock selection error")

func (m *mockSelector[T]) Select(population *Population[T]) (*Population[T], error) {
	if m.errorOnIndex >= 0 && m.callCount == m.errorOnIndex {
		m.callCount++
		return nil, errMockSelection
	}
	m.callCount++
	return m.populationToReturn, nil
}

func TestGeneticAlgorithmExecutor_PerformMutation(t *testing.T) {
	t.Run("mutates all individuals using mutator", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
		})

		executor := NewGeneticAlgorithmExecutor[int](population, nil, &mockMutator[int]{errorOnIndex: -1}, nil)
		mm := &mockMutator[int]{errorOnIndex: -1}
		executor.mutator = mm

		err := executor.PerformMutation()
		require.NoError(t, err)
		assert.Equal(t, 2, mm.callCount)
		assert.Equal(t, []int{2, 1, 3}, executor.population.Individuals[0].Chromosome)
		assert.Equal(t, []int{5, 4, 6}, executor.population.Individuals[1].Chromosome)
	})

	t.Run("returns error when mutator fails and stops further mutations", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		})

		executor := NewGeneticAlgorithmExecutor[int](population, nil, &mockMutator[int]{errorOnIndex: 1}, nil)
		mm := &mockMutator[int]{errorOnIndex: 1} // fail on second individual
		executor.mutator = mm

		err := executor.PerformMutation()
		require.Error(t, err)
		assert.ErrorIs(t, err, errMockMutation)
		// First individual should be mutated
		assert.Equal(t, []int{2, 1, 3}, executor.population.Individuals[0].Chromosome)
		// Second individual should remain unchanged due to error
		assert.Equal(t, []int{4, 5, 6}, executor.population.Individuals[1].Chromosome)
		// No further mutations should occur
		assert.Equal(t, 2, mm.callCount)
	})

	t.Run("returns ErrPopulationEmpty for nil or empty population", func(t *testing.T) {
		// nil population
		executorNil := NewGeneticAlgorithmExecutor[int](nil, nil, &mockMutator[int]{}, nil)
		err := executorNil.PerformMutation()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// empty population
		popFactory := NewPopulationFactory[int]()
		emptyPop := popFactory.CreateEmptyPopulation()
		executorEmpty := NewGeneticAlgorithmExecutor[int](emptyPop, nil, &mockMutator[int]{}, nil)
		err = executorEmpty.PerformMutation()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// population with nil Individuals slice
		populationNilIndividuals := &Population[int]{}
		executorNilIndividuals := NewGeneticAlgorithmExecutor[int](populationNilIndividuals, nil, &mockMutator[int]{}, nil)
		err = executorNilIndividuals.PerformMutation()
		assert.ErrorIs(t, err, ErrPopulationEmpty)
	})
}

func TestGeneticAlgorithmExecutor_PerformSelection(t *testing.T) {
	t.Run("returns new population on successful selection", func(t *testing.T) {
		originalPopulation := createTestPopulation([][]int{{1, 2}, {3, 4}})
		expectedPopulation := createTestPopulation([][]int{{5, 6}, {7, 8}})

		ms := &mockSelector[int]{populationToReturn: expectedPopulation, errorOnIndex: -1}
		executor := NewGeneticAlgorithmExecutor[int](originalPopulation, nil, nil, ms)

		newPopulation, err := executor.PerformSelection()
		require.NoError(t, err)
		assert.Equal(t, expectedPopulation, newPopulation)
		assert.Equal(t, 1, ms.callCount)
	})

	t.Run("returns error when selector fails", func(t *testing.T) {
		originalPopulation := createTestPopulation([][]int{{1, 2}, {3, 4}})

		ms := &mockSelector[int]{populationToReturn: nil, errorOnIndex: 0}
		executor := NewGeneticAlgorithmExecutor[int](originalPopulation, nil, nil, ms)

		newPopulation, err := executor.PerformSelection()
		require.Error(t, err)
		assert.ErrorIs(t, err, errMockSelection)
		assert.Nil(t, newPopulation)
		assert.Equal(t, 1, ms.callCount)
	})

	t.Run("returns ErrPopulationEmpty for nil or empty population", func(t *testing.T) {
		// nil population
		executorNil := NewGeneticAlgorithmExecutor[int](nil, nil, nil, &mockSelector[int]{})
		_, err := executorNil.PerformSelection()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// empty population
		popFactory := NewPopulationFactory[int]()
		emptyPop := popFactory.CreateEmptyPopulation()
		executorEmpty := NewGeneticAlgorithmExecutor[int](emptyPop, nil, nil, &mockSelector[int]{})
		_, err = executorEmpty.PerformSelection()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// population with nil Individuals slice
		populationNilIndividuals := &Population[int]{}
		executorNilIndividuals := NewGeneticAlgorithmExecutor[int](populationNilIndividuals, nil, nil, &mockSelector[int]{})
		_, err = executorNilIndividuals.PerformSelection()
		assert.ErrorIs(t, err, ErrPopulationEmpty)
	})
}
