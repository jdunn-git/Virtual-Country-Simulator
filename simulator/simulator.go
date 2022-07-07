package simulator

import (
	"CS5260_Final_Project/countries"
	"CS5260_Final_Project/quality"
	"CS5260_Final_Project/scheduler"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	X0 float64
	K  float64
	L  float64
)

// Simulator can be used to simulate actions against simulated worlds and maintains the states in a
//  priority queue
type Simulator struct {
	States       WorldStates // This is the Priority Queue for the simulator
	MaxSize      int
	MaxDepth     int
	InitialDepth int
	Randomizer   *rand.Rand
}

// InitializeSimulator will initialize the simulator to be able to process any number of simulated worlds
//  in the priority queue
func InitializeSimulator(worldStates WorldStates, maxSize, maxDepth, initialDepth int) *Simulator {
	s := &Simulator{
		MaxSize:      maxSize,
		MaxDepth:     maxDepth,
		States:       worldStates,
		InitialDepth: initialDepth,
		Randomizer:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	s.GenerateIntialRewards()

	heap.Init(&s.States)

	return s
}

// GenerateIntialRewards will set the rewards for the initial world state when the simulator is initialized
func (s *Simulator) GenerateIntialRewards() {
	for _, state := range s.States {
		if state.PreviousState != nil {
			for _, country := range state.SimulatedCountries {
				// Only need to change the Discounted reward since this is dependent on the depth
				country.SetDiscountedReward(quality.CalculateDiscountedReward(country,
					s.MaxDepth-state.Generation))
			}
		}
	}
}

// SimulateAllActionsFromState will add every action on to the priority queue
func (s *Simulator) SimulateAllActionsFromState(actionMap map[string]scheduler.ScheduleAction,
	countryMap map[string]scheduler.Country,
	startingState *WorldState,
	depth int) bool {

	// Break conditions
	if depth > s.MaxDepth || len(s.States) >= s.MaxSize {
		return len(s.States) >= s.MaxSize
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
		var country scheduler.Country
		var otherCountry scheduler.Country
		var err error

		country, err = countries.GetSelf(countryMap)
		if err != nil {
			panic("No self country")
		}

		var actionError error
		// Only need two countries when there's a Transfer
		if action.GetType() == scheduler.TransferType {
			// This approach isn't fast but it promotes the most random approach within reason
			countryList := make([]scheduler.Country, len(countryMap)-1)
			i = 0
			for _, c := range countryMap {
				// Skip the self country
				if c.GetName() != country.GetName() {
					countryList[i] = c
					i++
				}
			}
			countryRand := s.Randomizer.Intn(len(countryList))
			otherCountry = countryList[countryRand]

			// Perform one more random check to see whether the self country should be the primary or secondary country
			selfRand := s.Randomizer.Intn(2)
			if selfRand%2 == 0 {
				tempCountry := country
				country = otherCountry
				otherCountry = tempCountry
			}

			// Only perform an action if the action does not already exist in the priority queue
			actionId := fmt.Sprintf("%v | %v-%v", startingState.Id, action.GetName(), country.GetName())
			if !s.States.CheckIfActionExists(actionId) {
				// Simulate action if possible, adding to the priority queue
				newWorldState, actionError = s.SimulateAction(startingState, country, otherCountry, action)
				if actionError != nil {
					_ = actionError
				}
			}
		} else {
			// Only perform an action if the action does not already exist in the priority queue
			actionId := fmt.Sprintf("%v | %v-%v", startingState.Id, action.GetName(), country.GetName())
			if !s.States.CheckIfActionExists(actionId) {
				// Simulate action if possible, adding to the priority queue
				newWorldState, actionError = s.SimulateAction(startingState, country, nil, action)
				if actionError != nil {
					_ = actionError
				}
			}
		}

		if actionError == nil {
			// Next, recursively call SimulateAllActionsFromState so that we populate the priority queue in a depth-first manner
			s.SimulateAllActionsFromState(actionMap, countryMap, newWorldState, depth+1)
		} else {
			//fmt.Println(actionError)
		}
	}

	return len(s.States) >= s.MaxSize
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
	sim := make(map[string]scheduler.Country)

	for key, c := range ws.SimulatedCountries {
		sim[key] = c.Duplicate()
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
	for _, countrySim := range sim {
		countrySim.SetQualityRating(quality.PerformQualityCalculation(countrySim))
		countrySim.SetUndiscountedReward(quality.CalculateUndiscountedReward(countrySim))
		countrySim.SetDiscountedReward(quality.CalculateDiscountedReward(countrySim,
			s.MaxDepth-ws.Generation+1))
	}

	// This ID is used for uniqueness, but it's only concern with the previous state, the action, and the country
	// performing the action. If there's a second country for a transfer, then it doesn't make a significant enough
	// difference to warrant being considered a "unique" state for world quality
	//     - This is because country A transferring to country B is almost identical in the calculation to country A
	//		 transferring the same resource+amount to country C
	ID := fmt.Sprintf("%v | %v-%v", ws.Id, action.GetName(), country.GetName())

	// Push WorldState to priority queue
	newWorldState := &WorldState{
		SimulatedCountries: sim,
		PreviousState:      ws,
		Generation:         ws.Generation + 1,
		Index:              0,
		Id:                 ID,
		Action: scheduler.ScheduledAction{
			Action:       action,
			ThisCountry:  country,
			OtherCountry: otherCountry,
		},
	}

	newWorldState.CalculateExpectedUtility(s.InitialDepth)

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

// GetTopNLeaves will return the top leaves from the priority queue, as well as all previous nodes of these leaves
//  We will be refreshing the priority queue at this point so popping is okay
func (s *Simulator) GetTopNLeaves(num int, startingGeneration int) (map[int]*WorldState, map[int]WorldStates, int) {
	topStates := make(map[int]*WorldState)
	backupStates := make(map[int]WorldStates)
	//states := s.States

	topLevel := make(WorldStates, num)

	// Maintain a count of states that we're saving
	count := 0

	i := 0
	for i < num && len(s.States) > 0 {
		// Pop Top State
		state := heap.Pop(&s.States).(*WorldState)

		// Verify that it's a leaf
		if state.Generation == s.MaxDepth {
			// Add leaf to the retStates list
			topStates[i] = state

			// Also add this to the priority queue for the max depth
			topLevel[i] = state

			i++
			count++
		}
	}

	// Error condition that all states might pop out before enough top states are found
	if len(s.States) == 0 {
		topLevel = topLevel[:i]
	}

	heap.Init(&topLevel)
	backupStates[s.MaxDepth] = topLevel

	// Get all generations back from these states for when we rebuild the priority queue
	for _, state := range topStates {
		i = state.Generation
		for i > startingGeneration {
			i--
			saveState := state.GetGeneration(i)

			// Add this state to the backupStates
			// If the backupStates for this generation isn't initialize, then init
			if _, found := backupStates[saveState.Generation]; !found {
				newWorldStates := make(WorldStates, 1)
				newWorldStates[0] = saveState
				heap.Init(&newWorldStates)
				backupStates[saveState.Generation] = newWorldStates
				count++
			} else {
				// Check if state already exists in the WorldStates
				if !backupStates[saveState.Generation].CheckIfActionExists(state.Id) {
					newWorldStates := make(WorldStates, len(backupStates[saveState.Generation])+1)
					copy(newWorldStates, backupStates[saveState.Generation])
					newWorldStates[len(newWorldStates)-1] = saveState
					heap.Push(&newWorldStates, saveState)
					backupStates[saveState.Generation] = newWorldStates
					count++
				}
			}

		}
	}

	return topStates, backupStates, count
}

// FlushQueue will empty out the priority queue of all states
func (s *Simulator) FlushQueue() {
	for len(s.States) > 0 {
		heap.Pop(&s.States)
	}
}

// InitializePriorityQueue will rebuild the priority queue after a refresh
func (s *Simulator) InitializePriorityQueue(topStates []*WorldState) {
	s.States = topStates
	heap.Init(&s.States)
}

// ConvertToWorldStates will return a WorldStates pointer based on a countryMap passed in
func ConvertToWorldStates(countryMap map[string]scheduler.Country) WorldStates {
	// Duplicate the state of the world
	sim := make(map[string]scheduler.Country)

	for key, country := range countryMap {
		sim[key] = country.Duplicate()
	}

	ID := "Initial State"

	// Push WorldState to priority queue
	newWorldState := &WorldState{
		SimulatedCountries: sim,
		PreviousState:      nil,
		Generation:         0,
		Index:              -1,
		Id:                 ID,
		ExpectedUtility:    0,
	}

	return []*WorldState{newWorldState}
}

// WorldState contains a state of all countries in the world, and represents a potential schedule
type WorldState struct {
	SimulatedCountries map[string]scheduler.Country
	PreviousState      *WorldState
	Action             scheduler.ScheduledAction
	WorldValue         float64
	SuccessProbability float64
	ExpectedUtility    float64
	Generation         int
	Index              int
	Id                 string
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

//// GenerateWorldQualityRating will set the WorldValue for a WorldState
//func (ws *WorldState) GenerateWorldQualityRating() float64 {
//	// TODO: Calculate Undiscounted Reward
//
//	// TODO: Calculate Discounted Reward
//
//	// TODO: Calculate Success Probability
//
//	// TODO: Calculate Expected Utility
//
//	ws.WorldValue = ws.SimulatedCountries["Gondor"].GetQualityRating()
//	return ws.WorldValue
//}

// GenerateProbabilityOfSuccess will return a number between 0 and 1 with the success likelihood of this schedule
func (ws *WorldState) GenerateProbabilityOfSuccess(originalGeneration int) {
	successOdds := 1.0

	// Get each country involved in the schedule
	countryNames := make(map[string]interface{})
	states := ws
	for states.Generation >= originalGeneration {
		if states.Generation != 0 {
			break
		}

		countryNames[states.Action.ThisCountry.GetName()] = 0
		if states.Action.OtherCountry != nil {
			countryNames[states.Action.OtherCountry.GetName()] = 0
		}

		states = states.PreviousState
	}

	logisticRegressions := make([]float64, len(countryNames))

	i := 0
	// Iterate over each country
	for name, _ := range countryNames {

		// Perform a Logistic function to generate success probability
		lr := -K * (ws.SimulatedCountries[name].GetDiscountedReward() - X0)
		lr = 1 + math.Exp(lr)

		lr = L / lr

		logisticRegressions[i] = lr
		i++
	}

	// Get the product of the individual probabilities
	for _, lr := range logisticRegressions {
		successOdds *= lr
	}

	ws.SuccessProbability = successOdds
}

// CalculateExpectedUtility will calculate the ExpectedUtility of this schedule
func (ws *WorldState) CalculateExpectedUtility(originalGeneration int) {
	ws.GenerateProbabilityOfSuccess(originalGeneration)

	// Expected Utility is the Product of the success probability and the discounted reward of the self
	// country, summed together with the cost of the schedule failing

	// Get the self country's discounted reward
	selfCountry, err := countries.GetSelf(ws.SimulatedCountries)
	if err != nil {
		panic(err)
	}

	discountedReward := selfCountry.GetDiscountedReward()

	// Multiply by the success probability
	successValue := discountedReward * ws.SuccessProbability

	// Calculate the cost of failure
	failureValue := (1 - ws.SuccessProbability) * scheduler.FailureConstant

	// Add the two values
	ws.ExpectedUtility = successValue + failureValue

}

func (ws *WorldState) ToString() string {
	str := ""
	state := ws
	// TODO need to work my way up here
	for state.Generation > 0 {
		if state.Generation == 1 {
			str = fmt.Sprintf(" %s EU: %v\n%s", state.Action.ToString(), state.ExpectedUtility, str)
		} else {
			str = fmt.Sprintf("  %s EU: %v\n%s", state.Action.ToString(), state.ExpectedUtility, str)
		}
		state = state.PreviousState
	}
	str = fmt.Sprintf("[%s\n]\n\n", str)

	return str
}

// WorldStates is the type needed for the priority queue
type WorldStates []*WorldState

// CheckIfActionExists will return true if the action is already found in the WorldStates list
func (ws WorldStates) CheckIfActionExists(Id string) bool {
	for _, state := range ws {
		if Id == state.Id {
			return true
		}
	}

	return false
}

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
	return ws[i].ExpectedUtility >= ws[j].ExpectedUtility
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
