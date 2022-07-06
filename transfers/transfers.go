package transfers

import (
	"CS5260_Final_Project/scheduler"
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

// InitializeTransfers will initialize the map of Transfers
func InitializeTransfers() map[string]scheduler.ScheduleAction {
	transferMap := make(map[string]scheduler.ScheduleAction)

	// MetallicElements Template
	transferMap["MetallicElements Transfer"] = Transfer{
		Name: "MetallicElements Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "MetallicElements",
			Amount: 10,
		},
	}

	// Timber Template
	transferMap["Timber Transfer"] = Transfer{
		Name: "Timber Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "Timber",
			Amount: 10,
		},
	}

	// MetallicAlloys Template
	transferMap["MetallicAlloys Transfer"] = Transfer{
		Name: "MetallicAlloys Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "MetallicAlloys",
			Amount: 10,
		},
	}

	// MetallicAlloysWaste Template
	transferMap["MetallicAlloysWaste Transfer"] = Transfer{
		Name: "MetallicAlloysWaste Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "MetallicAlloysWaste",
			Amount: 10,
		},
	}

	// Electronics Template
	transferMap["Electronics Transfer"] = Transfer{
		Name: "Electronics Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "Electronics",
			Amount: 10,
		},
	}

	// ElectronicsWaste Template
	transferMap["ElectronicsWaste Transfer"] = Transfer{
		Name: "ElectronicsWaste Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "ElectronicsWaste",
			Amount: 10,
		},
	}

	// Food Template
	transferMap["Food Transfer"] = Transfer{
		Name: "Food Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "Food",
			Amount: 10,
		},
	}

	// FoodWaste Template
	transferMap["FoodWaste Transfer"] = Transfer{
		Name: "FoodWaste Transfer",
		Type: scheduler.TransferType,
		Resource: scheduler.ActionResource{
			Name:   "FoodWaste",
			Amount: 10,
		},
	}

	return transferMap
}
