package simulator

import (
	"CS5260_Final_Project/countries"
	"CS5260_Final_Project/quality"
	"CS5260_Final_Project/scheduler"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sort"
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
	CountriesMap      map[string]scheduler.Country
	States            WorldStates // This is the Priority Queue for the simulator
	stateCountPerBeam map[int]int // This is the number of states in each beam
	ResourceNameList  []string    // This will hold a list of all resources
	MaxSize           int
	MaxDepth          int
	InitialDepth      int
	Randomizer        *rand.Rand
}

// InitializeSimulator will initialize the simulator to be able to process any number of simulated worlds
//  in the priority queue
func InitializeSimulator(countriesMap map[string]scheduler.Country, worldStates WorldStates, resourceNameList []string, maxSize, maxDepth, initialDepth int) *Simulator {
	s := &Simulator{
		CountriesMap:     countriesMap,
		ResourceNameList: resourceNameList,
		MaxSize:          maxSize,
		MaxDepth:         maxDepth,
		States:           worldStates,
		InitialDepth:     initialDepth,
		Randomizer:       rand.New(rand.NewSource(time.Now().UnixNano())),
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
	depth int,
	maxSize int,
	beams int,
	beamID int) bool {

	// Break conditions
	if depth > s.MaxDepth || len(s.States) >= s.MaxSize || s.stateCountPerBeam[beamID] > maxSize {
		//fmt.Printf("len(s.States) (%v) is >= maxSize (%v)\n", len(s.States), maxSize)
		return len(s.States) >= maxSize
	}

	// Make an action list we can safely remove from
	actions := make([]scheduler.ScheduleAction, len(actionMap))
	i := 0
	for _, val := range actionMap {
		actions[i] = val
		i++
	}

	tempNewWorldStates := make(WorldStates, 0)
	heap.Init(&tempNewWorldStates)

	// Iterate over all actions, and remove each action after performing it
	for len(actions) > 0 && len(s.States) < s.MaxSize && s.stateCountPerBeam[beamID] < maxSize {
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
					//panic(actionError)
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
					//panic(actionError)
				}
			}
		}

		if actionError == nil {
			tmp := tempNewWorldStates
			tempNewWorldStates = make(WorldStates, len(tempNewWorldStates)+1)
			copy(tempNewWorldStates, tmp)
			tempNewWorldStates[len(tempNewWorldStates)-1] = newWorldState
			heap.Push(&tempNewWorldStates, newWorldState)

			if s.stateCountPerBeam != nil {
				s.stateCountPerBeam[beamID]++
				if s.stateCountPerBeam[beamID] > maxSize {
					break
				}
			}
		}
	}

	// Next, recursively call SimulateAllActionsFromState with the best state so that we populate the priority queue
	// in a depth-first manner
	//var bestWorldState *WorldState
	if len(tempNewWorldStates) > 0 {
		// If we don't have enough states to split across all beams, then do the best we can
		if beams > len(tempNewWorldStates) {
			beams = len(tempNewWorldStates)
		}

		bestWorldState := heap.Pop(&tempNewWorldStates).(*WorldState)

		if beams == 1 {
			// This is hit in a beam scenario - increment the state count in this beam
			if beamID > -1 {

			}
			s.SimulateAllActionsFromState(actionMap, countryMap, bestWorldState, depth+1, maxSize, 1, beamID)
		} else {
			// If we are using a beam search, then the first time this will split the available frontier into equal lengths for
			// the beams
			// Note that this block of code will only be hit the first time the search is happening
			bestWorldStates := make([]*WorldState, beams)
			bestWorldStates[0] = bestWorldState

			//fmt.Printf("bestWorldState[0].ExpectedUtility: %v\n", bestWorldStates[0].ExpectedUtility)
			for i := 1; i < beams; i++ {
				// No need to check for length of tempNewWorldStates here since there's no scenario where this should not be
				// full, since this code block only runs the first layer of a search
				bestWorldStates[i] = heap.Pop(&tempNewWorldStates).(*WorldState)
				//fmt.Printf("bestWorldState[%v].ExpectedUtility: %v\n", i, bestWorldStates[i].ExpectedUtility)
			}

			// Split the frontier size for each beam
			frontierSizePerBeam := maxSize / beams

			s.stateCountPerBeam = make(map[int]int, beams)
			for i, bestState := range bestWorldStates {
				// Set the initial state count in this beam
				s.stateCountPerBeam[i] = 0
				s.SimulateAllActionsFromState(actionMap, countryMap, bestState, depth+1, frontierSizePerBeam, 1, i)
			}
		}
	}
	return len(s.States) >= maxSize
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
		countrySim.SetUndiscountedReward(quality.CalculateUndiscountedReward(s.CountriesMap, countrySim))
		countrySim.SetDiscountedReward(quality.CalculateDiscountedReward(countrySim,
			ws.Generation))
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
		ResourceNameList:   ws.ResourceNameList,
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

