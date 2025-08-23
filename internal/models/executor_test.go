package models

import (
	"cmp"
	"context"
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

func (m *mockFitnessEvaluator[T]) Evaluate(ctx context.Context, chromosome *[]T) (float64, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return 0.0, ctx.Err()
	default:
	}

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		floatCrossover := &mockCrossover[float64]{}

		floatExecutor := NewGeneticAlgorithmExecutor(floatPopulation, floatEvaluator, floatMutator, floatSelector, floatCrossover, 10)
		require.NotNil(t, floatExecutor)

		// Test with string
		stringPopulation := createTestPopulation([][]string{{"1", "2"}, {"3", "4"}})
		stringEvaluator := &mockFitnessEvaluator[string]{
			fitnessValues: []float64{100.0, 200.0},
			errorOnIndex:  -1,
		}
		stringMutator := &mockMutator[string]{}
		stringSelector := &mockSelector[string]{}
		stringCrossover := &mockCrossover[string]{}

		stringExecutor := NewGeneticAlgorithmExecutor(stringPopulation, stringEvaluator, stringMutator, stringSelector, stringCrossover, 10)
		require.NotNil(t, stringExecutor)
	})

	t.Run("creates executor with various numWorkers values", func(t *testing.T) {
		population := createTestPopulation([][]int{{1, 2}, {3, 4}})
		fitnessEvaluator := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{10.0, 20.0},
			errorOnIndex:  -1,
		}
		mutator := &mockMutator[int]{}
		selector := &mockSelector[int]{}
		crossover := &mockCrossover[int]{}

		// Test with default numWorkers
		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)
		require.NotNil(t, executor)
		assert.Equal(t, 1, executor.numWorkers, "Should default to 1 worker")

		// Test with explicit numWorkers
		executor = NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10, 4)
		require.NotNil(t, executor)
		assert.Equal(t, 4, executor.numWorkers, "Should set numWorkers to specified value")

		// Test with 0 workers (edge case)
		executorZero := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10, 0)
		require.NotNil(t, executorZero)
		assert.Equal(t, 0, executorZero.numWorkers, "Should set numWorkers to 0 when specified")

		// Test with negative workers (edge case)
		executorNeg := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10, -1)
		require.NotNil(t, executorNeg)
		assert.Equal(t, -1, executorNeg.numWorkers, "Should set numWorkers to -1 when specified")
	})
}

func TestGeneticAlgorithmExecutor_RefreshFitness(t *testing.T) {
	t.Run("successfully refreshes fitness for all individuals with default workers", func(t *testing.T) {
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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[float64]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[string]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[uint8]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		crossover := NewSinglePointCrossover[int]()

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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
		executor.population = newPopulation // Update executor's population

		// Check that the best solution was carried over by elitism
		assert.Contains(t, newPopulation.Individuals, bestSolution)

		// 4. Verify chromosome integrity after selection
		for i, individual := range newPopulation.Individuals {
			currentSum := 0
			for _, gene := range individual.Chromosome {
				currentSum += gene
			}
			assert.Contains(t, initialSums, currentSum, "Sum of genes for individual %d changed after selection and crossover", i)
		}

		// 4. Perform Crossover
		crossoverPopulation, err := executor.PerformCrossover()
		require.NoError(t, err)
		require.NotNil(t, crossoverPopulation)
		assert.Len(t, crossoverPopulation.Individuals, 4)

		// 6. Perform mutation
		err = executor.PerformMutation()
		require.NoError(t, err)
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
		crossover := &mockCrossover[CustomInt]{}

		executor := NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossover, 10)

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

type mockCrossover[T cmp.Ordered] struct {
	CrossoverFunc func(p1, p2 []T) ([]T, []T, error)
}

func (m *mockCrossover[T]) Crossover(p1, p2 []T) ([]T, []T, error) {
	if m.CrossoverFunc != nil {
		return m.CrossoverFunc(p1, p2)
	}
	return p1, p2, nil
}

func TestGeneticAlgorithmExecutor_PerformMutation(t *testing.T) {
	t.Run("mutates all individuals using mutator", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
		})

		executor := NewGeneticAlgorithmExecutor[int](population, nil, &mockMutator[int]{errorOnIndex: -1}, nil, &mockCrossover[int]{}, 10)
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

		executor := NewGeneticAlgorithmExecutor[int](population, nil, &mockMutator[int]{errorOnIndex: 1}, nil, &mockCrossover[int]{}, 10)
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
		executorNil := NewGeneticAlgorithmExecutor[int](nil, nil, &mockMutator[int]{}, nil, &mockCrossover[int]{}, 10)
		err := executorNil.PerformMutation()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// empty population
		popFactory := NewPopulationFactory[int]()
		emptyPop := popFactory.CreateEmptyPopulation()
		executorEmpty := NewGeneticAlgorithmExecutor[int](emptyPop, nil, &mockMutator[int]{}, nil, &mockCrossover[int]{}, 10)
		err = executorEmpty.PerformMutation()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// population with nil Individuals slice
		populationNilIndividuals := &Population[int]{}
		executorNilIndividuals := NewGeneticAlgorithmExecutor[int](populationNilIndividuals, nil, &mockMutator[int]{}, nil, &mockCrossover[int]{}, 10)
		err = executorNilIndividuals.PerformMutation()
		assert.ErrorIs(t, err, ErrPopulationEmpty)
	})
}

