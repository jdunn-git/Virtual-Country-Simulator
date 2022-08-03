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
	Name                string
	InitialStateFile    string
	TweakedParameter    string
	Value               float64
	MinTime             float64
	MaxTime             float64
	AverageTime         float64
	TestRuns            int
	MinQualityDelta     float64
	MaxQualityDelta     float64
	AverageQualityDelta float64
	//MinQualityDelta float64
	//MaxQualityDelta float64
	//AverageQualityDelta float64
}

func (od OutputData) toString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v", od.Name, od.InitialStateFile, od.TweakedParameter,
		od.Value, od.MinQualityDelta, od.MaxQualityDelta, od.AverageQualityDelta, od.MinTime,
		od.MaxTime, od.AverageTime, od.TestRuns)
}

func (od OutputData) SaveToOutputFile(filename string) {
	//filename = fmt.Sprintf("output_data/%s", filename)

	outputString := fmt.Sprintf("%s\n", od.toString())

	if filename != "" {
		f, err := os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			//fmt.Println("Here")
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(outputString); err != nil {
			//fmt.Println("Here 2")
			log.Println(err)
		}
	}
}

func (od OutputData) SaveHeaderToOutputFile(filename string) {
	//filename = fmt.Sprintf("output_data/%s", filename)
	str := fmt.Sprintf("Name,InitialStateFile,TweakedParameter,Value,MinQualityDelta," +
		"MaxQualityDelta,AverageQualityDelta,MinTime,MaxTime,AverageTime,TestRuns")

	outputString := fmt.Sprintf("%s\n", str)

	if filename != "" {
		f, err := os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			//fmt.Println("Here")
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(outputString); err != nil {
			//fmt.Println("Here 2")
			log.Println(err)
		}
	}
}

type OutputList []OutputData

func (ol OutputList) toString() string {
	str := fmt.Sprintf("Name,InitialStateFile,TweakedParameter,Value,MinQualityDelta," +
		"MaxQualityDelta,AverageQualityDelta,MinTime,MaxTime,AverageTime,TestRuns")
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
			//fmt.Println("Here")
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(outputString); err != nil {
			//fmt.Println("Here 2")
			log.Println(err)
		}
	}
}

func testSetup() {
	// Import the country quality weights map
	importQualityWeights(countryWeightsFile)

}

const (
	countryName              = "Gondor"
	resourcesFile            = "../inputs/resources.csv"
	countryWeightsFile       = "../inputs/country_quality_weights.csv"
	balancedCountriesFile    = "../inputs/countries_balanced.csv"
	imbalancedCountriesFile  = "../inputs/countries_imbalanced.csv"
	randomCountriesFile      = "../inputs/countries_random.csv"
	outputScheduleFile       = ""
	proposedScheduleFilename = ""
	numOutputSchedules       = 3
	depthBound               = 8
	frontierMaxSize          = 800
	roundsToSimulate         = 10
	beamCount                = 3

	testIterations = 10

	testMaxDepth           = 10
	testMaxFrontierSize    = 1000
	testMaxGamma           = 1
	testMinFailureConstant = -3
	testMaxBeamCount       = 5
)

//
// Max Depth Tests
//

// TestTweakMaxDepthWithBalancedCountries
func TestTweakMaxDepthWithBalancedCountries(t *testing.T) {
	testSetup()
	outputFilename := "../test_outputs/tweak_max_depth_with_balanced_countries.csv"
	MaxIteration := testIterations

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

	output := OutputData{}
	output.SaveHeaderToOutputFile(outputFilename)

	for maxDepth := 3; maxDepth <= testMaxDepth; maxDepth++ {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		finalQualityDelta := .0
		var minQualityDelta float64
		var maxQualityDelta float64

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, maxDepth, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for max depth: %v\n", i, maxDepth)
		}
		fmt.Printf("Finished all iterations for max depth: %v\n", maxDepth)

		output = OutputData{
			Name:                "TweakMaxDepthWithBalancedCountries",
			InitialStateFile:    strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:    "MaxDepth",
			Value:               float64(maxDepth),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}

		output.SaveToOutputFile(outputFilename)
	}
}

// TestTweakMaxDepthWithImbalancedCountries
func TestTweakMaxDepthWithImbalancedCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for maxDepth := 3; maxDepth <= testMaxDepth; maxDepth++ {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, maxDepth, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for max depth: %v\n", i, maxDepth)
		}
		fmt.Printf("Finished all iterations for max depth: %v\n", maxDepth)

		output := OutputData{
			Name:                "TweakMaxDepthWithImbalancedCountries",
			InitialStateFile:    strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:    "MaxDepth",
			Value:               float64(maxDepth),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_depth_with_imbalanced_countries.csv")
}

