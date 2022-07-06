package scheduler

import "fmt"

const (
	TransformType string = "Transform"
	TransferType         = "Transfer"
)

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
	GetType() string
}

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

func PerformAction(country Country, otherCountry Country, action string) error {
	if AvailableActions[action] == nil {
		return fmt.Errorf("%s not found in list of available actions: %v", action, AvailableActions)
	}
	return AvailableActions[action].Take(country, otherCountry)
}
