package models

import (
	"cmp"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleSwapMutator_Int tests the SimpleSwapMutator with integer chromosomes.
func TestSimpleSwapMutator_Int(t *testing.T) {
	mut := NewSimpleSwapMutator[int](1.0)

	t.Run("swaps two positions", func(t *testing.T) {
		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		err := mut.Mutate(context.Background(), &chromosome)
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

		err := mut.Mutate(context.Background(), &chromosome)
		require.Error(t, err)

		var me *MutationError
		require.True(t, errors.As(err, &me), "expected MutationError")
	})

	t.Run("nil chromosome returns error", func(t *testing.T) {
		var chromosome []int = nil

		err := mut.Mutate(context.Background(), &chromosome)
		require.Error(t, err)

		var me *MutationError
		require.True(t, errors.As(err, &me), "expected MutationError")
	})

	t.Run("single-gene chromosome returns error", func(t *testing.T) {
		chromosome := []int{7}

		err := mut.Mutate(context.Background(), &chromosome)
		require.Error(t, err)

		var me *MutationError
		require.True(t, errors.As(err, &me), "expected MutationError")
	})

	t.Run("2-gene chromosome mutates correctly", func(t *testing.T) {
		chromosome := []int{1, 2}

		err := mut.Mutate(context.Background(), &chromosome)
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

		err := mut.Mutate(context.Background(), &chromosome)
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

		err := mut.Mutate(context.Background(), &chromosome)
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

		err := mut.Mutate(context.Background(), &chromosome)
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

// Mock sleeping mutator for testing context cancellation
type mockSleepingMutator[T cmp.Ordered] struct {
	sleepDuration time.Duration
}

func (m *mockSleepingMutator[T]) Mutate(ctx context.Context, chromosome *[]T) error {
	// Sleep for the specified duration to simulate long-running mutation
	select {
	case <-time.After(m.sleepDuration):
		// Sleep completed, continue with mutation
	case <-ctx.Done():
		// Context was cancelled during sleep
		return NewMutationError("context cancelled", ctx.Err())
	}

	// Simple mutation: swap first two genes if possible
	if chromosome != nil && *chromosome != nil && len(*chromosome) >= 2 {
		ch := *chromosome
		ch[0], ch[1] = ch[1], ch[0]
	}
	return nil
}

// TestSimpleSwapMutator_ContextCancellation tests context cancellation functionality
func TestSimpleSwapMutator_ContextCancellation(t *testing.T) {
	t.Run("immediate context cancellation returns error", func(t *testing.T) {
		mut := NewSimpleSwapMutator[int](1.0)
		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		// Create a context that's already cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := mut.Mutate(ctx, &chromosome)
		require.Error(t, err)

		var me *MutationError
		assert.ErrorAs(t, err, &me)
		assert.Contains(t, me.Message, "context cancelled")

		// Chromosome should remain unchanged
		assert.Equal(t, original, chromosome)
	})

	t.Run("context cancellation during sleep returns error", func(t *testing.T) {
		// Use a sleeping mock mutator
		mut := &mockSleepingMutator[int]{
			sleepDuration: 100 * time.Millisecond,
		}

		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		// Create a context that will be cancelled during the sleep
		ctx, cancel := context.WithCancel(context.Background())

		// Start mutation in a goroutine
		var err error
		done := make(chan bool)

		go func() {
			err = mut.Mutate(ctx, &chromosome)
			done <- true
		}()

		// Cancel the context after a short delay (shorter than sleep duration)
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		// Wait for mutation to complete
		<-done

		require.Error(t, err)
		var me *MutationError
		assert.ErrorAs(t, err, &me)
		assert.Contains(t, me.Message, "context cancelled")

		// Chromosome should remain unchanged
		assert.Equal(t, original, chromosome)
	})

	t.Run("context timeout returns error", func(t *testing.T) {
		// Use a sleeping mock mutator
		mut := &mockSleepingMutator[int]{
			sleepDuration: 100 * time.Millisecond,
		}

		chromosome := []int{1, 2, 3, 4, 5}
		original := append([]int(nil), chromosome...)

		// Create a context with a timeout shorter than the sleep duration
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := mut.Mutate(ctx, &chromosome)
		require.Error(t, err)

		var me *MutationError
		assert.ErrorAs(t, err, &me)
		assert.Contains(t, me.Message, "context cancelled")

		// Chromosome should remain unchanged
		assert.Equal(t, original, chromosome)
	})
}
