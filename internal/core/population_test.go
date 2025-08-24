package core

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for testing
func mockFitnessFunction(ch []int) (float64, error) {
	return 999.0, nil
}

func mockRandomIntGenerator() int {
	return rand.Intn(100)
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
		solution1 := solutionFactory.CreateSolution([]int{1, 2, 3})
		solution2 := solutionFactory.CreateSolution([]int{4, 5, 6})
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
