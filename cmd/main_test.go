package main

import (
	"CS5260_Final_Project/util"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

type OutputData struct {
	Name                        string
	InitialStateFile            string
	TweakedParameter            string
	Value                       float64
	AverageTime                 float64
	TestRuns                    int
	AverageFinalExpectedUtility float64
}

func (od OutputData) toString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v", od.Name, od.InitialStateFile, od.TweakedParameter, od.Value,
		od.AverageFinalExpectedUtility, od.AverageTime, od.TestRuns)
}

type OutputList []OutputData

func (ol OutputList) toString() string {
	str := fmt.Sprintf("Name,IntialStateFile,TweakedParameter,Value,AverageFinalExpectedUtility,AverageTime,TestRuns")
	for _, data := range ol {
		str = fmt.Sprintf("%s\n%s", str, data.toString())
	}
	return str
}

func (ol OutputList) SaveToOutputFile(filename string) {
	//filename = fmt.Sprintf("output_data/%s", filename)

	outputString := ol.toString()

	if filename != "" {
		f, err := os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Here")
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(outputString); err != nil {
			fmt.Println("Here 2")
			log.Println(err)
		}
	}
}

const (
	countryName             = "Gondor"
	resourcesFile           = "../inputs/resources.csv"
	balancedCountriesFile   = "../inputs/countries_balanced.csv"
	imbalancedCountriesFile = "../inputs/countries_imbalanced.csv"
	randomCountriesFile     = "../inputs/countries_random.csv"
	outputScheduleFile      = ""
	numOutputSchedules      = 4
	depthBound              = 5
	frontierMaxSize         = 300
	roundsToSimulate        = 10

	testMaxDepth           = 10
	testMaxFrontierSize    = 500
	testMaxGamma           = 1
	testMinFailureConstant = -3
)

//
// Max Depth Tests
//

// TestTweakMaxDepthWithBalancedCountries
func TestTweakMaxDepthWithBalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for maxDepth := 3; maxDepth <= testMaxDepth; maxDepth++ {
		var totalTime time.Duration
		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, balancedCountriesFile,
				outputScheduleFile, numOutputSchedules, maxDepth, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for max depth: %v\n", i, maxDepth)
		}
		fmt.Printf("Finished all iterations for max depth: %v\n", maxDepth)

		output := OutputData{
			Name:                        "TweakMaxDepthWithBalancedCountries",
			InitialStateFile:            strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:            "MaxDepth",
			Value:                       float64(maxDepth),
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_depth_with_balanced_countries.csv")
}

// TestTweakMaxDepthWithImbalancedCountries
func TestTweakMaxDepthWithImbalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for maxDepth := 3; maxDepth <= testMaxDepth; maxDepth++ {
		var totalTime time.Duration
		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, imbalancedCountriesFile,
				outputScheduleFile, numOutputSchedules, maxDepth, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for max depth: %v\n", i, maxDepth)
		}
		fmt.Printf("Finished all iterations for max depth: %v\n", maxDepth)

		output := OutputData{
			Name:                        "TweakMaxDepthWithImbalancedCountries",
			InitialStateFile:            strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:            "MaxDepth",
			Value:                       float64(maxDepth),
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_depth_with_imbalanced_countries.csv")
}

// TestTweakMaxDepthWithRandomCountries
func TestTweakMaxDepthWithRandomCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for maxDepth := 3; maxDepth <= testMaxDepth; maxDepth++ {
		var totalTime time.Duration
		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, randomCountriesFile,
				outputScheduleFile, numOutputSchedules, maxDepth, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for max depth: %v\n", i, maxDepth)
		}
		fmt.Printf("Finished all iterations for max depth: %v\n", maxDepth)

		output := OutputData{
			Name:                        "TweakMaxDepthWithRandomCountries",
			InitialStateFile:            strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:            "MaxDepth",
			Value:                       float64(maxDepth),
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_depth_with_random_countries.csv")
}

//
// Max Frontier Size Tests
//

