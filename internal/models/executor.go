package models

import "cmp"

type GeneticAlgorithmExecutor[T cmp.Ordered] struct {
	population       Population[T]
	fitnessEvaluator IFitnessEvaluator[T]
}

func NewGeneticAlgorithmExecutor[T cmp.Ordered](population Population[T], fitnessEvaluator IFitnessEvaluator[T]) *GeneticAlgorithmExecutor[T] {
	return &GeneticAlgorithmExecutor[T]{
		population:       population,
		fitnessEvaluator: fitnessEvaluator,
	}
}

func (e *GeneticAlgorithmExecutor[T]) RefreshFitness() error {
	for i, ind := range e.population.Individuals {
		fitness, err := e.fitnessEvaluator.Evaluate(ind.Chromosome)
		if err != nil {
			return err
		}
		e.population.Individuals[i].Fitness = fitness
	}
	return nil
}
