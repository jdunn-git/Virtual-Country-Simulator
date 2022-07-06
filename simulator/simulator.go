package simulator

import (
	"CS5260_Final_Project/countries"
	"CS5260_Final_Project/quality"
	"CS5260_Final_Project/scheduler"
	"container/heap"
	"fmt"
	"math/rand"
)

// Simulator can be used to simulate actions against simulated worlds and maintains the states in a
//  priority queue
type Simulator struct {
	States   WorldStates // This is the Priority Queue for the simulator
	MaxSize  int
	MaxDepth int
}

// InitializeSimulator will initialize the simulator to be able to process any number of simulated worlds
//  in the priority queue
func InitializeSimulator(countriesMap map[string]*countries.Country, maxSize int, maxDepth int) (*Simulator, *WorldState) {
	return &Simulator{
			MaxSize:  maxSize,
			MaxDepth: maxDepth,
		}, &WorldState{
			SimulatedCountries: countriesMap,
			Generation:         0,
		}
}

// SimulateAllActionsFromState will add every action on to the priority queue
func (s *Simulator) SimulateAllActionsFromState(actionMap map[string]scheduler.ScheduleAction,
	countryMap map[string]*countries.Country,
	startingState *WorldState,
	depth int) {

	// Break conditions
	if depth > s.MaxDepth || len(s.States) >= s.MaxSize {
		return
	}

	// Make an action list we can safely remove from
	actions := make([]scheduler.ScheduleAction, len(actionMap))
	i := 0
	for _, val := range actionMap {
		actions[i] = val
		i++
	}

	// Iterate over all actions, and remove each action after performing it
	for len(actions) > 0 && len(s.States) < s.MaxSize {
		// This will hold the new world state after the simulation
		var newWorldState *WorldState

		// Randomly grab an action
		randomIndex := rand.Intn(len(actions))
		action := actions[randomIndex]

		// Remove action for list of actions
		newActions := make([]scheduler.ScheduleAction, len(actions)-1)
		copy(newActions[:randomIndex], actions[:randomIndex])
		copy(newActions[randomIndex:], actions[randomIndex+1:])
		actions = newActions

		// Randomly get two countries, with one being self
		var country *countries.Country
		var otherCountry *countries.Country
		var err error

		country, err = countries.GetSelf(countryMap)
		if err != nil {
			panic("No self country")
		}

		// Only need two countries when there's a Transfer
		if action.GetType() == scheduler.TransferType {
			// This approach isn't fast but it promotes the most random approach within reason
			countryList := make([]*countries.Country, len(countryMap)-1)
			i = 0
			for _, c := range countryMap {
				// Skip the self country
				if c.Name != country.Name {
					countryList[i] = c
					i++
				}
			}
			countryRand := rand.Intn(len(countryList))
			otherCountry = countryList[countryRand]

			// Perform one more random check to see whether the self country should be the primary or secondary country
			selfRand := rand.Intn(2)
			if selfRand%2 == 0 {
				tempCountry := country
				country = otherCountry
				otherCountry = tempCountry
			}

			// Simulate action if possible, adding to the priority queue
			newWorldState = s.SimulateActionIfPossible(startingState, country, otherCountry, action)
		} else {
			// Simulate action if possible, adding to the priority queue
			newWorldState = s.SimulateActionIfPossible(startingState, country, nil, action)
		}
		// Next, recursively call SimulateAllActionsFromState so that we populate the priority queue in a depth-first manner
		s.SimulateAllActionsFromState(actionMap, countryMap, newWorldState, depth+1)
	}
}

// SimulateActionIfPossible will perform an action against a simulated (duplicated) copy of the world
//  but will not return an error if the action was not possible
func (s *Simulator) SimulateActionIfPossible(ws *WorldState,
	country scheduler.Country,
	otherCountry scheduler.Country,
	action scheduler.ScheduleAction) *WorldState {

	newWorldState, err := s.SimulateAction(ws, country, otherCountry, action)
	if err != nil {
		fmt.Println(err)
	}

	return newWorldState
}

// SimulateAction will perform an action against a simulated (duplicated) copy of the world
func (s *Simulator) SimulateAction(ws *WorldState,
	country scheduler.Country,
	otherCountry scheduler.Country,
	action scheduler.ScheduleAction) (*WorldState, error) {

	// Duplicate the state of the world
	sim := make(map[string]*countries.Country)

	for key, country := range ws.SimulatedCountries {
		sim[key] = country.Duplicate()
	}

	newCountry := sim[country.GetName()]
	var newOtherCountry scheduler.Country
	if otherCountry != nil {
		newOtherCountry = sim[otherCountry.GetName()]
	}

	// Perform action on duplicated world
	err := action.Take(newCountry, newOtherCountry)
	if err != nil {
		return &WorldState{}, fmt.Errorf("error while simulating action: %v", err)
	}

	// Retrieve the Quality Scores
	// TODO: Enhance this by having the correct score. For now, just use the gondor quality rating
	for _, country := range sim {
		country.SetQualityRating(quality.PerformQualityCalculation(country))
	}

	// Push WorldState to priority queue
	newWorldState := &WorldState{
		SimulatedCountries: sim,
		PreviousState:      ws,
		WorldValue:         sim["Gondor"].GetQualityRating(),
		Generation:         ws.Generation + 1,
		Index:              0,
	}

	//fmt.Printf("Adding World Value %v\n", newWorldState.WorldValue)
	// Init heap if this is the first time that the priority queue is being pushed to
	if s.States == nil {
		states := WorldStates{newWorldState}
		//states[0] = newWorldState
		heap.Init(&states)
		s.States = states
		//fmt.Println("First Run!!")
	} else {
		newWorldStates := make(WorldStates, len(s.States)+1)
		//fmt.Printf("*** Adding %v\n", newWorldState.WorldValue)
		copy(newWorldStates, s.States)
		newWorldStates[len(s.States)] = newWorldState
		heap.Push(&newWorldStates, newWorldState)
		//heap.Fix(&newWorldStates, len(s.States))
		s.States = newWorldStates
		//for _, state := range s.States {
		//	fmt.Printf("\tIndex: %v, Value: %v\n", state.Index, state.WorldValue)
		//}
		//time.Sleep(500 * time.Millisecond)
	}

	return newWorldState, nil
}

