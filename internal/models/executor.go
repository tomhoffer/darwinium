package models

import (
	"cmp"
)

type GeneticAlgorithmExecutor[T cmp.Ordered] struct {
	population       *Population[T]
	fitnessEvaluator IFitnessEvaluator[T]
}

func NewGeneticAlgorithmExecutor[T cmp.Ordered](population *Population[T], fitnessEvaluator IFitnessEvaluator[T]) *GeneticAlgorithmExecutor[T] {
	return &GeneticAlgorithmExecutor[T]{
		population:       population,
		fitnessEvaluator: fitnessEvaluator,
	}
}

func (e *GeneticAlgorithmExecutor[T]) RefreshFitness() error {

	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return ErrPopulationEmpty
	}
	for i := range e.population.Individuals {
		fitness, err := e.fitnessEvaluator.Evaluate(&e.population.Individuals[i].Chromosome)
		if err != nil {
			return err
		}
		e.population.Individuals[i].Fitness = fitness
	}
	return nil
}
