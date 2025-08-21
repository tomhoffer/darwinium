package models

import (
	"cmp"
)

type GeneticAlgorithmExecutor[T cmp.Ordered] struct {
	population       *Population[T]
	fitnessEvaluator IFitnessEvaluator[T]
	mutator          IMutator[T]
	selector         ISelector[T]
}

func NewGeneticAlgorithmExecutor[T cmp.Ordered](population *Population[T], fitnessEvaluator IFitnessEvaluator[T], mutator IMutator[T], selector ISelector[T]) *GeneticAlgorithmExecutor[T] {
	return &GeneticAlgorithmExecutor[T]{
		population:       population,
		fitnessEvaluator: fitnessEvaluator,
		mutator:          mutator,
		selector:         selector,
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

func (e *GeneticAlgorithmExecutor[T]) PerformMutation() error {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return ErrPopulationEmpty
	}
	for i := range e.population.Individuals {
		err := e.mutator.Mutate(&e.population.Individuals[i].Chromosome)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *GeneticAlgorithmExecutor[T]) PerformSelection() (*Population[T], error) {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return nil, ErrPopulationEmpty
	}
	newPopulation, err := e.selector.Select(e.population)
	if err != nil {
		return nil, err
	}
	return newPopulation, nil
}
