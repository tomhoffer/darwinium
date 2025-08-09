package models

import (
	"cmp"

	"github.com/tomhoffer/darwinium/internal/utils"
)

type IFitnessEvaluator[T cmp.Ordered] interface {
	Evaluate(chromosome []T) (float64, error)
}

type SimpleSumFitnessEvaluator[T cmp.Ordered] struct{}

func NewSimpleSumFitnessEvaluator[T cmp.Ordered]() *SimpleSumFitnessEvaluator[T] {
	return &SimpleSumFitnessEvaluator[T]{}
}

func (s SimpleSumFitnessEvaluator[T]) Evaluate(chromosome []T) (float64, error) {
	var sum float64
	for _, val := range chromosome {
		converted, err := utils.ConvertToFloat64(val)
		if err != nil {
			return 0, err
		}
		sum += converted
	}
	return sum, nil
}