// TestTweakFrontierMaxSizeWithBalancedCountries
func TestTweakFrontierMaxSizeWithBalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for maxFrontierSize := 50; maxFrontierSize <= testMaxFrontierSize; maxFrontierSize += 50 {
		var totalTime time.Duration
		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, balancedCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, maxFrontierSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for max frontier size: %v\n", i, maxFrontierSize)
		}
		fmt.Printf("Finished all iterations max frontier size: %v\n", maxFrontierSize)

		output := OutputData{
			Name:                        "TweakMaxFrontierSizeWithBalancedCountries",
			InitialStateFile:            strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:            "FrontierMaxSize",
			Value:                       float64(maxFrontierSize),
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_frontier_size_with_balanced_countries.csv")
}

// TestTweakFrontierMaxSizeWithImbalancedCountries
func TestTweakFrontierMaxSizeWithImbalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for maxFrontierSize := 50; maxFrontierSize <= testMaxFrontierSize; maxFrontierSize += 50 {
		var totalTime time.Duration
		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, imbalancedCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, maxFrontierSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for max frontier size: %v\n", i, maxFrontierSize)
		}
		fmt.Printf("Finished all iterations max frontier size: %v\n", maxFrontierSize)

		output := OutputData{
			Name:                        "TweakFrontierMaxSizeWithImbalancedCountries",
			InitialStateFile:            strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:            "FrontierMaxSize",
			Value:                       float64(maxFrontierSize),
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_frontier_size_with_imbalanced_countries.csv")
}

// TestTweakFrontierMaxSizeWithRandomCountries
func TestTweakFrontierMaxSizeWithRandomCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for maxFrontierSize := 50; maxFrontierSize <= testMaxFrontierSize; maxFrontierSize += 50 {
		var totalTime time.Duration
		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, randomCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, maxFrontierSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for max frontier size: %v\n", i, maxFrontierSize)
		}
		fmt.Printf("Finished all iterations max frontier size: %v\n", maxFrontierSize)

		output := OutputData{
			Name:                        "TweakFrontierMaxSizeWithRandomCountries",
			InitialStateFile:            strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:            "FrontierMaxSize",
			Value:                       float64(maxFrontierSize),
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_frontier_size_with_random_countries.csv")
}

//
// Gamma Tests
//

// TestTweakGammaWithBalancedCountries
func TestTweakGammaWithBalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			FailureConstant: util.FailureCost,
			//Gamma:           util.Gamma,
			X0: util.X0,
			K:  util.K,
			L:  util.L,
		},
		UntweakableConstants: UntweakableConstants{
			TransformType: util.TransformType,
			TransferType:  util.TransferType,
		},
	}

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for gamma := 0.2; gamma <= testMaxGamma; gamma += 0.2 {
		var totalTime time.Duration
		testConstants.TweakableConstants.Gamma = gamma

		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, balancedCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for gamma: %v\n", i, gamma)
		}
		fmt.Printf("Finished all iterations for gamma: %v\n", gamma)

		output := OutputData{
			Name:                        "TweakGammaWithBalancedCountries",
			InitialStateFile:            strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:            "Gamma",
			Value:                       gamma,
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_gamma_with_balanced_countries.csv")
}

// TestTweakGammaWithImbalancedCountries
func TestTweakGammaWithImbalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for gamma := 0.2; gamma <= testMaxGamma; gamma += 0.2 {
		var totalTime time.Duration
		testConstants.TweakableConstants.Gamma = gamma

		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, imbalancedCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for gamma: %v\n", i, gamma)
		}
		fmt.Printf("Finished all iterations for gamma: %v\n", gamma)

		output := OutputData{
			Name:                        "TweakGammaWithBalancedCountries",
			InitialStateFile:            strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:            "Gamma",
			Value:                       gamma,
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_gamma_with_imbalanced_countries.csv")
}

// TestTweakGammaWithRandomCountries
func TestTweakGammaWithRandomCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
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

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for gamma := 0.2; gamma <= testMaxGamma; gamma += 0.2 {
		var totalTime time.Duration
		testConstants.TweakableConstants.Gamma = gamma

		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, randomCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for gamma: %v\n", i, gamma)
		}
		fmt.Printf("Finished all iterations for gamma: %v\n", gamma)

		output := OutputData{
			Name:                        "TweakGammaWithBalancedCountries",
			InitialStateFile:            strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:            "Gamma",
			Value:                       gamma,
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_gamma_with_random_countries.csv")
}

