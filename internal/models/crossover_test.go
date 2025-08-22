package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSinglePointCrossover_Int tests the SinglePointCrossover with integer chromosomes.
func TestSinglePointCrossover_Int(t *testing.T) {
	crossover := NewSinglePointCrossover[int]()

	t.Run("valid crossover passes", func(t *testing.T) {
		parent1 := []int{1, 2, 3, 4, 5}
		parent2 := []int{6, 7, 8, 9, 10}

		offspring1, offspring2, err := crossover.Crossover(parent1, parent2)
		require.NoError(t, err)

		assert.Equal(t, len(parent1), len(offspring1))
		assert.Equal(t, len(parent2), len(offspring2))
		assert.NotEqual(t, parent1, offspring1)
		assert.NotEqual(t, parent2, offspring2)
		assert.NotEqual(t, parent1, offspring2)
		assert.NotEqual(t, parent2, offspring1)

		// Check that the sum of elements is conserved across parents and offspring.
		sumParents := 0
		for _, v := range parent1 {
			sumParents += v
		}
		for _, v := range parent2 {
			sumParents += v
		}

		sumOffspring := 0
		for _, v := range offspring1 {
			sumOffspring += v
		}
		for _, v := range offspring2 {
			sumOffspring += v
		}
		assert.Equal(t, sumParents, sumOffspring)
	})

	t.Run("single gene chromosomes pass", func(t *testing.T) {
		parent1 := []int{1}
		parent2 := []int{2}
		offspring1, offspring2, err := crossover.Crossover(parent1, parent2)
		require.NoError(t, err)
		assert.Equal(t, parent1, offspring1)
		assert.Equal(t, parent2, offspring2)
	})

	t.Run("different length chromosomes return error", func(t *testing.T) {
		parent1 := []int{1, 2, 3}
		parent2 := []int{4, 5}
		_, _, err := crossover.Crossover(parent1, parent2)
		require.Error(t, err)

		var ce *CrossoverError
		require.True(t, errors.As(err, &ce), "expected CrossoverError")

		var ice *InvalidChromosomeError
		require.True(t, errors.As(err, &ice), "expected InvalidChromosomeError")
	})

	t.Run("invalid parents return error", func(t *testing.T) {
		testCases := []struct {
			name    string
			parent1 []int
			parent2 []int
		}{
			{"parent1 is nil", nil, []int{1, 2, 3}},
			{"parent2 is nil", []int{1, 2, 3}, nil},
			{"both parents are nil", nil, nil},
			{"parent1 is empty", []int{}, []int{1, 2, 3}},
			{"parent2 is empty", []int{1, 2, 3}, []int{}},
			{"both parents are empty", []int{}, []int{}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, _, err := crossover.Crossover(tc.parent1, tc.parent2)
				require.Error(t, err)

				var ce *CrossoverError
				require.True(t, errors.As(err, &ce), "expected CrossoverError")

				var ice *InvalidChromosomeError
				require.True(t, errors.As(err, &ice), "expected InvalidChromosomeError to be wrapped")
			})
		}
	})
}

// TestSinglePointCrossover_String tests the SinglePointCrossover with string chromosomes.
func TestSinglePointCrossover_String(t *testing.T) {
	crossover := NewSinglePointCrossover[string]()

	t.Run("valid crossover on a string chromosome passes", func(t *testing.T) {
		parent1 := []string{"a", "b", "c", "d"}
		parent2 := []string{"e", "f", "g", "h"}

		offspring1, offspring2, err := crossover.Crossover(parent1, parent2)
		require.NoError(t, err)

		assert.Equal(t, len(parent1), len(offspring1))
		assert.Equal(t, len(parent2), len(offspring2))
		assert.NotEqual(t, parent1, offspring1)
		assert.NotEqual(t, parent2, offspring2)
	})
}
