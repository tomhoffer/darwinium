package models

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for testing
func mockFitnessFunction(ch []int) (float64, error) {
	return 999.0, nil
}

func mockErrorFitnessFunction(ch []int) (float64, error) {
	return 0.0, errors.New("fitness calculation error")
}

func mockRandomIntGenerator() int {
	return rand.Intn(100)
}

func TestPopulation_RefreshFitness(t *testing.T) {
	solutionFactory := NewSolutionFactory[int]()
	populationFactory := NewPopulationFactory[int]()

	t.Run("fitness is refreshed for all individuals in population", func(t *testing.T) {
		// Create solutions with mock fitness function
		solution1 := solutionFactory.CreateSolution([]int{1, 2, 3}, mockFitnessFunction)
		solution2 := solutionFactory.CreateSolution([]int{4, 5, 6}, mockFitnessFunction)

		population := populationFactory.CreatePopulation([]Solution[int]{*solution1, *solution2})

		// Verify initial fitness is 0
		assert.Equal(t, 0.0, population.Individuals[0].Fitness)
		assert.Equal(t, 0.0, population.Individuals[1].Fitness)

		err := population.RefreshFitness()

		assert.NoError(t, err)
		assert.Equal(t, 999.0, population.Individuals[0].Fitness)
		assert.Equal(t, 999.0, population.Individuals[1].Fitness)
	})

	t.Run("fitness calculation error is returned", func(t *testing.T) {
		// Create solutions with error-producing fitness function
		solution1 := solutionFactory.CreateSolution([]int{1, 2, 3}, mockErrorFitnessFunction)
		solution2 := solutionFactory.CreateSolution([]int{4, 5, 6}, mockFitnessFunction)

		population := populationFactory.CreatePopulation([]Solution[int]{*solution1, *solution2})

		err := population.RefreshFitness()

		assert.Error(t, err)
		assert.Equal(t, "fitness calculation error", err.Error())
	})

	t.Run("refreshing fitness over an empty population succeeds", func(t *testing.T) {
		population := populationFactory.CreateEmptyPopulation()
		err := population.RefreshFitness()
		assert.NoError(t, err)
	})
}

func TestNewPopulationFactory(t *testing.T) {
	factory := NewPopulationFactory[int]()

	assert.NotNil(t, factory)
	assert.IsType(t, &PopulationFactory[int]{}, factory)
}

func TestPopulationFactory_CreatePopulation(t *testing.T) {
	populationFactory := NewPopulationFactory[int]()
	solutionFactory := NewSolutionFactory[int]()

	t.Run("create population with individuals is correct", func(t *testing.T) {
		solution1 := solutionFactory.CreateSolution([]int{1, 2, 3}, mockFitnessFunction)
		solution2 := solutionFactory.CreateSolution([]int{4, 5, 6}, mockFitnessFunction)
		individuals := []Solution[int]{*solution1, *solution2}

		population := populationFactory.CreatePopulation(individuals)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 2)
		assert.Equal(t, []int{1, 2, 3}, population.Individuals[0].Chromosome)
		assert.Equal(t, []int{4, 5, 6}, population.Individuals[1].Chromosome)
	})

	t.Run("create population with empty slice produces empty population", func(t *testing.T) {
		individuals := []Solution[int]{}

		population := populationFactory.CreatePopulation(individuals)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 0)
	})
}

func TestPopulationFactory_CreateEmptyPopulation(t *testing.T) {
	factory := NewPopulationFactory[int]()

	population := factory.CreateEmptyPopulation()

	require.NotNil(t, population)
	assert.Len(t, population.Individuals, 0)
	assert.NotNil(t, population.Individuals)
}

