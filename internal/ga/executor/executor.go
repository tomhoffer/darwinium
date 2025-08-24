package executor

import (
	"cmp"
	"context"
	"fmt"
	"math/rand"

	progressbar "github.com/schollz/progressbar/v3"
	"github.com/tomhoffer/darwinium/internal/core"
	"github.com/tomhoffer/darwinium/internal/ga/crossover"
	"github.com/tomhoffer/darwinium/internal/ga/fitness"
	"github.com/tomhoffer/darwinium/internal/ga/mutation"
	"github.com/tomhoffer/darwinium/internal/ga/selection"
	"golang.org/x/sync/errgroup"
)

type GeneticAlgorithmExecutor[T cmp.Ordered] struct {
	population       *core.Population[T]
	fitnessEvaluator fitness.IFitnessEvaluator[T]
	mutator          mutation.IMutator[T]
	selector         selection.ISelector[T]
	crossover        crossover.ICrossover[T]
	generations      int
	numWorkers       int
}

func NewGeneticAlgorithmExecutor[T cmp.Ordered](population *core.Population[T], fitnessEvaluator fitness.IFitnessEvaluator[T], mutator mutation.IMutator[T], selector selection.ISelector[T], crossover crossover.ICrossover[T], generations int, numWorkers ...int) *GeneticAlgorithmExecutor[T] {
	// Default to 1 worker if not specified
	workerCount := 1
	if len(numWorkers) > 0 {
		workerCount = numWorkers[0]
	}

	return &GeneticAlgorithmExecutor[T]{
		population:       population,
		fitnessEvaluator: fitnessEvaluator,
		mutator:          mutator,
		selector:         selector,
		crossover:        crossover,
		generations:      generations,
		numWorkers:       workerCount,
	}
}

func (e *GeneticAlgorithmExecutor[T]) RefreshFitness(ctx context.Context) error {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return core.ErrPopulationEmpty
	}

	// Run fitness evaluation in goroutines with limited concurrency
	g, gCtx := errgroup.WithContext(ctx)

	if e.numWorkers != -1 {
		g.SetLimit(e.numWorkers)
	}

	for i := range e.population.Individuals {
		individualIndex := i // explicit capture
		g.Go(func() error {
			fitness, err := e.fitnessEvaluator.Evaluate(gCtx, &e.population.Individuals[individualIndex].Chromosome)
			if err != nil {
				return err
			}
			e.population.Individuals[individualIndex].Fitness = fitness
			return nil
		})
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return fitness.NewFitnessEvaluationError("failed to evaluate fitness", err)
	}
	return nil
}

func (e *GeneticAlgorithmExecutor[T]) PerformMutation(ctx context.Context) error {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return core.ErrPopulationEmpty
	}

	// Run mutation in goroutines with limited concurrency
	g, gCtx := errgroup.WithContext(ctx)

	if e.numWorkers != -1 {
		g.SetLimit(e.numWorkers)
	}

	for i := range e.population.Individuals {
		individualIndex := i // explicit capture
		g.Go(func() error {
			err := e.mutator.Mutate(gCtx, &e.population.Individuals[individualIndex].Chromosome)
			if err != nil {
				return err
			}
			return nil
		})
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return fitness.NewFitnessEvaluationError("failed to evaluate fitness", err)
	}
	return nil
}

func (e *GeneticAlgorithmExecutor[T]) PerformSelection() (*core.Population[T], error) {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return nil, core.ErrPopulationEmpty
	}
	newPopulation, err := e.selector.Select(e.population)
	if err != nil {
		return nil, err
	}
	return newPopulation, nil
}

func (e *GeneticAlgorithmExecutor[T]) PerformCrossover() (*core.Population[T], error) {
	if e.population == nil || e.population.Individuals == nil || len(e.population.Individuals) == 0 {
		return nil, crossover.NewCrossoverError("cannot perform crossover on empty population", core.ErrPopulationEmpty)
	}

	individuals := e.population.Individuals
	rand.Shuffle(len(individuals), func(i, j int) {
		individuals[i], individuals[j] = individuals[j], individuals[i]
	})

	offspringPopulation := &core.Population[T]{
		Individuals: make([]core.Solution[T], 0, len(individuals)),
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

		offspringPopulation.Individuals = append(offspringPopulation.Individuals, core.Solution[T]{Chromosome: offspringChr1}, core.Solution[T]{Chromosome: offspringChr2})
	}

	return offspringPopulation, nil
}

// Loop runs the genetic algorithm for the specified number of generations.
// It performs fitness evaluation, selection, crossover, and mutation in each generation.
// The method returns the final population and any error that occurred during execution.
func (e *GeneticAlgorithmExecutor[T]) Loop(ctx context.Context, generations int) (*core.Population[T], error) {
	bar := progressbar.Default(int64(generations))
	fmt.Println("Starting genetic algorithm...")
	for i := 0; i < generations; i++ {
		err := bar.Add(1)
		if err != nil {
			return nil, err
		}

		// a. Refresh fitness for the current population
		if err := e.RefreshFitness(ctx); err != nil {
			return nil, fmt.Errorf("failed to refresh fitness at generation %d: %w", i, err)
		}

		// b. Find and print the best fitness in the current generation
		_, err = e.population.BestFitness()
		if err != nil {
			return nil, fmt.Errorf("failed to get best fitness at generation %d: %w", i, err)
		}

		// c. Perform selection
		selectedPopulation, err := e.PerformSelection()
		if err != nil {
			return nil, fmt.Errorf("failed to perform selection at generation %d: %w", i, err)
		}
		e.population = selectedPopulation

		// d. Perform crossover
		offspringPopulation, err := e.PerformCrossover()
		if err != nil {
			return nil, fmt.Errorf("failed to perform crossover at generation %d: %w", i, err)
		}
		e.population = offspringPopulation

		// e. Perform mutation
		if err := e.PerformMutation(ctx); err != nil {
			return nil, fmt.Errorf("failed to perform mutation at generation %d: %w", i, err)
		}
	}
	// f. Re-evaluate fitness for the new population (after crossover + mutation)
	if err := e.RefreshFitness(ctx); err != nil {
		return nil, fmt.Errorf("failed to refresh fitness: %w", err)
	}

	fmt.Println("\nFinished genetic algorithm!")
	return e.population, nil
}
