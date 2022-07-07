package main

import (
	"CS5260_Final_Project/countries"
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
	"time"
)

// Functions to sort the keys for the backup state map
type stateEntry []int

func (se stateEntry) Len() int           { return len(se) }
func (se stateEntry) Less(i, j int) bool { return se[i] > se[j] }
func (se stateEntry) Swap(i, j int)      { se[i], se[j] = se[j], se[i] }

func myCountryScheduler(myCountryName, resourcesFilename, initialStateFilename,
	outputSchduleFilename string, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate int,
	constants Constants) TestRun {

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
		SelfCountryName:       myCountryName,
	}

	resourcesMap := make(map[string]*resources.Resource)
	scheduler.CountriesMap = make(map[string]scheduler.Country)

	// Import resources
	lines, err := ReadCsv(resourcesFilename)
	if err != nil {
		panic(err)
	}

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

		scheduler.CountriesMap[countryName] = newCountry
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
	outputSchedules := 0
	//numOutputSchedules = 20 // Comes from params and can be removed
	printedOutputScheduleCount := 0
	numberOfTopLeaves := numOutputSchedules

	//depthBound = 3
	//frontierMaxSize = 200

	previousTopStates := make(map[int]simulator.WorldStates)
	//simulatorStartingStates := make(map[int]*simulator.WorldState)

	initialWorldState := simulator.ConvertToWorldStates(scheduler.CountriesMap)
	heap.Init(&initialWorldState)
	previousTopStates[0] = initialWorldState
	//previousTopStates := initialWorldState
	previousDepth := 0

	timeStart := time.Now()

	var expectedUtilities []float64

	for currentGeneration < roundsToSimulate {
		sim := simulator.InitializeSimulator(initialWorldState, frontierMaxSize, depthBound, currentGeneration)

		// Get all keys and sort
		var stateKeys stateEntry
		for key, _ := range previousTopStates {
			stateKeys = append(stateKeys, key)
		}
		sort.Sort(stateKeys)

		depth := previousDepth + 1
		full := false
		// Grab the top branches from the states, then work backwards if there's more room
		for _, key := range stateKeys {
			//if key == 3 {
			//	panic("key is 3!")
			//}
			for _, state := range previousTopStates[key] {
				// This will run the simulations to mock the schedules for each of the top states starting with the top
				full = sim.SimulateAllActionsFromState(scheduler.AvailableActions, scheduler.CountriesMap, state, state.Generation)
				if full {
					break
				}
			}
			if full {
				break
			}
		}
		//time.Sleep(time.Duration(500))

		// Increment depth bound and update previous depth - this will ensure we're also the same distance from the max depth
		// with each iteration of this loop
		depthBound++
		previousDepth = depth

		startingGeneration := depth
		topStates, backupStatesMap, size := sim.GetTopNLeaves(numberOfTopLeaves, startingGeneration)

		previousTopStates = backupStatesMap
		outputSchedules += len(topStates)

		sim.FlushQueue()
		rebuildStates := make([]*simulator.WorldState, size)
		i := 0
		for _, backupStates := range backupStatesMap {
			for _, state := range backupStates {
				//fmt.Printf("Backup States (%v), index: %v, ExpectedUtility: %v\n", state.Generation, state.Index, state.ExpectedUtility)
				state.Index = -1
				rebuildStates[i] = state
				i++
			}
		}

		initialWorldState = rebuildStates
		heap.Init(&initialWorldState)

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
		bestState := topStates[0]
		bestAction := bestState.GetGeneration(nextGeneration).Action
		bestAction.Action.Take(bestAction.ThisCountry, bestAction.OtherCountry)

		// Update the Quality score for the countries involved
		scheduler.CountriesMap[bestAction.ThisCountry.GetName()].SetQualityRating(
			quality.PerformQualityCalculation(bestAction.ThisCountry))
		if bestAction.OtherCountry != nil {
			scheduler.CountriesMap[bestAction.OtherCountry.GetName()].SetQualityRating(
				quality.PerformQualityCalculation(bestAction.OtherCountry))
		}

		// Increment to next generation
		nextGeneration++
		currentGeneration++
		time.Sleep(1 * time.Second)

		roundInformation := fmt.Sprintf("***************************  Round %v  ***************************",
			currentGeneration)
		roundInformation = fmt.Sprintf("%s\nDepth bound this round: %v\nFrontier size: %v\nTop "+
			"schedules count: %v\nSchedules:\n\n", roundInformation, depthBound-1, frontierMaxSize, numOutputSchedules)
		printToFile(outputSchduleFilename, roundInformation)

		expectedUtilities = make([]float64, len(topStates))
		// Print the top states for a check
		for i, state := range topStates {
			//if printedOutputScheduleCount > numOutputSchedules {
			//	break
			//}

			expectedUtilities[i] = state.ExpectedUtility

			//fmt.Println(state.ToString())
			//time.Sleep(1 * time.Second)
			//fmt.Println("**********")

			printToFile(outputSchduleFilename, state.ToString())
			printedOutputScheduleCount++
		}

		//if printedOutputScheduleCount > numOutputSchedules {
		//	break
		//}

		// Final round outputs
		//if currentGeneration == roundsToSimulate {

		//}
	}

	run.FinalExceptedUtilities = expectedUtilities
	run.MaxDepthAllRounds = depthBound
	run.Time = time.Since(timeStart)

	return run
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
	Time                   time.Duration
	DepthEachRound         int
	NumberOfRounds         int
	MaxDepthAllRounds      int
	FrontierMaxSize        int
	TopSchedulesEachRound  int
	FinalExceptedUtilities []float64

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

func main() {
	if len(os.Args) < 9 {
		str := "<selfCountryName>, <resourcesFilename>, <initialStateFilename>, <outputScheduleFilename>, " +
			"<numOutputSchedules>, <depthBound>, <frontierMaxSize>, <roundsToSimulate>"
		panic(fmt.Errorf("not enough command line parameters. need: %s", str))
	}

	//for i, arg := range os.Args {
	//	fmt.Printf("os.Args[%v]: '%v'\n", i, arg)
	//}

	selfCountry := strings.Trim(os.Args[1], ",")
	resourcesFilename := strings.Trim(os.Args[2], ",")
	initialStateFilename := strings.Trim(os.Args[3], ",")
	outputScheduleFilename := strings.Trim(os.Args[4], ",")
	numOutputSchedules := strings.Trim(os.Args[5], ",")
	depthBound := strings.Trim(os.Args[6], ",")
	frontierMaxSize := strings.Trim(os.Args[7], ",")
	roundsToSimulate := strings.Trim(os.Args[8], ",")

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

	myCountryScheduler(selfCountry, resourcesFilename, initialStateFilename, outputScheduleFilename,
		numOutputSchedulesInt, depthBoundInt, frontierMaxSizeInt, roundsToSimulateInt, constants)
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