//func (s *Simulator) PrintWorldStates() {
//	for _, worldState := range s.States {
//		fmt.Println("***** Gondor resources *****")
//		worldState.SimulatedCountries["Gondor"].Print()
//		fmt.Println("***** Rohan resources *****")
//		worldState.SimulatedCountries["Rohan"].Print()
//		fmt.Printf("generation: %v\nWorldValue: %v\n\n", worldState.Generation, worldState.WorldValue)
//	}
//}

// TODO: Will also need a way to dump the entire priority queue and rebuild
//  - Maybe rebuild based off of the biggest jumps in value?

// GetTopNLeaves will return the top leaves from the priority queue, as well as all previous nodes of these leaves
//  We will be refreshing the priority queue at this point so popping is okay
func (s *Simulator) GetTopNLeaves(num int, startingGeneration int) ([]*WorldState, map[string]*WorldState) {
	topStates := make([]*WorldState, num)
	retStates := make(map[string]*WorldState)
	//states := s.States

	i := 0
	for i < num {
		// Pop Top State
		state := heap.Pop(&s.States).(*WorldState)

		// Verify that it's a leaf
		if state.Generation == s.MaxDepth {
			// Add leaf to the retStates list
			topStates[i] = state
			retStates[fmt.Sprintf("%v - %v", state.Index, state.WorldValue)] = state
			i++
		}
	}

	// Get all generations back from these states for when we rebuild the priority queue
	for _, state := range topStates {
		i = state.Generation
		for i > startingGeneration {
			i--
			saveState := state.GetGeneration(i)
			retStates[fmt.Sprintf("%v - %v", saveState.Index, saveState.WorldValue)] = saveState
		}
	}

	return topStates, retStates
}

// FlushQueue will empty out the priority queue of all states
func (s *Simulator) FlushQueue() {
	s.States = make([]*WorldState, 0)
	heap.Init(&s.States)
}

// RebuildQueue will rebuild the priority queue after a refresh
func (s *Simulator) RebuildQueue(topStates []*WorldState) {
	s.States = topStates
	heap.Init(&s.States)
}

// WorldState contains a state of all countries in the world
type WorldState struct {
	SimulatedCountries map[string]*countries.Country
	PreviousState      *WorldState
	//NextState          *WorldState
	WorldValue float64
	Generation int
	Index      int
}

// GetGeneration returns the PreviousState with num Generation value
func (ws *WorldState) GetGeneration(num int) *WorldState {
	state := *ws
	gen := state.Generation
	for gen > num {
		state = *state.PreviousState
		gen = state.Generation
	}

	return &state
}

// WorldStates is the type needed for the priority queue
type WorldStates []*WorldState

// Len is implemented for the sort interface
func (ws WorldStates) Len() int {
	return len(ws)
}

// Swap implementation is for the sort interface
func (ws WorldStates) Swap(i, j int) {
	//fmt.Printf("Swapping %v and %v\n", ws[i].WorldValue, ws[j].WorldValue)
	ws[i], ws[j] = ws[j], ws[i]
	ws[i].Index = i
	ws[j].Index = j
}

// Less implementation is for the sort interface
func (ws WorldStates) Less(i, j int) bool {
	//fmt.Printf("Comparing %v >= %v: result is %v\n", ws[i].WorldValue, ws[j].WorldValue, ws[i].WorldValue >= ws[j].WorldValue)
	return ws[i].WorldValue >= ws[j].WorldValue
}

// Push implementation is for the priority queue / heap interface
func (ws *WorldStates) Push(worldState any) {
	n := len(*ws) - 1
	newWorldState := worldState.(*WorldState)
	newWorldState.Index = n
	//newWorldStates := make(WorldStates, len(*ws)+1)
	//copy(newWorldStates, *ws)
	//newWorldStates[len(newWorldStates)-1] = &newWorldState
	//*ws = newWorldStates
	(*ws)[len(*ws)-1] = newWorldState
}

// Pop implementation is for the priority queue / heap interface
func (ws *WorldStates) Pop() any {
	tmp := *ws
	l := len(tmp)
	res := tmp[l-1]
	tmp[l-1] = nil
	//res.Index = -1
	*ws = tmp[0 : l-1]
	return res
}
