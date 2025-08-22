package models

import (
	"cmp"
	"math/rand"
)

type GeneticAlgorithmExecutor[T cmp.Ordered] struct {
	population       *Population[T]
	fitnessEvaluator IFitnessEvaluator[T]
	mutator          IMutator[T]
	selector         ISelector[T]
	crossover        ICrossover[T]
}

func NewGeneticAlgorithmExecutor[T cmp.Ordered](population *Population[T], fitnessEvaluator IFitnessEvaluator[T], mutator IMutator[T], selector ISelector[T], crossover ICrossover[T]) *GeneticAlgorithmExecutor[T] {
	return &GeneticAlgorithmExecutor[T]{
		population:       population,
		fitnessEvaluator: fitnessEvaluator,
		mutator:          mutator,
		selector:         selector,
		crossover:        crossover,
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

func (e *GeneticAlgorithmExecutor[T]) PerformCrossover() (*Population[T], error) {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return nil, NewCrossoverError("cannot perform crossover on empty population", ErrPopulationEmpty)
	}

	individuals := e.population.Individuals
	rand.Shuffle(len(individuals), func(i, j int) {
		individuals[i], individuals[j] = individuals[j], individuals[i]
	})

	offspringPopulation := &Population[T]{
		Individuals: make([]Solution[T], 0, len(individuals)),
	}

	for i := 0; i < len(individuals); i += 2 {
		if i+1 >= len(individuals) {
			offspringPopulation.Individuals = append(offspringPopulation.Individuals, individuals[i])
			break
		}

		parent1 := individuals[i]
		parent2 := individuals[i+1]

		offspringChr1, offspringChr2, err := e.crossover.Crossover(parent1.Chromosome, parent2.Chromosome)
		if err != nil {
			return nil, err
		}

		offspringPopulation.Individuals = append(offspringPopulation.Individuals, Solution[T]{Chromosome: offspringChr1}, Solution[T]{Chromosome: offspringChr2})
	}

	return offspringPopulation, nil
}