// TestTweakMaxDepthWithRandomCountries
func TestTweakMaxDepthWithRandomCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for maxDepth := 3; maxDepth <= testMaxDepth; maxDepth++ {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, maxDepth, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for max depth: %v\n", i, maxDepth)
		}
		fmt.Printf("Finished all iterations for max depth: %v\n", maxDepth)

		output := OutputData{
			Name:                "TweakMaxDepthWithRandomCountries",
			InitialStateFile:    strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:    "MaxDepth",
			Value:               float64(maxDepth),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
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
	testSetup()
	MaxIteration := testIterations

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

	for maxFrontierSize := 100; maxFrontierSize <= testMaxFrontierSize; maxFrontierSize += 100 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, maxFrontierSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for max frontier size: %v\n", i, maxFrontierSize)
		}
		fmt.Printf("Finished all iterations max frontier size: %v\n", maxFrontierSize)

		output := OutputData{
			Name:                "TweakMaxFrontierSizeWithBalancedCountries",
			InitialStateFile:    strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:    "FrontierMaxSize",
			Value:               float64(maxFrontierSize),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_frontier_size_with_balanced_countries.csv")
}

// TestTweakFrontierMaxSizeWithImbalancedCountries
func TestTweakFrontierMaxSizeWithImbalancedCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for maxFrontierSize := 100; maxFrontierSize <= testMaxFrontierSize; maxFrontierSize += 100 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, maxFrontierSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for max frontier size: %v\n", i, maxFrontierSize)
		}
		fmt.Printf("Finished all iterations max frontier size: %v\n", maxFrontierSize)

		output := OutputData{
			Name:                "TweakFrontierMaxSizeWithImbalancedCountries",
			InitialStateFile:    strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:    "FrontierMaxSize",
			Value:               float64(maxFrontierSize),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_max_frontier_size_with_imbalanced_countries.csv")
}

// TestTweakFrontierMaxSizeWithRandomCountries
func TestTweakFrontierMaxSizeWithRandomCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for maxFrontierSize := 100; maxFrontierSize <= testMaxFrontierSize; maxFrontierSize += 100 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, maxFrontierSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for max frontier size: %v\n", i, maxFrontierSize)
		}
		fmt.Printf("Finished all iterations max frontier size: %v\n", maxFrontierSize)

		output := OutputData{
			Name:                "TweakFrontierMaxSizeWithRandomCountries",
			InitialStateFile:    strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:    "FrontierMaxSize",
			Value:               float64(maxFrontierSize),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
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
	testSetup()
	MaxIteration := testIterations

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

	for gamma := 0.4; gamma < testMaxGamma; gamma += 0.1 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		testConstants.TweakableConstants.Gamma = gamma

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for gamma: %v\n", i, gamma)
		}
		fmt.Printf("Finished all iterations for gamma: %v\n", gamma)

		output := OutputData{
			Name:                "TweakGammaWithBalancedCountries",
			InitialStateFile:    strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:    "Gamma",
			Value:               gamma,
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_gamma_with_balanced_countries.csv")
}

// TestTweakGammaWithImbalancedCountries
func TestTweakGammaWithImbalancedCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for gamma := 0.4; gamma < testMaxGamma; gamma += 0.1 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		testConstants.TweakableConstants.Gamma = gamma

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for gamma: %v\n", i, gamma)
		}
		fmt.Printf("Finished all iterations for gamma: %v\n", gamma)

		output := OutputData{
			Name:                "TweakGammaWithImbalancedCountries",
			InitialStateFile:    strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:    "Gamma",
			Value:               gamma,
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_gamma_with_imbalanced_countries.csv")
}

// TestTweakGammaWithRandomCountries
func TestTweakGammaWithRandomCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for gamma := 0.4; gamma < testMaxGamma; gamma += 0.1 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		testConstants.TweakableConstants.Gamma = gamma

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for gamma: %v\n", i, gamma)
		}
		fmt.Printf("Finished all iterations for gamma: %v\n", gamma)

		output := OutputData{
			Name:                "TweakGammaWithRandomCountries",
			InitialStateFile:    strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:    "Gamma",
			Value:               gamma,
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
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
	testSetup()
	MaxIteration := testIterations

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			//FailureConstant: util.FailureCost,
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

	for failureConstant := -0.5; failureConstant >= testMinFailureConstant; failureConstant -= 0.5 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		testConstants.TweakableConstants.FailureConstant = failureConstant

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for failure constant: %v\n", i, failureConstant)
		}
		fmt.Printf("Finished all iterations for failure constant: %v\n", failureConstant)

		output := OutputData{
			Name:                "TweakFailureConstantWithBalancedCountries",
			InitialStateFile:    strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:    "FailureCost",
			Value:               failureConstant,
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_failure_constant_with_balanced_countries.csv")
}