func TestGeneticAlgorithmExecutor_PerformSelection(t *testing.T) {
	t.Run("returns new population on successful selection", func(t *testing.T) {
		originalPopulation := createTestPopulation([][]int{{1, 2}, {3, 4}})
		expectedPopulation := createTestPopulation([][]int{{5, 6}, {7, 8}})

		ms := &mockSelector[int]{populationToReturn: expectedPopulation, errorOnIndex: -1}
		executor := NewGeneticAlgorithmExecutor[int](originalPopulation, nil, nil, ms, &mockCrossover[int]{}, 10)

		newPopulation, err := executor.PerformSelection()
		require.NoError(t, err)
		assert.Equal(t, expectedPopulation, newPopulation)
		assert.Equal(t, 1, ms.callCount)
	})

	t.Run("returns error when selector fails", func(t *testing.T) {
		originalPopulation := createTestPopulation([][]int{{1, 2}, {3, 4}})

		ms := &mockSelector[int]{populationToReturn: nil, errorOnIndex: 0}
		executor := NewGeneticAlgorithmExecutor[int](originalPopulation, nil, nil, ms, &mockCrossover[int]{}, 10)

		newPopulation, err := executor.PerformSelection()
		require.Error(t, err)
		assert.ErrorIs(t, err, errMockSelection)
		assert.Nil(t, newPopulation)
		assert.Equal(t, 1, ms.callCount)
	})

	t.Run("returns ErrPopulationEmpty for nil or empty population", func(t *testing.T) {
		// nil population
		executorNil := NewGeneticAlgorithmExecutor[int](nil, nil, nil, &mockSelector[int]{}, &mockCrossover[int]{}, 10)
		_, err := executorNil.PerformSelection()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// empty population
		popFactory := NewPopulationFactory[int]()
		emptyPop := popFactory.CreateEmptyPopulation()
		executorEmpty := NewGeneticAlgorithmExecutor[int](emptyPop, nil, nil, &mockSelector[int]{}, &mockCrossover[int]{}, 10)
		_, err = executorEmpty.PerformSelection()
		assert.ErrorIs(t, err, ErrPopulationEmpty)

		// population with nil Individuals slice
		populationNilIndividuals := &Population[int]{}
		executorNilIndividuals := NewGeneticAlgorithmExecutor[int](populationNilIndividuals, nil, nil, &mockSelector[int]{}, &mockCrossover[int]{}, 10)
		_, err = executorNilIndividuals.PerformSelection()
		assert.ErrorIs(t, err, ErrPopulationEmpty)
	})
}

