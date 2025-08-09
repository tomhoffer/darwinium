package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleSumFitnessEvaluator_Int(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[int]{}

	t.Run("positive integers", func(t *testing.T) {
		chromosome := []int{1, 2, 3, 4, 5}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 15.0, result)
	})

	t.Run("negative integers", func(t *testing.T) {
		chromosome := []int{-1, -2, -3}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, -6.0, result)
	})

	t.Run("mixed positive and negative", func(t *testing.T) {
		chromosome := []int{10, -5, 3, -2}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 6.0, result)
	})

	t.Run("empty chromosome", func(t *testing.T) {
		chromosome := []int{}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("zero values", func(t *testing.T) {
		chromosome := []int{0, 0, 0}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Int8(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[int8]{}

	t.Run("int8 values", func(t *testing.T) {
		chromosome := []int8{10, 20, -5, 15}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 40.0, result)
	})

	t.Run("int8 max and min values", func(t *testing.T) {
		chromosome := []int8{127, -128}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, -1.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Int16(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[int16]{}

	t.Run("int16 values", func(t *testing.T) {
		chromosome := []int16{1000, 2000, -500}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 2500.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Int32(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[int32]{}

	t.Run("int32 values", func(t *testing.T) {
		chromosome := []int32{100000, 200000, -50000}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 250000.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Int64(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[int64]{}

	t.Run("int64 values", func(t *testing.T) {
		chromosome := []int64{1000000000, 2000000000, -500000000}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 2500000000.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Uint(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[uint]{}

	t.Run("uint values", func(t *testing.T) {
		chromosome := []uint{1, 2, 3, 4, 5}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 15.0, result)
	})

	t.Run("large uint values", func(t *testing.T) {
		chromosome := []uint{1000, 2000, 3000}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 6000.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Uint8(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[uint8]{}

	t.Run("uint8 values", func(t *testing.T) {
		chromosome := []uint8{10, 20, 30, 40}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 100.0, result)
	})

	t.Run("uint8 max value", func(t *testing.T) {
		chromosome := []uint8{255, 255}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 510.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Uint16(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[uint16]{}

	t.Run("uint16 values", func(t *testing.T) {
		chromosome := []uint16{1000, 2000, 3000}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 6000.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Uint32(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[uint32]{}

	t.Run("uint32 values", func(t *testing.T) {
		chromosome := []uint32{100000, 200000, 300000}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 600000.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Uint64(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[uint64]{}

	t.Run("uint64 values", func(t *testing.T) {
		chromosome := []uint64{1000000000, 2000000000, 3000000000}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 6000000000.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Uintptr(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[uintptr]{}

	t.Run("uintptr values", func(t *testing.T) {
		chromosome := []uintptr{100, 200, 300}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 600.0, result)
	})
}

func TestSimpleSumFitnessEvaluator_Float32(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[float32]{}

	t.Run("positive float32 values", func(t *testing.T) {
		chromosome := []float32{1.5, 2.5, 3.0}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 7.0, result)
	})

	t.Run("negative float32 values", func(t *testing.T) {
		chromosome := []float32{-1.5, -2.5, -1.0}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, -5.0, result)
	})

	t.Run("mixed float32 values", func(t *testing.T) {
		chromosome := []float32{10.5, -5.5, 2.0}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 7.0, result)
	})

	t.Run("float32 precision", func(t *testing.T) {
		chromosome := []float32{0.1, 0.2, 0.3}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.InDelta(t, 0.6, result, 0.0001)
	})
}

func TestSimpleSumFitnessEvaluator_Float64(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[float64]{}

	t.Run("positive float64 values", func(t *testing.T) {
		chromosome := []float64{1.5, 2.5, 3.0}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 7.0, result)
	})

	t.Run("negative float64 values", func(t *testing.T) {
		chromosome := []float64{-1.5, -2.5, -1.0}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, -5.0, result)
	})

	t.Run("high precision float64", func(t *testing.T) {
		chromosome := []float64{1.123456789, 2.987654321, 0.888888888}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.InDelta(t, 4.999999998, result, 0.000000001)
	})

	t.Run("very large float64 values", func(t *testing.T) {
		chromosome := []float64{1e10, 2e10, 3e10}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 6e10, result)
	})

	t.Run("very small float64 values", func(t *testing.T) {
		chromosome := []float64{1e-10, 2e-10, 3e-10}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.InDelta(t, 6e-10, result, 1e-15)
	})
}

func TestSimpleSumFitnessEvaluator_String(t *testing.T) {
	evaluator := SimpleSumFitnessEvaluator[string]{}

	t.Run("valid numeric strings", func(t *testing.T) {
		chromosome := []string{"1.5", "2.5", "3.0"}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 7.0, result)
	})

	t.Run("integer strings", func(t *testing.T) {
		chromosome := []string{"10", "20", "30"}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 60.0, result)
	})

	t.Run("negative numeric strings", func(t *testing.T) {
		chromosome := []string{"-1.5", "-2.5", "4.0"}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("scientific notation strings", func(t *testing.T) {
		chromosome := []string{"1e2", "2e1", "3e0"}
		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 123.0, result)
	})

	t.Run("invalid string should return error", func(t *testing.T) {
		chromosome := []string{"1.5", "invalid", "3.0"}
		result, err := evaluator.Evaluate(chromosome)

		assert.Error(t, err)
		assert.Equal(t, 0.0, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("empty string should return error", func(t *testing.T) {
		chromosome := []string{"1.5", "", "3.0"}
		result, err := evaluator.Evaluate(chromosome)

		assert.Error(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("non-numeric string should return error", func(t *testing.T) {
		chromosome := []string{"abc", "def"}
		result, err := evaluator.Evaluate(chromosome)

		assert.Error(t, err)
		assert.Equal(t, 0.0, result)
	})
}

// Test with custom ordered types
type CustomInt int
type CustomFloat float64
type CustomString string

func TestSimpleSumFitnessEvaluator_CustomTypes(t *testing.T) {
	t.Run("custom int type", func(t *testing.T) {
		evaluator := SimpleSumFitnessEvaluator[CustomInt]{}
		chromosome := []CustomInt{CustomInt(10), CustomInt(20), CustomInt(30)}

		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 60.0, result)
	})

	t.Run("custom float type", func(t *testing.T) {
		evaluator := SimpleSumFitnessEvaluator[CustomFloat]{}
		chromosome := []CustomFloat{CustomFloat(1.5), CustomFloat(2.5), CustomFloat(3.0)}

		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 7.0, result)
	})

	t.Run("custom string type", func(t *testing.T) {
		evaluator := SimpleSumFitnessEvaluator[CustomString]{}
		chromosome := []CustomString{CustomString("10"), CustomString("20"), CustomString("30")}

		result, err := evaluator.Evaluate(chromosome)

		require.NoError(t, err)
		assert.Equal(t, 60.0, result)
	})

	t.Run("custom string type with invalid value", func(t *testing.T) {
		evaluator := SimpleSumFitnessEvaluator[CustomString]{}
		chromosome := []CustomString{CustomString("10"), CustomString("invalid"), CustomString("30")}

		result, err := evaluator.Evaluate(chromosome)

		assert.Error(t, err)
		assert.Equal(t, 0.0, result)
	})
}
