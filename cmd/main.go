package main

import (
	"fmt"
	"math/rand"

	"github.com/tomhoffer/darwinium/internal/models"
)

const (
	populationSize   = 100000
	chromosomeLength = 20
	generations      = 500
	tournamentSize   = 5
	elitismCount     = 1
	mutationRate     = 0.01 // Per-gene mutation probability
	geneMin          = -100
	geneMax          = 100
	numWorkers       = -1
)

// Custom chromosome type
type chromosomeType int

func generateRandomPopulation(populationFactory *models.PopulationFactory[chromosomeType], solutionFactory *models.SolutionFactory[chromosomeType]) *models.Population[chromosomeType] {
	var individuals []models.Solution[chromosomeType]
	for i := 0; i < populationSize; i++ {
		chromosome := make([]chromosomeType, chromosomeLength)
		for j := 0; j < chromosomeLength; j++ {
			chromosome[j] = chromosomeType(rand.Intn(geneMax-geneMin+1) + geneMin)
		}
		individuals = append(individuals, *solutionFactory.CreateSolution(chromosome))
	}
	return populationFactory.CreatePopulation(individuals)
}

func main() {
	// 1. Dependency injection
	solutionFactory := models.NewSolutionFactory[chromosomeType]()
	populationFactory := models.NewPopulationFactory[chromosomeType]()
	fitnessEvaluator := models.NewSimpleSumFitnessEvaluator[chromosomeType]()
	crossoverer := models.NewSinglePointCrossover[chromosomeType]()
	mutator := models.NewSimpleSwapMutator[chromosomeType](mutationRate)
	selector, err := models.NewTournamentSelector[chromosomeType](tournamentSize, elitismCount)
	if err != nil {
		panic(fmt.Sprintf("failed to create selector: %v", err))
	}

	// 2. Generate a random population
	population := generateRandomPopulation(populationFactory, solutionFactory)

	// 3. Instantiate the GeneticAlgorithmExecutor
	executor := models.NewGeneticAlgorithmExecutor(population, fitnessEvaluator, mutator, selector, crossoverer, generations, numWorkers)

	// 4. Run the GA loop
	finalPopulation, err := executor.Loop(generations)
	if err != nil {
		panic(fmt.Sprintf("genetic algorithm failed: %v", err))
	}

	// 5. Print the final result
	bestSolution, err := finalPopulation.BestSolution()
	if err != nil {
		panic(fmt.Sprintf("failed to get best solution: %v", err))
	}

	fmt.Printf("Best solution found with fitness %.2f:\n", bestSolution.Fitness)
	fmt.Printf("Chromosome: %v\n", bestSolution.Chromosome)
}