// GetTopLeaf will return the top leaf from the priority queue, or the best non-leaf if we didn't get to the final
//  layer. It will also return the generation that was returned from.
func (s *Simulator) GetTopLeaf() (*WorldState, int) {
	topStatePerGeneration := make(map[int]*WorldState)

	// Maintain a count of states that we're saving
	topGeneration := 0

	for len(s.States) > 0 {
		// Pop Top State
		state := heap.Pop(&s.States).(*WorldState)

		if topStatePerGeneration[state.Generation] == nil {
			topStatePerGeneration[state.Generation] = state
			if topGeneration < state.Generation {
				topGeneration = state.Generation
			}
		} else if topStatePerGeneration[state.Generation].ExpectedUtility < state.ExpectedUtility {
			topStatePerGeneration[state.Generation] = state
			if state.Generation > topGeneration {
				topGeneration = state.Generation
			}
		}

	}

	return topStatePerGeneration[topGeneration], topGeneration
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
func ConvertToWorldStates(countryMap map[string]scheduler.Country, resourceNameList []string) WorldStates {
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
		ResourceNameList:   resourceNameList,
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
	ResourceNameList   []string
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

// GetExpectedUtility will return the expected utility of this state
func (ws *WorldState) GetExpectedUtility() float64 {
	return ws.ExpectedUtility
}

// GetExpectedUtilitySum will return a sum of all expected utilities in this state path
func (ws *WorldState) GetExpectedUtilitySum() float64 {
	sum := .0

	//selfCountry, _ := countries.GetSelf(ws.SimulatedCountries)
	//tmpQualityRating := selfCountry.GetQualityRating()
	//tmpUndiscountedReward := selfCountry.GetUndiscountedReward()
	//tmpDiscountedReward := selfCountry.GetDiscountedReward()
	//fmt.Printf("\tPost-simulation self.Quality: %v, self.UndiscountedReward: %v, self.DiscountedReward: %v\n",
	//	tmpQualityRating, tmpUndiscountedReward, tmpDiscountedReward)

	//fmt.Printf("* Getting Expected Utility Sum\n")
	state := ws
	for state.Generation > -1 {
		sum += state.ExpectedUtility
		//fmt.Printf("\tGeneration: %v, Expected utility: %v\n", state.Generation, state.ExpectedUtility)
		if state.PreviousState == nil {
			break
		}

		state = state.PreviousState
	}
	//fmt.Printf("\tSum: %v\n", sum)

	return sum
}

// GetFinalCountryQualitiesDelta will return the final delta of the quality ratings of each country
func (ws *WorldState) GetFinalCountryQualitiesDelta() map[string]float64 {
	retMap := make(map[string]float64, len(ws.SimulatedCountries))

	// Add the final quality of each country into the map
	for countryName, country := range ws.SimulatedCountries {
		retMap[countryName] = country.GetQualityRating()
	}

	// Get the initial state
	state := ws
	for state.PreviousState != nil {
		state = state.PreviousState
	}

	// Subtract the initial state's quality of each country
	for countryName, country := range state.SimulatedCountries {
		retMap[countryName] -= country.GetQualityRating()
	}

	return retMap
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
		if states.Generation == 0 {
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

	// Part Two: Since the game manager is evenly distributing the accepted schedules, I am no longer going to consider
	// other countries for the success probability. Instead, it will always be 0.95 to represent some degree of it not
	// working out.
	// Countries are now fully allowed to be selfish and propose whatever schedule they want, without considering how it
	// will impact other countries
	//ws.SuccessProbability = successOdds
	ws.SuccessProbability = 0.95
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

	if math.IsNaN(successValue) {
		_ = successValue
	}

	// Add the two values
	ws.ExpectedUtility = successValue + failureValue

}

// SetAsNewInitialState will set the state values to that of an initial state, but the simulated countries remain
func (ws *WorldState) SetAsNewInitialState() {
	ws.Index = -1
	ws.Generation = 0
	ws.PreviousState = nil
	ws.ExpectedUtility = .0
	ws.Id = ""
	ws.SuccessProbability = 1.00
}

// GetAllActions will return an ordered list of all scheduled actions in this state
func (ws *WorldState) GetAllActions() []scheduler.ScheduledAction {
	retActions := make([]scheduler.ScheduledAction, ws.Generation)
	tmpState := ws

	for i := ws.Generation - 1; i >= 0; i-- {
		retActions[i] = tmpState.Action
		tmpState = tmpState.PreviousState
	}

	return retActions
}

// ToString will convert the state of the world into a string
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

// StringifyCountries will transform the country resource maps to a string
func (ws *WorldState) StringifyCountries() string {
	str := ""

	str = " Country Name ||"
	for i, resourceName := range ws.ResourceNameList {
		str = fmt.Sprintf("%s %s", str, resourceName)
		if i < len(ws.ResourceNameList) {
			str = fmt.Sprintf("%s |", str)
		}
	}

	tmp_str := ""
	for i := -5; i < len(str); i++ {
		_ = i
		tmp_str = fmt.Sprintf("%s_", tmp_str)
	}

	str = fmt.Sprintf("%s\n%s\n", str, tmp_str)

	orderedCountryNames := make([]string, len(ws.SimulatedCountries))
	i := 0
	for countryName, _ := range ws.SimulatedCountries {
		orderedCountryNames[i] = countryName
		i++
	}

	sort.Strings(orderedCountryNames)

	for _, countryName := range orderedCountryNames {
		resourceStr := ""
		country := ws.SimulatedCountries[countryName]
		for i, resourceName := range ws.ResourceNameList {
			resourceAmountStr := fmt.Sprintf("%v", country.GetResourceAmount(resourceName))

			// Generate pre- and post-fix whitespace based on length of string
			prefixLen := 0
			postfixLen := 0
			diff := len(resourceName) - len(resourceAmountStr) + 2
			if len(resourceName)%2 == 0 {
				if len(resourceAmountStr)%2 == 1 {
					prefixLen = diff / 2
					postfixLen = diff/2 + 1
				} else {
					prefixLen = diff / 2
					postfixLen = diff / 2
				}
			} else {
				if len(resourceAmountStr)%2 == 1 {
					prefixLen = diff / 2
					postfixLen = diff / 2
				} else {
					prefixLen = diff / 2
					postfixLen = diff/2 + 1
				}
			}

			prefix := ""
			postfix := ""
			for i := 0; i < prefixLen; i++ {
				prefix = fmt.Sprintf("%s ", prefix)
			}
			for i := 0; i < postfixLen; i++ {
				postfix = fmt.Sprintf("%s ", postfix)
			}

			if i < len(ws.ResourceNameList) {
				resourceAmountStr = fmt.Sprintf("%s%s%s|", prefix, resourceAmountStr, postfix)
			}

			//whiteSpaceCount := len(resourceName) - len(resourceAmountStr)
			//for j := 0; j < whiteSpaceCount; j++ {
			//	_ = j
			//	resourceStr = fmt.Sprintf("%s ", resourceStr)
			//}

			resourceStr = fmt.Sprintf("%s%s", resourceStr, resourceAmountStr)

		}

		// Generate country name string with tabs based on length
		countryNameStr := countryName
		if len(countryNameStr) > 7 {
			countryNameStr = fmt.Sprintf("%s\t", countryNameStr)
		} else {
			countryNameStr = fmt.Sprintf("%s\t\t", countryNameStr)
		}
		str = fmt.Sprintf("%s %s  ||%s\n", str, countryNameStr, resourceStr)
	}

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
