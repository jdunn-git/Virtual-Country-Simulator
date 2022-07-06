package main

import (
	"CS5260_Final_Project/countries"
	"CS5260_Final_Project/quality"
	"CS5260_Final_Project/resources"
	"CS5260_Final_Project/scheduler"
	"CS5260_Final_Project/simulator"
	"CS5260_Final_Project/transfers"
	"CS5260_Final_Project/transformations"
	"container/heap"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	resourcesMap := make(map[string]*resources.Resource)
	countriesMap := make(map[string]*countries.Country)

	// Import resources
	lines, err := ReadCsv("resources.csv")
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
	lines, err = ReadCsv("countries_2.csv")
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
		self, err := strconv.ParseBool(line[1])
		if err != nil {
			log.Fatal(err)
		}

		countryResources := make([]*resources.CountryResource, len(resourcesMap))
		for i, value := range line {
			// Skip the first two columns since this is the country name and self columns respectively
			if i < 2 {
				continue
			}

			amount, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}

			countryResources[i-2] = resources.NewCountryResource(
				resourcesMap[headerMap[i]],
				amount,
			)
		}

		newCountry := countries.NewCountry(countryName, countryResources, self)
		newCountry.SetQualityRating(quality.PerformQualityCalculation(newCountry))

		countriesMap[countryName] = newCountry
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

	// TODO: Make this more configurable
	maxSize := 50
	maxDepth := 4
	sim, worldState := simulator.InitializeSimulator(countriesMap, maxSize, maxDepth)

	// Grab the top three branches from the states

	sim.SimulateAllActionsFromState(scheduler.AvailableActions, countriesMap, worldState, 1)
	time.Sleep(time.Duration(500))
	topStates, keepStates := sim.GetTopNLeaves(3, 1)
	sim.FlushQueue()
	rebuildStates := make([]*simulator.WorldState, len(keepStates))
	i := 0
	for _, state := range keepStates {
		fmt.Printf("Keep States (%v), index: %v, WorldValue: %v\n", state.Generation, state.Index, state.WorldValue)
		state.Index = -1
		rebuildStates[i] = state
		i++
	}
	sim.RebuildQueue(rebuildStates)
	_ = topStates

	// TODO: Don't forget to progress the "actual" world alongside the simulation
	//  - Will need to get the first gen parent and then use that to take the action

	fmt.Println()
	//sim.PrintWorldStates()
	states := sim.States
	for _, _ = range states {
		heap.Init(&states)
	}
	//heap.Init(&states)
	for len(states) > 0 {
		state := heap.Pop(&states)
		fmt.Printf("World value: %v\n", state.(*simulator.WorldState).WorldValue)
		//heap.Init(&states)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()

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
