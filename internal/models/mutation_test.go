package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleSwapMutator_Int tests the SimpleSwapMutator with integer chromosomes.
func TestSimpleSwapMutator_Int(t *testing.T) {
	mut := NewSimpleSwapMutator[int](1.0)

	t.Run("swaps two positions", func(t *testing.T) {
		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		err := mut.Mutate(&chromosome)
		require.NoError(t, err)

		// Verify that at least one position changed
		changed := false
		for i := range chromosome {
			if chromosome[i] != original[i] {
				changed = true
				break
			}
		}
		assert.True(t, changed, "chromosome should have been mutated")
	})

	t.Run("empty chromosome returns error", func(t *testing.T) {
		chromosome := []int{}

		err := mut.Mutate(&chromosome)
		require.Error(t, err)

		var me *MutationError
		require.True(t, errors.As(err, &me), "expected MutationError")
	})

	t.Run("nil chromosome returns error", func(t *testing.T) {
		var chromosome []int = nil

		err := mut.Mutate(&chromosome)
		require.Error(t, err)

		var me *MutationError
		require.True(t, errors.As(err, &me), "expected MutationError")
	})

	t.Run("single-gene chromosome returns error", func(t *testing.T) {
		chromosome := []int{7}

		err := mut.Mutate(&chromosome)
		require.Error(t, err)

		var me *MutationError
		require.True(t, errors.As(err, &me), "expected MutationError")
	})

	t.Run("2-gene chromosome mutates correctly", func(t *testing.T) {
		chromosome := []int{1, 2}

		err := mut.Mutate(&chromosome)
		assert.NoError(t, err)
		assert.Equal(t, 2, chromosome[0])
		assert.Equal(t, 1, chromosome[1])
	})
}

// TestSimpleSwapMutator_MutationRate tests the mutation rate functionality.
func TestSimpleSwapMutator_MutationRate(t *testing.T) {
	t.Run("default mutation rate is 0.01", func(t *testing.T) {
		mut := NewSimpleSwapMutator[int]()
		assert.Equal(t, 0.01, mut.mutationRate)
	})

	t.Run("custom mutation rate is set correctly", func(t *testing.T) {
		mut := NewSimpleSwapMutator[int](0.5)
		assert.Equal(t, 0.5, mut.mutationRate)
	})

	t.Run("high mutation rate always mutates", func(t *testing.T) {
		mut := NewSimpleSwapMutator[int](1.0) // 100% mutation rate
		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		err := mut.Mutate(&chromosome)
		require.NoError(t, err)

		// With 100% mutation rate, chromosome should always be mutated
		changed := false
		for i := range chromosome {
			if chromosome[i] != original[i] {
				changed = true
				break
			}
		}
		assert.True(t, changed, "chromosome should have been mutated with 100% mutation rate")
	})

	t.Run("zero mutation rate never mutates", func(t *testing.T) {
		mut := NewSimpleSwapMutator[int](0.0) // 0% mutation rate
		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		err := mut.Mutate(&chromosome)
		require.NoError(t, err)

		// With 0% mutation rate, chromosome should never be mutated
		assert.Equal(t, original, chromosome, "chromosome should not have been mutated with 0% mutation rate")
	})
}

// TestSimpleSwapMutator_String tests the SimpleSwapMutator with string chromosomes.
func TestSimpleSwapMutator_String(t *testing.T) {
	mut := NewSimpleSwapMutator[string](1.0)

	t.Run("swaps two positions", func(t *testing.T) {
		chromosome := []string{"a", "b", "c", "d"}
		original := append([]string(nil), chromosome...)

		err := mut.Mutate(&chromosome)
		require.NoError(t, err)

		// Verify that at least one position changed
		changed := false
		for i := range chromosome {
			if chromosome[i] != original[i] {
				changed = true
				break
			}
		}
		assert.True(t, changed, "chromosome should have been mutated")
	})
}
