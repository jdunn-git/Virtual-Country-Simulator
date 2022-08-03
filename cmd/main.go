package main

import (
	"CS5260_Final_Project/countries"
	"CS5260_Final_Project/manager"
	"CS5260_Final_Project/quality"
	"CS5260_Final_Project/resources"
	"CS5260_Final_Project/scheduler"
	"CS5260_Final_Project/simulator"
	"CS5260_Final_Project/transfers"
	"CS5260_Final_Project/transformations"
	"CS5260_Final_Project/util"
	"container/heap"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Functions to sort the keys for the backup state map
type stateEntry []int

func (se stateEntry) Len() int           { return len(se) }
func (se stateEntry) Less(i, j int) bool { return se[i] > se[j] }
func (se stateEntry) Swap(i, j int)      { se[i], se[j] = se[j], se[i] }

var wg sync.WaitGroup

func allCountrySchedulersWithGameManager(resourcesFilename, initialStateFilename, outputScheduleFilename,
	proposedScheduleFilename string, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beams int,
	constants Constants) TestRun {
	// Initialize Game Manager
	endChan := make(chan interface{})
	gameManager := manager.InitializeGameManager(endChan, outputScheduleFilename, proposedScheduleFilename)

	// Disburse Constants
	quality.Gamma = constants.Gamma
	scheduler.TransformType = constants.TransformType
	scheduler.TransferType = constants.TransferType
	scheduler.FailureConstant = constants.FailureConstant
	simulator.K = constants.K
	simulator.L = constants.L
	simulator.X0 = constants.X0

	// Generate initial Run struct
	run := TestRun{
		Time:                  0,
		DepthEachRound:        depthBound,
		NumberOfRounds:        roundsToSimulate,
		MaxDepthAllRounds:     0,
		FrontierMaxSize:       frontierMaxSize,
		TopSchedulesEachRound: numOutputSchedules,
		Constants:             constants,
		SelfCountryName:       "All Countries",
	}

	timeStart := time.Now()

	// Run each country as a separate go routine
	for countryName, _ := range quality.CountryQualityWeightsMap {
		wg.Add(1)
		go myCountryScheduler(countryName, resourcesFilename, initialStateFilename, outputScheduleFilename,
			numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beams, constants)
		time.Sleep(1 * time.Second)
		break
	}
	//go myCountryScheduler("Rohan", resourcesFilename, initialStateFilename, outputScheduleFilename,
	//	numOutputSchedulesInt, depthBoundInt, frontierMaxSizeInt, roundsToSimulateInt, beamCountInt, constants)
	//time.Sleep(1 * time.Second)
	// TODO will need to run the manager to listen for countries and schedules each round. Can control rounds here?
	fmt.Printf("Waiting for countries to finish\n")

	// Run the Game Manager
	fmt.Printf("Started the game manager\n")
	run.AverageQualityDelta = gameManager.Run(roundsToSimulate)

	// Wait for all countries to terminate
	wg.Wait()
	fmt.Printf("Done waiting for countries to finish\n")

	run.MaxDepthAllRounds = depthBound
	run.Time = time.Since(timeStart)

	//endChan <- 0
	fmt.Printf("Closed the game manager\n")

	return run
}

func myCountryScheduler(myCountryName, resourcesFilename, initialStateFilename,
	outputScheduleFilename string, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beams int,
	constants Constants) {

	// Close the waitgroup for this country whenever this process terminates
	defer wg.Done()
	//return TestRun{}

	selfCountry := &countries.Country{}

	resourcesMap := make(map[string]*resources.Resource)
	countryScheduler := &scheduler.CountryScheduler{
		CountriesMap: make(map[string]scheduler.Country),
	}

	// Import resources
	lines, err := ReadCsv(resourcesFilename)
	if err != nil {
		panic(err)
	}

	resourceNameList := make([]string, len(lines)-1)
	resourceCount := 0
	for _, line := range lines {
		// Skip the header line
		if line[0] == "Resource" {
			continue
		}

		resourceName := line[0]

		weight, err := strconv.ParseFloat(line[1], 32)
		if err != nil {
			panic(err)
		}

		resourcesMap[resourceName] = resources.NewResource(resourceName, float32(weight), false)
		resourceNameList[resourceCount] = resourceName
		resourceCount++
	}

	// Import countries
	lines, err = ReadCsv(initialStateFilename)
	if err != nil {
		panic(err)
	}

	headerMap := make(map[int]string)
	for _, line := range lines {
		// Generate the headerMap and skip the header line
		if line[0] == "Country" {
			for i, name := range line {
				headerMap[i] = name
			}
			continue
		}

		countryName := line[0]
		self := countryName == myCountryName

		countryResources := make([]*resources.CountryResource, len(resourcesMap))
		for i, value := range line {
			// Skip the first column since this is the country name
			if i < 1 {
				continue
			}

			amount, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}

			countryResources[i-1] = resources.NewCountryResource(
				resourcesMap[headerMap[i]],
				amount,
			)
		}

		newCountry := countries.NewCountry(countryName, countryResources, self)
		newCountry.SetQualityRating(quality.PerformQualityCalculation(newCountry))

		// If this country is "self" for its instance, then register with the game manager
		if self {
			man := manager.GetGameManager()
			man.Register(newCountry, newCountry.GetSendChan())
			selfCountry = newCountry
		}

		countryScheduler.CountriesMap[countryName] = newCountry
	}

	//for _, res := range resourcesMap {
	//	res.Print()
	//	fmt.Println()
	//}
	//
	//for _, country := range countriesMap {
	//	country.Print()
	//	fmt.Println()
	//}

	transformationMap := transformations.InitializeTransformations()
	transferMap := transfers.InitializeTransfers()

	scheduler.InitializeAvailableActions(transformationMap, transferMap)

	//gondor := countriesMap["Gondor"]
	//gondor.Print()
	//err = scheduler.PerformAction(gondor, nil, "Alloys Transformation")
	//err = scheduler.PerformAction(gondor, nil, "Alloys Transformation")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println()
	//gondor.Print()
	//fmt.Println()
	//
	//rohan := countriesMap["Rohan"]
	//rohan.Print()
	//err = scheduler.PerformAction(gondor, rohan, "Timber Transfer")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println()
	//gondor.Print()
	//fmt.Println()
	//rohan.Print()
	//fmt.Println()

	currentGeneration := 0
	nextGeneration := 1

	// Perform multiple iterations for frontier
	//outputSchedules := 0
	//numOutputSchedules = 20 // Comes from params and can be removed
	//printedOutputScheduleCount := 0
	//numberOfTopLeaves := numOutputSchedules

	//depthBound = 3
	//frontierMaxSize = 200

	previousTopStates := make(map[int]simulator.WorldStates)
	//simulatorStartingStates := make(map[int]*simulator.WorldState)

	initialWorldState := simulator.ConvertToWorldStates(countryScheduler.CountriesMap, resourceNameList)
	heap.Init(&initialWorldState)
	previousTopStates[0] = initialWorldState
	//previousTopStates := initialWorldState
	//previousDepth := 0

	//var expectedUtilities []float64
	round := 0

	// Schedule the world in rounds
	for round < roundsToSimulate {
		//maxDepthThisRound := depthBound * (round + 1)
		startingGeneration := 0

		if round == roundsToSimulate-1 {
			_ = round
		}

		// Simulate the world for maxDepthThisRound times
		sim := simulator.InitializeSimulator(countryScheduler.CountriesMap, initialWorldState, resourceNameList, frontierMaxSize, depthBound, startingGeneration)

		//tmpQualityRating := sim.CountriesMap[selfCountry.GetName()].GetQualityRating()
		//tmpUndiscountedReward := sim.CountriesMap[selfCountry.GetName()].GetUndiscountedReward()
		//tmpDiscountedReward := sim.CountriesMap[selfCountry.GetName()].GetDiscountedReward()
		//fmt.Printf("\tPre-simulation self.Quality: %v, self.UndiscountedReward: %v, self.DiscountedReward: %v\n",
		//	tmpQualityRating, tmpUndiscountedReward, tmpDiscountedReward)

		// Get all keys and sort
		var stateKeys stateEntry
		for key, _ := range previousTopStates {
			stateKeys = append(stateKeys, key)
		}
		sort.Sort(stateKeys)

		//depth := previousDepth + 1
		full := false
		// Grab the top branches from the states, then work backwards if there's more room
		for _, key := range stateKeys {
			//if key == 3 {
			//	panic("key is 3!")
			//}
			for _, state := range previousTopStates[key] {
				// This will run the simulations to mock the schedules for each of the top states starting with the top
				full = sim.SimulateAllActionsFromState(scheduler.AvailableActions, countryScheduler.CountriesMap, state,
					state.Generation, sim.MaxSize, beams, -1)
				_ = full
			}
		}
		//time.Sleep(time.Duration(500))

		//// Increment depth bound and update previous depth - this will ensure we're also the same distance from the max depth
		//// with each iteration of this loop
		//depthBound++
		//previousDepth = depth

		//startingGeneration := round * depthBound
		bestState, _ := sim.GetTopLeaf()
		//previousTopStates = backupStatesMap
		//outputSchedules = 1

		//sim.FlushQueue()
		//rebuildStates := make([]*simulator.WorldState, size)
		//i := 0
		//for _, backupStates := range backupStatesMap {
		//	for _, state := range backupStates {
		//		//fmt.Printf("Backup States (%v), index: %v, ExpectedUtility: %v\n", state.Generation, state.Index, state.ExpectedUtility)
		//		state.Index = -1
		//		rebuildStates[i] = state
		//		i++
		//	}
		//}

		//fmt.Println()
		//sim.PrintWorldStates()
		//states := sim.States
		//
		////heap.Init(&states)
		//for len(states) > 0 {
		//	state := heap.Pop(&states)
		//	_ = state
		//	//fmt.Printf("World value: %v\n", state.(*simulator.WorldState).WorldValue)
		//	//heap.Init(&states)
		//	//time.Sleep(200 * time.Millisecond)
		//}

		//fmt.Println()

		// Next we need to update the "real" world based on our best simulation
		// TODD: Still need to go through this and correct the best value calculation

		// Get the first action on the best simulated state
		//bestState := topStates[0]
		//bestAction := bestState.GetGeneration(nextGeneration).Action
		bestActions := bestState.GetAllActions()
		//bestAction.Action.Take(bestAction.ThisCountry, bestAction.OtherCountry)

		// PART TWO -- Instead of taking the best action, we will create a schedule and propose it, then take the responded action
		actionsToTake, worldState, proposalErr := selfCountry.ProposeSchedule(bestActions, bestState)
		if proposalErr != nil {
			panic(proposalErr)
		}

		// Take the returned action
		for _, action := range actionsToTake {
			// Get the country names from the action, but be careful to take them on the countryScheduler for this
			// country. If we use the country that's inside the action, then it will be the other country's instance
			// of the countryScheduler
			country1 := countryScheduler.CountriesMap[action.ThisCountry.GetName()]
			var country2 scheduler.Country
			if action.OtherCountry != nil {
				country2 = countryScheduler.CountriesMap[action.OtherCountry.GetName()]
			}
			takeErr := action.Action.Take(country1, country2)
			if takeErr != nil {
				panic(takeErr)
			}

			// Update the Quality score for the countries involved
			countryScheduler.CountriesMap[action.ThisCountry.GetName()].SetQualityRating(
				quality.PerformQualityCalculation(action.ThisCountry))
			if action.OtherCountry != nil {
				countryScheduler.CountriesMap[action.OtherCountry.GetName()].SetQualityRating(
					quality.PerformQualityCalculation(action.OtherCountry))
			}
		}

		resultWorldStates := simulator.WorldStates{worldState.(*simulator.WorldState)}
		previousTopStates[0] = resultWorldStates

		// Construct rebuildStates with the state from the schedule
		sim.FlushQueue()
		rebuildStates := make([]*simulator.WorldState, 1)
		//fmt.Printf("Backup States (%v), index: %v, ExpectedUtility: %v\n", state.Generation, state.Index, state.ExpectedUtility)
		worldState.SetAsNewInitialState()
		rebuildStates[0] = worldState.(*simulator.WorldState)

		initialWorldState = rebuildStates
		heap.Init(&initialWorldState)

		// Increment to next generation
		nextGeneration++
		currentGeneration++
		time.Sleep(1 * time.Second)

		//roundInformation := fmt.Sprintf("***************************  Round %v  ***************************",
		//	currentGeneration)
		//roundInformation = fmt.Sprintf("%s\nDepth bound this round: %v\nFrontier size: %v\nTop "+
		//	"schedules count: %v\nSchedules:\n\n", roundInformation, depthBound-1, frontierMaxSize, numOutputSchedules)
		//printToFile(outputScheduleFilename, roundInformation)

		//if printedOutputScheduleCount > numOutputSchedules {
		//	break
		//}

		// Final round outputs
		//if currentGeneration == roundsToSimulate {

		//}

		// Increment the round
		round++
	}

	return
}

