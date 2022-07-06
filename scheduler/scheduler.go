package scheduler

import "fmt"

type Country interface {
	Transform(inputs []ActionResource,
		outputs []ActionResource,
		transformation string) error

	Transfer(otherCountry Country, resource ActionResource) error

	GetName() string

	AdjustResource(resourceName string, adjustment int)
}

type ActionResource struct {
	Name   string
	Amount int
}

// ScheduleAction represents an action that can be scheduled and taken
type ScheduleAction interface {
	Take(country Country, otherCountry Country) error
}

var availableActions map[string]ScheduleAction

// InitializeAvailableActions will intialize all transformations and transfers
func InitializeAvailableActions(transformations map[string]ScheduleAction, transfers map[string]ScheduleAction) {
	availableActions = make(map[string]ScheduleAction)

	for key, val := range transformations {
		availableActions[key] = val
	}

	for key, val := range transfers {
		availableActions[key] = val
	}
}

func PerformAction(country Country, otherCountry Country, action string) error {
	if availableActions[action] == nil {
		return fmt.Errorf("%s not found in list of available actions: %v", action, availableActions)
	}
	return availableActions[action].Take(country, otherCountry)
}
