package transfers

import (
	"CS5260_Final_Project/scheduler"
	"CS5260_Final_Project/util"
	"fmt"
)

// Transfer represents a resource transfer between two countries
type Transfer struct {
	Name     string
	Type     string
	Resource scheduler.ActionResource
}

// Take will perform the transformation on the provided country
func (t Transfer) Take(country scheduler.Country, otherCountry scheduler.Country) error {
	//fmt.Printf("Action Taken: %s\n", t.Name)
	return country.Transfer(otherCountry, t.Resource)
}

// GetType returns the type of this action
func (t Transfer) GetType() string {
	return t.Type
}

// GetName returns the name of this action
func (t Transfer) GetName() string {
	return t.Name
}

// InitializeTransfers will initialize the map of Transfers
func InitializeTransfers() map[string]scheduler.ScheduleAction {
	transferMap := make(map[string]scheduler.ScheduleAction)

	// MetallicElements Template
	transferMap["MetallicElements Transfer"] = Transfer{
		Name: "MetallicElements Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "MetallicElements",
			Amount: util.TransferMult * 10,
		},
	}

	// Timber Template
	transferMap["Timber Transfer"] = Transfer{
		Name: "Timber Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "Timber",
			Amount: util.TransferMult * 10,
		},
	}

	// MetallicAlloys Template
	transferMap["MetallicAlloys Transfer"] = Transfer{
		Name: "MetallicAlloys Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "MetallicAlloys",
			Amount: util.TransferMult * 10,
		},
	}

	// MetallicAlloysWaste Template
	transferMap["MetallicAlloysWaste Transfer"] = Transfer{
		Name: "MetallicAlloysWaste Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "MetallicAlloysWaste",
			Amount: util.TransferMult * 10,
		},
	}

	// Electronics Template
	transferMap["Electronics Transfer"] = Transfer{
		Name: "Electronics Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "Electronics",
			Amount: util.TransferMult * 10,
		},
	}

	// ElectronicsWaste Template
	transferMap["ElectronicsWaste Transfer"] = Transfer{
		Name: "ElectronicsWaste Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "ElectronicsWaste",
			Amount: util.TransferMult * 10,
		},
	}

	// Food Template
	transferMap["Food Transfer"] = Transfer{
		Name: "Food Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "Food",
			Amount: util.TransferMult * 10,
		},
	}

	// FoodWaste Template
	transferMap["FoodWaste Transfer"] = Transfer{
		Name: "FoodWaste Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "FoodWaste",
			Amount: util.TransferMult * 10,
		},
	}

	return transferMap
}

func (t Transfer) ToString(country scheduler.Country, otherCountry scheduler.Country) string {
	c1 := fmt.Sprintf("\"%s\"", country.GetName())
	c2 := fmt.Sprintf("\"%s\"", otherCountry.GetName())

	if otherCountry.GetSelf() {
		c2 = fmt.Sprintf("\"%s\"", otherCountry.GetName())
		c1 = fmt.Sprintf("\"%s\"", country.GetName())
	}

	str := fmt.Sprintf("(TRANSFER %s %s ((%s %v))", c1, c2, t.Resource.Name, t.Resource.Amount)

	str = fmt.Sprintf("%s))", str)

	return str
}