//
// Failure Constant Tests
//

// TestTweakFailureConstantWithBalancedCountries
func TestTweakFailureConstantWithBalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			//FailureCost: util.FailureCost,
			Gamma: util.Gamma,
			X0:    util.X0,
			K:     util.K,
			L:     util.L,
		},
		UntweakableConstants: UntweakableConstants{
			TransformType: util.TransformType,
			TransferType:  util.TransferType,
		},
	}

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for failureConstant := -0.5; failureConstant >= testMinFailureConstant; failureConstant -= 0.5 {
		var totalTime time.Duration
		testConstants.TweakableConstants.FailureConstant = failureConstant

		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, balancedCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for failure constant: %v\n", i, failureConstant)
		}
		fmt.Printf("Finished all iterations for failure constant: %v\n", failureConstant)

		output := OutputData{
			Name:                        "TweakFailureConstantWithBalancedCountries",
			InitialStateFile:            strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:            "FailureCost",
			Value:                       failureConstant,
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_failure_constant_with_balanced_countries.csv")
}

// TestTweakFailureConstantWithImbalancedCountries
func TestTweakFailureConstantWithImbalancedCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			//FailureCost: util.FailureCost,
			Gamma: util.Gamma,
			X0:    util.X0,
			K:     util.K,
			L:     util.L,
		},
		UntweakableConstants: UntweakableConstants{
			TransformType: util.TransformType,
			TransferType:  util.TransferType,
		},
	}

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for failureConstant := -0.5; failureConstant >= testMinFailureConstant; failureConstant -= 0.5 {
		var totalTime time.Duration
		testConstants.TweakableConstants.FailureConstant = failureConstant

		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, imbalancedCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for failure constant: %v\n", i, failureConstant)
		}
		fmt.Printf("Finished all iterations for failure constant: %v\n", failureConstant)

		output := OutputData{
			Name:                        "TweakFailureConstantWithBalancedCountries",
			InitialStateFile:            strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:            "Gamma",
			Value:                       failureConstant,
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_failure_constant_with_imbalanced_countries.csv")
}

// TestTweakFailureConstantWithRandomCountries
func TestTweakFailureConstantWithRandomCountries(t *testing.T) {
	MaxIteration := 3

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			//FailureCost: util.FailureCost,
			Gamma: util.Gamma,
			X0:    util.X0,
			K:     util.K,
			L:     util.L,
		},
		UntweakableConstants: UntweakableConstants{
			TransformType: util.TransformType,
			TransferType:  util.TransferType,
		},
	}

	outputList := make(OutputList, 0)
	expectedUtilityTotal := 0.0
	expectedUtilityCount := 0

	for failureConstant := -0.5; failureConstant >= testMinFailureConstant; failureConstant -= 0.5 {
		var totalTime time.Duration
		testConstants.TweakableConstants.FailureConstant = failureConstant

		for i := 0; i < MaxIteration; i++ {
			res := myCountryScheduler(countryName, resourcesFile, randomCountriesFile,
				outputScheduleFile, numOutputSchedules, depthBound, frontierMaxSize,
				roundsToSimulate, testConstants)
			totalTime += res.Time

			for _, eu := range res.FinalExceptedUtilities {
				expectedUtilityCount++
				expectedUtilityTotal += eu
			}

			fmt.Printf("Finished iteration %v for failure constant: %v\n", i, failureConstant)
		}
		fmt.Printf("Finished all iterations for failure constant: %v\n", failureConstant)

		output := OutputData{
			Name:                        "TweakFailureConstantWithBalancedCountries",
			InitialStateFile:            strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:            "FailureCost",
			Value:                       failureConstant,
			AverageTime:                 totalTime.Seconds() / float64(MaxIteration),
			TestRuns:                    MaxIteration,
			AverageFinalExpectedUtility: expectedUtilityTotal / float64(expectedUtilityCount),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_failure_constant_with_random_countries.csv")
}
