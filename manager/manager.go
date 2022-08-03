package manager

import (
	"CS5260_Final_Project/scheduler"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// GameManager will manage the countries in the simulation and the execution of schedules
type GameManager struct {
	// scheduleProposalChans is a map of country name to channel that country will send its schedule proposals on
	scheduleProposalChans map[string]chan Schedule

	// scheduleSendChans is a map of country name to channel that country will receive a schedule on
	scheduleSendChans map[string]chan Schedule

	// terminateChan will receive a signal to stop the game manager
	terminateChan chan interface{}

	// outputScheduleFilename is the filename of the schedule output
	outputScheduleFilename string

	//proposedScheduleFilename is the filename of all the proposed schedule outputs]
	proposedScheduleFilename string
}

// Schedule is a schedule that a country wants to execute
type Schedule struct {
	Actions         []scheduler.ScheduledAction
	State           SimulatedWorldState
	ProposedCountry SchedulingCountry
}

type SimulatedWorldState interface {
	// GenerateProbabilityOfSuccess will generate the probability of a state being accepted by all countries
	GenerateProbabilityOfSuccess(int)

	// CalculateExpectedUtility will generate the ExpectedUtility values for each country in a SimulatedWorldState
	CalculateExpectedUtility(int)

	// ToString will convert the Simulated world state to a string
	ToString() string

	// SetAsNewInitialState will set the state values to that of an initial state
	SetAsNewInitialState()

	// StringifyCountries will turn the country resource maps into strings
	StringifyCountries() string

	// GetExpectedUtility will return the ExpectedUtility of the world state
	GetExpectedUtility() float64

	// GetExpectedUtilitySum will return the sum of all expected utilities from the initial state to this state
	GetExpectedUtilitySum() float64

	// GetFinalCountryQualitiesDelta will return the final delta of the quality ratings of each country
	GetFinalCountryQualitiesDelta() map[string]float64
}

// SchedulingCountry is a country that has a channel to send schedule proposals on a channel on which to receive a
// schedule to execute
type SchedulingCountry interface {
	// ProposeSchedule will send a schedule proposal on an interface and block until a schedule to execute is received
	ProposeSchedule([]scheduler.ScheduledAction, SimulatedWorldState) ([]scheduler.ScheduledAction, SimulatedWorldState, error)

	// RegisterRecvChan will register the game manager send channel with the country
	RegisterRecvChan(chan Schedule)

	// GetName will return the name of the country
	GetName() string
}

// man contains the singleton GameManager
var man GameManager

var r *rand.Rand

// InitializeGameManager will be used to create the game manager singleton
func InitializeGameManager(terminateChan chan interface{}, outputScheduleFilename, proposedScheduleFilename string) GameManager {
	man = GameManager{
		scheduleProposalChans:    make(map[string]chan Schedule),
		scheduleSendChans:        make(map[string]chan Schedule),
		terminateChan:            terminateChan,
		outputScheduleFilename:   outputScheduleFilename,
		proposedScheduleFilename: proposedScheduleFilename,
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	return man
}

// GetGameManager will return the singleton game manager
func GetGameManager() GameManager {
	return man
}

// Register will register a schedule channel with a game manager
func (gm *GameManager) Register(country SchedulingCountry, proposalChan chan Schedule) {
	// Add proposal channel to map
	gm.scheduleProposalChans[country.GetName()] = proposalChan

	// Create channel to send official schedule on and add to map
	sendChan := make(chan Schedule)
	gm.scheduleSendChans[country.GetName()] = sendChan

	// Register the newly created channel to the country
	country.RegisterRecvChan(sendChan)
	// TODO: Add this function to countries.country
}

// Run will run the game manager
func (gm *GameManager) Run(rounds int) float64 {
	time.Sleep(5 * time.Second)

	var averageQualityDelta float64

	// Infinite loop while running, but will close whenever main terminates
	for round := 1; round <= rounds; round++ {
		fmt.Printf("Running the game manager for round %v out of %v\n", round, rounds)
		time.Sleep(500 * time.Millisecond)

		if round == rounds {
			_ = round
		}

		// Keep track of the sum of the delta from start to finish of the quality metrics, and reset each round
		finalQualityDeltaSum := .0

		// Check for game termination
		if gm.terminateChan != nil {
			select {
			case _, ok := <-gm.terminateChan:
				if ok {
					break
				} else {
					fmt.Println("Channel closed")
				}
				//break
			default:
				fmt.Println("Not yet terminated")
			}
		}

		proposedSchedules := make([]Schedule, len(gm.scheduleProposalChans))
		proposals := 0

		proposalRoundInformation := fmt.Sprintf("***********************************************************************************************\n")
		proposalRoundInformation = fmt.Sprintf("%s*************************************  Round %v Proposals *************************************\n",
			proposalRoundInformation, round)
		gm.printToFile(gm.proposedScheduleFilename, proposalRoundInformation)

		roundInformation := fmt.Sprintf("***********************************************************************************************\n")
		roundInformation = fmt.Sprintf("%s******************************************  Round %v  ******************************************\n",
			roundInformation, round)
		//roundInformation = fmt.Sprintf("%s\nDepth bound this round: %v\nFrontier size: %v\nTop "+
		//	"schedules count: %v\nSchedules:\n\n", roundInformation, depthBound-1, frontierMaxSize, numOutputSchedules)
		gm.printToFile(gm.outputScheduleFilename, roundInformation)

		// Listen on each registered channel until a message is received from each country
		for _, propChan := range gm.scheduleProposalChans {
			propSchedule := <-propChan
			proposedSchedules[proposals] = propSchedule
			proposals++

			// Add the quality delta of each country
			for _, qualityRating := range propSchedule.State.GetFinalCountryQualitiesDelta() {
				finalQualityDeltaSum += qualityRating
			}

			// print the proposal to the file
			gm.printToFile(gm.proposedScheduleFilename, fmt.Sprintf("Proposed Country:\t%s\n", propSchedule.ProposedCountry.GetName()))
			gm.printToFile(gm.proposedScheduleFilename, propSchedule.State.ToString())
		}

		// Get the average expected utility across each country
		averageQualityDelta = finalQualityDeltaSum / float64(len(proposedSchedules))

		// Decide on best schedule
		bestSchedule := gm.GetBestSchedule(proposedSchedules)

		// Print best schedule information to the file
		scheduleInformation := fmt.Sprintf("Chose schedule from:\t\t    %s\nExpected Utility for %s:\t    %v\n",
			bestSchedule.ProposedCountry.GetName(), bestSchedule.ProposedCountry.GetName(), bestSchedule.State.GetExpectedUtility())
		gm.printToFile(gm.outputScheduleFilename, fmt.Sprintf("%s\n", scheduleInformation))
		gm.printToFile(gm.outputScheduleFilename, bestSchedule.State.ToString())
		gm.printToFile(gm.outputScheduleFilename, fmt.Sprintf("%v\n\n\n", bestSchedule.State.StringifyCountries()))

		// Send the best schedule to each of the countries
		for _, sendChan := range gm.scheduleSendChans {
			sendChan <- bestSchedule
		}

		// We will now continue on with the next iteration of the for loop and block when listening for new proposals
	}

	return averageQualityDelta
}

// GetBestSchedule will return the "Best" schedule, which will be randomly assigned for now. Other algorithms could be
// used to perform better "Best" optimizations. For now though, I am using a random algorithm to ensure each country
// will have an even distribution of weights when proposing schedules.
func (gm *GameManager) GetBestSchedule(proposedSchedules []Schedule) Schedule {
	return proposedSchedules[r.Intn(len(proposedSchedules))]
}

// printToFile will print the proposed schedules to the output file
func (gm *GameManager) printToFile(filename, data string) {
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