// TestTweakFailureConstantWithImbalancedCountries
func TestTweakFailureConstantWithImbalancedCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			//FailureConstant: util.FailureCost,
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

	for failureConstant := -0.5; failureConstant >= testMinFailureConstant; failureConstant -= 0.5 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		testConstants.TweakableConstants.FailureConstant = failureConstant

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for failure constant: %v\n", i, failureConstant)
		}
		fmt.Printf("Finished all iterations for failure constant: %v\n", failureConstant)

		output := OutputData{
			Name:                "TweakFailureConstantWithImbalancedCountries",
			InitialStateFile:    strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:    "FailureCost",
			Value:               failureConstant,
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_failure_constant_with_imbalanced_countries.csv")
}

// TestTweakFailureConstantWithRandomCountries
func TestTweakFailureConstantWithRandomCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

	testConstants := Constants{
		TweakableConstants: TweakableConstants{
			//FailureConstant: util.FailureCost,
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

	for failureConstant := -0.5; failureConstant >= testMinFailureConstant; failureConstant -= 0.5 {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		testConstants.TweakableConstants.FailureConstant = failureConstant

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for failure constant: %v\n", i, failureConstant)
		}
		fmt.Printf("Finished all iterations for failure constant: %v\n", failureConstant)

		output := OutputData{
			Name:                "TweakFailureConstantWithRandomCountries",
			InitialStateFile:    strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:    "FailureCost",
			Value:               failureConstant,
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_failure_constant_with_random_countries.csv")
}

//
// Number of Beams Tests
//

// TestTweakBeamCountWithBalancedCountries
func TestTweakBeamCountWithBalancedCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for testBeamCount := 1; testBeamCount < testMaxBeamCount; testBeamCount++ {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, testBeamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for beam count: %v\n", i, testBeamCount)
		}
		fmt.Printf("Finished all iterations for beam count: %v\n", testBeamCount)

		output := OutputData{
			Name:                "TweakBeamCountWithBalancedCountries",
			InitialStateFile:    strings.Trim(balancedCountriesFile, "../"),
			TweakedParameter:    "BeamCount",
			Value:               float64(testBeamCount),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_beam_count_with_balanced_countries.csv")
}

// TestTweakBeamCountWithImbalancedCountries
func TestTweakBeamCountWithImbalancedCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for testBeamCount := 1; testBeamCount < testMaxBeamCount; testBeamCount++ {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for beam count: %v\n", i, testBeamCount)
		}
		fmt.Printf("Finished all iterations for beam count: %v\n", testBeamCount)

		output := OutputData{
			Name:                "TweakBeamCountWithImbalancedCountries",
			InitialStateFile:    strings.Trim(imbalancedCountriesFile, "../"),
			TweakedParameter:    "BeamCount",
			Value:               float64(testBeamCount),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_beam_count_with_imbalanced_countries.csv")
}

// TestTweakBeamCountWithRandomCountries
func TestTweakBeamCountWithRandomCountries(t *testing.T) {
	testSetup()
	MaxIteration := testIterations

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

	for testBeamCount := 1; testBeamCount < testMaxBeamCount; testBeamCount++ {
		var totalTime time.Duration
		var minTime time.Duration
		var maxTime time.Duration

		var finalQualityDelta float64
		var minQualityDelta float64
		var maxQualityDelta float64

		finalQualityDelta = .0

		for i := 0; i < MaxIteration; i++ {
			res := allCountrySchedulersWithGameManager(resourcesFile, balancedCountriesFile, outputScheduleFile,
				proposedScheduleFilename, numOutputSchedules, depthBound, frontierMaxSize, roundsToSimulate, beamCount,
				testConstants)

			totalTime += res.Time

			finalQualityDelta += res.AverageQualityDelta

			if i == 0 {
				minTime = res.Time
				maxTime = res.Time

				minQualityDelta = res.AverageQualityDelta
				maxQualityDelta = res.AverageQualityDelta
			} else {
				if res.Time < minTime {
					minTime = res.Time
				} else if res.Time > maxTime {
					maxTime = res.Time
				}

				if res.AverageQualityDelta < minQualityDelta {
					minQualityDelta = res.AverageQualityDelta
				} else if res.AverageQualityDelta > maxQualityDelta {
					maxQualityDelta = res.AverageQualityDelta
				}
			}

			fmt.Printf("Finished iteration %v for beam count: %v\n", i, testBeamCount)
		}
		fmt.Printf("Finished all iterations for beam count: %v\n", testBeamCount)

		output := OutputData{
			Name:                "TweakBeamCountWithRandomCountries",
			InitialStateFile:    strings.Trim(randomCountriesFile, "../"),
			TweakedParameter:    "BeamCount",
			Value:               float64(testBeamCount),
			MinTime:             minTime.Seconds(),
			MaxTime:             maxTime.Seconds(),
			AverageTime:         totalTime.Seconds() / float64(MaxIteration),
			TestRuns:            MaxIteration,
			MinQualityDelta:     minQualityDelta,
			MaxQualityDelta:     maxQualityDelta,
			AverageQualityDelta: finalQualityDelta / float64(MaxIteration),
		}
		outputList = append(outputList, output)
	}

	outputList.SaveToOutputFile("../test_outputs/tweak_beam_count_with_random_countries.csv")
}