func TestPopulationFactory_CreateRandomPopulation(t *testing.T) {
	populationFactory := NewPopulationFactory[int]()
	solutionFactory := NewSolutionFactory[int]()

	t.Run("create random population with valid parameters", func(t *testing.T) {
		size := 5
		chromosomeLength := 10

		population := populationFactory.CreateRandomPopulation(
			size,
			chromosomeLength,
			solutionFactory,
			mockRandomIntGenerator,
			mockFitnessFunction,
		)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, size)

		for i, individual := range population.Individuals {
			assert.Len(t, individual.Chromosome, chromosomeLength, "Individual %d should have correct chromosome length", i)
			assert.NotNil(t, individual.fitnessFn, "Individual %d should have fitness function", i)
		}
	})

	t.Run("create random population with zero size produces empty population", func(t *testing.T) {
		population := populationFactory.CreateRandomPopulation(
			0,
			10,
			solutionFactory,
			mockRandomIntGenerator,
			mockFitnessFunction,
		)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 0)
	})

	t.Run("create random population with zero chromosome length creates individuals with empty chromosome", func(t *testing.T) {
		population := populationFactory.CreateRandomPopulation(
			3,
			0,
			solutionFactory,
			mockRandomIntGenerator,
			mockFitnessFunction,
		)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 3)

		for i, individual := range population.Individuals {
			assert.Len(t, individual.Chromosome, 0, "Individual %d should have empty chromosome", i)
		}
	})
}

func TestPopulation_WithDifferentTypes(t *testing.T) {
	t.Run("float64 chromosome type works", func(t *testing.T) {
		factory := NewPopulationFactory[float64]()
		solutionFactory := NewSolutionFactory[float64]()

		fitnessFunc := func(ch []float64) (float64, error) {
			sum := 0.0
			for _, val := range ch {
				sum += val
			}
			return sum, nil
		}

		randomGen := func() float64 {
			return rand.Float64() * 100
		}

		population := factory.CreateRandomPopulation(
			3, 4, solutionFactory, randomGen, fitnessFunc,
		)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 3)

		for _, individual := range population.Individuals {
			assert.Len(t, individual.Chromosome, 4)
		}
	})

	t.Run("int chromosome type works", func(t *testing.T) {
		factory := NewPopulationFactory[int]()
		solutionFactory := NewSolutionFactory[int]()

		fitnessFunc := func(ch []int) (float64, error) {
			return 10.0, nil
		}

		randomGen := func() int {
			return rand.Int() * 100
		}

		population := factory.CreateRandomPopulation(
			3, 4, solutionFactory, randomGen, fitnessFunc,
		)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 3)

		for _, individual := range population.Individuals {
			assert.Len(t, individual.Chromosome, 4)
		}
	})

	t.Run("string chromosome type works", func(t *testing.T) {
		factory := NewPopulationFactory[string]()
		solutionFactory := NewSolutionFactory[string]()

		fitnessFunc := func(ch []string) (float64, error) {
			return float64(len(ch)), nil
		}

		randomGen := func() string {
			return "test"
		}

		population := factory.CreateRandomPopulation(
			2, 3, solutionFactory, randomGen, fitnessFunc,
		)

		require.NotNil(t, population)
		assert.Len(t, population.Individuals, 2)

		for _, individual := range population.Individuals {
			assert.Len(t, individual.Chromosome, 3)
			for _, gene := range individual.Chromosome {
				assert.Equal(t, "test", gene)
			}
		}
	})
}

func TestPopulation_Integration(t *testing.T) {
	t.Run("full workflow test", func(t *testing.T) {
		populationFactory := NewPopulationFactory[int]()
		solutionFactory := NewSolutionFactory[int]()

		// Create a random oldPopulation
		oldPopulation := populationFactory.CreateRandomPopulation(
			3, 5, solutionFactory, mockRandomIntGenerator, mockFitnessFunction,
		)

		// Refresh fitness for all individuals
		err := oldPopulation.RefreshFitness()
		require.NoError(t, err)

		// Verify all individuals have calculated fitness
		for i, individual := range oldPopulation.Individuals {
			assert.GreaterOrEqual(t, individual.Fitness, 0.0, "Individual %d should have non-negative fitness", i)
		}

		// Create a new population from existing individuals
		newPopulation := populationFactory.CreatePopulation(oldPopulation.Individuals)
		assert.Equal(t, len(oldPopulation.Individuals), len(newPopulation.Individuals))

		// Verify that the new oldPopulation's individuals have the same fitness values
		for i, individual := range newPopulation.Individuals {
			assert.Equal(t, oldPopulation.Individuals[i].Fitness, individual.Fitness, "Individual %d should have the same fitness as the original", i)
			assert.Equal(t, oldPopulation.Individuals[i].Chromosome, individual.Chromosome, "Individual %d should have the same chromosome as the original", i)
		}
	})
}