func printToFile(filename, data string) {
	if filename != "" {
		f, err := os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(data); err != nil {
			log.Println(err)
		}
	}
}

type TestRun struct {
	Time                  time.Duration
	DepthEachRound        int
	NumberOfRounds        int
	MaxDepthAllRounds     int
	FrontierMaxSize       int
	TopSchedulesEachRound int
	AverageQualityDelta   float64

	SelfCountryName string

	Constants
}

type Constants struct {
	TweakableConstants
	UntweakableConstants
}

type TweakableConstants struct {
	FailureConstant float64
	Gamma           float64
	X0              float64
	K               float64
	L               float64
}

type UntweakableConstants struct {
	TransformType string
	TransferType  string
}

// importQualityWeights will read in quality weights for each country from filename, and will store the weights into the
// quality package CountryQualityWeightsMap. This is done one time for all goroutines to use, so it's a singleton
func importQualityWeights(filename string) {
	// Import resources
	lines, err := ReadCsv(filename)
	if err != nil {
		panic(err)
	}

	quality.CountryQualityWeightsMap = make(map[string]quality.CountryQualityWeights, len(lines)-1)

	headerMap := make(map[int]string)
	for _, line := range lines {
		qualityWeightsMap := make(map[string]float64, len(line))
		// Generate the headerMap and skip the header line
		if line[0] == "Country" {
			for i, name := range line {
				headerMap[i] = name
			}
			continue
		}

		countryName := line[0]

		for i, value := range line {
			// Skip the first column since this is the country name
			if i < 1 {
				continue
			}

			weight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic(err)
			}

			qualityWeightsMap[headerMap[i-1]] = weight
		}

		// Construct the quality weights struct for this country
		countryQualityWeights := quality.CountryQualityWeights{
			FoodQualityWeight:  qualityWeightsMap["FoodQualityWeight"],
			HousingQuality:     qualityWeightsMap["HousingQuality"],
			ElectronicsQuality: qualityWeightsMap["ElectronicsQuality"],
			MetalsQuality:      qualityWeightsMap["MetalsQuality"],
			LandQuality:        qualityWeightsMap["LandQuality"],
			MilitaryQuality:    qualityWeightsMap["MilitaryQuality"],
		}

		// Add the quality weights struct to the map
		quality.CountryQualityWeightsMap[countryName] = countryQualityWeights
	}

}