func TestPerformCrossover(t *testing.T) {
	t.Run("successful crossover on even population", func(t *testing.T) {
		pop := createTestPopulation([][]int{
			{1, 2},
			{3, 4},
		})
		crossover := &mockCrossover[int]{
			CrossoverFunc: func(p1, p2 []int) ([]int, []int, error) {
				return []int{1}, []int{2}, nil
			},
		}
		executor := NewGeneticAlgorithmExecutor(pop, &mockFitnessEvaluator[int]{}, &mockMutator[int]{}, &mockSelector[int]{}, crossover, 10)
		offspring, err := executor.PerformCrossover()
		require.NoError(t, err)
		assert.Len(t, offspring.Individuals, 2)
	})

	t.Run("successful crossover on odd population", func(t *testing.T) {
		pop := createTestPopulation([][]int{
			{1, 2},
			{3, 4},
			{5, 6},
		})
		crossover := &mockCrossover[int]{
			CrossoverFunc: func(p1, p2 []int) ([]int, []int, error) {
				return []int{1}, []int{2}, nil
			},
		}
		executor := NewGeneticAlgorithmExecutor(pop, &mockFitnessEvaluator[int]{}, &mockMutator[int]{}, &mockSelector[int]{}, crossover, 10)
		offspring, err := executor.PerformCrossover()
		require.NoError(t, err)
		assert.Len(t, offspring.Individuals, 3)
	})

	t.Run("crossover failure results in correctly propagated error", func(t *testing.T) {
		pop := createTestPopulation([][]int{
			{1, 2},
			{3, 4},
		})
		crossover := &mockCrossover[int]{
			CrossoverFunc: func(p1, p2 []int) ([]int, []int, error) {
				return nil, nil, NewCrossoverError("mock crossover error", assert.AnError)
			},
		}
		executor := NewGeneticAlgorithmExecutor(pop, &mockFitnessEvaluator[int]{}, &mockMutator[int]{}, &mockSelector[int]{}, crossover, 10)
		_, err := executor.PerformCrossover()
		require.Error(t, err)

		var ce *CrossoverError
		assert.ErrorAs(t, err, &ce)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("empty population returns error", func(t *testing.T) {
		pop := &Population[int]{}
		crossover := &mockCrossover[int]{}
		executor := NewGeneticAlgorithmExecutor(pop, &mockFitnessEvaluator[int]{}, &mockMutator[int]{}, &mockSelector[int]{}, crossover, 10)
		_, err := executor.PerformCrossover()
		require.Error(t, err)
		var ce *CrossoverError
		assert.ErrorAs(t, err, &ce)
		assert.ErrorIs(t, err, ErrPopulationEmpty)
	})

	t.Run("nil population returns error", func(t *testing.T) {
		crossover := &mockCrossover[int]{}
		executor := NewGeneticAlgorithmExecutor(nil, &mockFitnessEvaluator[int]{}, &mockMutator[int]{}, &mockSelector[int]{}, crossover, 10)
		_, err := executor.PerformCrossover()
		require.Error(t, err)
		var ce *CrossoverError
		assert.ErrorAs(t, err, &ce)
		assert.ErrorIs(t, err, ErrPopulationEmpty)
	})
}

func TestGeneticAlgorithmExecutor_Loop(t *testing.T) {
	t.Run("executes all genetic algorithm steps correctly", func(t *testing.T) {
		// Create a population with known fitness values
		population := createTestPopulation([][]int{
			{1, 2, 3}, // Sum: 6
			{4, 5, 6}, // Sum: 15
			{7, 8, 9}, // Sum: 24
		})

		// Create mock components that track call order and count
		mockFitness := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{6.0, 15.0, 24.0, 8.0, 16.0, 25.0, 10.0, 20.0, 30.0}, // Values for 2 iterations (6) + final evaluation (3)
			errorOnIndex:  -1,
		}

		mockMutator := &mockMutator[int]{errorOnIndex: -1}

		// Mock selector that returns the same population to maintain size
		mockSelector := &mockSelector[int]{
			populationToReturn: population,
			errorOnIndex:       -1,
		}
		mockCrossover := &mockCrossover[int]{}

		// Create executor with mocked components
		executor := NewGeneticAlgorithmExecutor(population, mockFitness, mockMutator, mockSelector, mockCrossover, 2)

		// Execute the loop
		resultPopulation, err := executor.Loop(2)

		// Verify no errors occurred
		require.NoError(t, err)
		require.NotNil(t, resultPopulation)

		// Verify that fitness evaluation was called for each iteration
		// First iteration: 3 individuals (initial fitness)
		// Second iteration: 3 individuals (initial fitness)
		// Final call: 3 individuals (after all generations complete)
		// Total: 9 calls
		assert.Equal(t, 9, mockFitness.callCount, "Fitness evaluator should be called 9 times (3 individuals × 2 iterations + 3 individuals final evaluation)")

		// Verify that mutation was called for each iteration
		// First iteration: 3 individuals, Second iteration: 3 individuals
		assert.Equal(t, 6, mockMutator.callCount, "Mutator should be called 6 times (3 individuals × 2 iterations)")

		// Verify that selection was called for each iteration
		assert.Equal(t, 2, mockSelector.callCount, "Selector should be called 2 times (once per iteration)")

		// Verify that crossover was called for each iteration
		// This would need a more sophisticated mock to track, but we can verify the population structure
		assert.Len(t, executor.population.Individuals, 3, "Population should maintain its size through iterations")
	})

	t.Run("stops execution when fitness evaluation fails", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
		})

		// Mock fitness evaluator that fails on the first call
		mockFitness := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{6.0}, // Only one value before failure
			errorOnIndex:  0,              // Fail on 1st call
		}

		mockMutator := &mockMutator[int]{errorOnIndex: -1}

		// Mock selector that returns the same population to maintain size
		mockSelector := &mockSelector[int]{
			populationToReturn: population,
			errorOnIndex:       -1,
		}
		mockCrossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, mockFitness, mockMutator, mockSelector, mockCrossover, 2)

		// Execute the loop - should fail immediately during fitness evaluation
		resultPopulation, err := executor.Loop(2)

		// Verify error occurred
		require.Error(t, err)
		assert.Nil(t, resultPopulation)
		assert.ErrorIs(t, err, ErrFitnessEvaluationFailed)

		// Verify that fitness evaluation was called once before failing
		assert.Equal(t, 1, mockFitness.callCount, "Fitness evaluator should be called 1 time before failing")

		// Verify that no other components were called due to early failure
		assert.Equal(t, 0, mockMutator.callCount, "Mutator should not be called due to fitness failure")
		assert.Equal(t, 0, mockSelector.callCount, "Selector should not be called due to fitness failure")
	})

	t.Run("stops execution when mutation fails", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
		})

		mockFitness := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{6.0, 15.0},
			errorOnIndex:  -1,
		}

		// Mock mutator that fails on the first call
		mockMutator := &mockMutator[int]{errorOnIndex: 0} // Fail on 1st call

		// Mock selector that returns the same population to maintain size
		mockSelector := &mockSelector[int]{
			populationToReturn: population,
			errorOnIndex:       -1,
		}
		mockCrossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, mockFitness, mockMutator, mockSelector, mockCrossover, 2)

		// Execute the loop - should fail during first iteration
		resultPopulation, err := executor.Loop(2)

		// Verify error occurred
		require.Error(t, err)
		assert.ErrorIs(t, err, errMockMutation)
		assert.Nil(t, resultPopulation)

		// Verify that fitness evaluation was called for first iteration
		assert.Equal(t, 2, mockFitness.callCount, "Fitness evaluator should be called 2 times (first iteration)")

		// Verify that mutation was called once before failing
		assert.Equal(t, 1, mockMutator.callCount, "Mutator should be called 1 time before failing")

		// Verify that selection was called once before mutation failure
		assert.Equal(t, 1, mockSelector.callCount, "Selector should be called 1 time before mutation failure")
	})

	t.Run("stops execution when selection fails", func(t *testing.T) {
		population := createTestPopulation([][]int{
			{1, 2, 3},
			{4, 5, 6},
		})

		mockFitness := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{6.0, 15.0},
			errorOnIndex:  -1,
		}

		mockMutator := &mockMutator[int]{errorOnIndex: -1}

		// Mock selector that fails on the first call
		mockSelector := &mockSelector[int]{errorOnIndex: 0} // Fail on 1st call

		mockCrossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(population, mockFitness, mockMutator, mockSelector, mockCrossover, 2)

		// Execute the loop - should fail during first iteration
		resultPopulation, err := executor.Loop(2)

		// Verify error occurred
		require.Error(t, err)
		assert.ErrorIs(t, err, errMockSelection)
		assert.Nil(t, resultPopulation)

		// Verify that fitness evaluation was called for first iteration
		assert.Equal(t, 2, mockFitness.callCount, "Fitness evaluator should be called 2 times (first iteration)")

		// Verify that mutation was not called due to selection failure
		assert.Equal(t, 0, mockMutator.callCount, "Mutator should not be called due to selection failure")

		// Verify that selection was called once before failing
		assert.Equal(t, 1, mockSelector.callCount, "Selector should be called 1 time before failing")
	})

	t.Run("handles empty population", func(t *testing.T) {
		emptyPopulation := createTestPopulation([][]int{})

		mockFitness := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{},
			errorOnIndex:  -1,
		}
		mockMutator := &mockMutator[int]{errorOnIndex: -1}

		// Mock selector that returns the same population to maintain size
		mockSelector := &mockSelector[int]{
			populationToReturn: emptyPopulation,
			errorOnIndex:       -1,
		}
		mockCrossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(emptyPopulation, mockFitness, mockMutator, mockSelector, mockCrossover, 1)

		// Execute the loop - should fail immediately
		resultPopulation, err := executor.Loop(1)

		// Verify error occurred
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPopulationEmpty)
		assert.Nil(t, resultPopulation)

		// Verify that no components were called
		assert.Equal(t, 0, mockFitness.callCount, "Fitness evaluator should not be called")
		assert.Equal(t, 0, mockMutator.callCount, "Mutator should not be called")
		assert.Equal(t, 0, mockSelector.callCount, "Selector should not be called")
	})

	t.Run("handles nil population", func(t *testing.T) {
		mockFitness := &mockFitnessEvaluator[int]{
			fitnessValues: []float64{},
			errorOnIndex:  -1,
		}
		mockMutator := &mockMutator[int]{errorOnIndex: -1}

		// Mock selector that returns nil population (will cause error)
		mockSelector := &mockSelector[int]{
			populationToReturn: nil,
			errorOnIndex:       -1,
		}
		mockCrossover := &mockCrossover[int]{}

		executor := NewGeneticAlgorithmExecutor(nil, mockFitness, mockMutator, mockSelector, mockCrossover, 1)

		// Execute the loop - should fail immediately
		resultPopulation, err := executor.Loop(1)

		// Verify error occurred
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPopulationEmpty)
		assert.Nil(t, resultPopulation)

		// Verify that no components were called
		assert.Equal(t, 0, mockFitness.callCount, "Fitness evaluator should not be called")
		assert.Equal(t, 0, mockMutator.callCount, "Mutator should not be called")
		assert.Equal(t, 0, mockSelector.callCount, "Selector should not be called")
	})
}

// BenchmarkExecutor_Loop benchmarks the Loop method with real dependencies
func BenchmarkExecutor_Loop(b *testing.B) {
	populationSize := 1000
	chromosomeLength := 100
	generations := 10
	tournamentSize := 3

	// Use real genetic algorithm components
	fitnessEvaluator := NewSimpleSumFitnessEvaluator[int]()
	mutator := NewSimpleSwapMutator[int](0.01)                    // 1% mutation rate
	selector, _ := NewTournamentSelector[int](tournamentSize, 10) // 10 elitism
	crossover := NewSinglePointCrossover[int]()

	// Create executor with real components
	executor := NewGeneticAlgorithmExecutor(nil, fitnessEvaluator, mutator, selector, crossover, generations)

	// Pre-create all populations to avoid timing them
	populations := make([]*Population[int], b.N)
	for i := 0; i < b.N; i++ {
		populations[i] = createBenchmarkPopulation(populationSize, chromosomeLength)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Just assign the pre-created population (no creation, no timing)
		executor.population = populations[i]

		// Run the genetic algorithm loop (this IS timed)
		_, err := executor.Loop(generations)
		if err != nil {
			b.Fatalf("Loop failed: %v", err)
		}
	}
}
