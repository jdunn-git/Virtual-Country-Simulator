package scheduler

var TransformType string
var TransferType string

var FailureConstant = -0.05

type Country interface {
	// Transform will transform a set of input resources into a set of output resources for this country
	Transform(inputs []ActionResource,
		outputs []ActionResource,
		transformation string) error

	// Transfer will move resources from this country to otherCountry
	Transfer(otherCountry Country, resource ActionResource) error

	// GetQualityRating will return the Quality of the country
	GetQualityRating() float64
	// SetQualityRating will set the Quality of the country
	SetQualityRating(float64)

	// GetDiscountedReward will return the DiscountedReward of the country
	GetDiscountedReward() float64
	// SetDiscountedReward will set the DiscountedReward of the country
	SetDiscountedReward(float64)

	// GetUndiscountedReward will return the DiscountedReward of the country
	GetUndiscountedReward() float64
	// SetUndiscountedReward will set the DiscountedReward of the country
	SetUndiscountedReward(float64)

	// GetName will return the Name of the country
	GetName() string

	// Self will return true if this country is self
	GetSelf() bool

	// AdjustResource will make a change the amount of the resource with resourceName by the adjustment amount
	AdjustResource(resourceName string, adjustment int)

	// Duplicate will duplicate this country for a simulation
	Duplicate() Country

	// GetResourceAmount will return this country's amount of the passed in resource
	GetResourceAmount(resourceName string) int
	// GetResourceAmountAndWeight will return this country's amount and weight of the passed in resource
	GetResourceAmountAndWeight(resourceName string) (float64, float64)
}

// ActionResource will store a resource name and amount for a given action
type ActionResource struct {
	Name   string
	Amount int
}

// ScheduleAction represents an action that can be scheduled and taken
type ScheduleAction interface {
	// Take will cause country to "take" the action, using otherCountry if needed
	Take(country Country, otherCountry Country) error

	// GetType will return the type constant for this action
	GetType() string

	// GetName will return the name of this action
	GetName() string

	// ToString will convert the action to a string representation
	ToString(country Country, otherCountry Country) string
}

// ScheduledAction represents an action scheduled for one or two countries
type ScheduledAction struct {
	Action       ScheduleAction
	ThisCountry  Country
	OtherCountry Country
}

// ToString will convert the ScheduledAction to a string format
func (sa *ScheduledAction) ToString() string {
	return sa.Action.ToString(sa.ThisCountry, sa.OtherCountry)
}

// CountryScheduler contains the countries map for this scheduler instance
type CountryScheduler struct {
	// CountriesMap contains all countries that can be used in the schedule
	CountriesMap map[string]Country
}

// AvailableActions is the global map of all actions that a country can take
var AvailableActions map[string]ScheduleAction

// InitializeAvailableActions will initialize all transformations and transfers
func InitializeAvailableActions(transformations map[string]ScheduleAction, transfers map[string]ScheduleAction) {
	AvailableActions = make(map[string]ScheduleAction)

	for key, val := range transformations {
		AvailableActions[key] = val
	}

	for key, val := range transfers {
		AvailableActions[key] = val
	}
}