func main() {
	if len(os.Args) < 11 {
		str := "<selfCountryName>, <resourcesFilename>, <countryQualityWeightsFilename>, <initialStateFilename>, " +
			"<outputScheduleFilename>, <proposedScheduleFilename>, <numOutputSchedules>, <depthBound>, " +
			"<frontierMaxSize>, <roundsToSimulate>, <beamCount>"
		panic(fmt.Errorf("not enough command line parameters. need: %s", str))
	}

	//for i, arg := range os.Args {
	//	fmt.Printf("os.Args[%v]: '%v'\n", i, arg)
	//}

	selfCountry := strings.Trim(os.Args[1], ",")
	resourcesFilename := strings.Trim(os.Args[2], ",")
	countryQualityWeightsFilename := strings.Trim(os.Args[3], ",")
	initialStateFilename := strings.Trim(os.Args[4], ",")
	outputScheduleFilename := strings.Trim(os.Args[5], ",")
	proposedScheduleFilename := strings.Trim(os.Args[6], ",")
	numOutputSchedules := strings.Trim(os.Args[7], ",")
	depthBound := strings.Trim(os.Args[8], ",")
	frontierMaxSize := strings.Trim(os.Args[9], ",")
	roundsToSimulate := strings.Trim(os.Args[10], ",")
	beamCount := strings.Trim(os.Args[11], ",")

	numOutputSchedulesInt, err := strconv.Atoi(numOutputSchedules)
	if err != nil {
		panic(err)
	}

	depthBoundInt, err := strconv.Atoi(depthBound)
	if err != nil {
		panic(err)
	}

	frontierMaxSizeInt, err := strconv.Atoi(frontierMaxSize)
	if err != nil {
		panic(err)
	}

	roundsToSimulateInt, err := strconv.Atoi(roundsToSimulate)
	if err != nil {
		panic(err)
	}

	beamCountInt, err := strconv.Atoi(beamCount)
	if err != nil {
		panic(err)
	}

	constants := Constants{
		TweakableConstants: TweakableConstants{
			FailureConstant: util.FailureCost,
			Gamma:           util.Gamma,
			X0:              util.X0,
			K:               util.K,
			L:               util.L,
		},
		UntweakableConstants: UntweakableConstants{
			TransformType: util.TransformType,
			TransferType:  util.TransferType,
		},
	}

	// Import the country quality weights map
	importQualityWeights(countryQualityWeightsFilename)

	_ = selfCountry

	// Start all of the countries and game manager
	allCountrySchedulersWithGameManager(resourcesFilename, initialStateFilename, outputScheduleFilename,
		proposedScheduleFilename, numOutputSchedulesInt, depthBoundInt, frontierMaxSizeInt, roundsToSimulateInt,
		beamCountInt, constants)
}

// ReadCsv accepts a file and returns its content as a multi-dimensional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
