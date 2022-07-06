package transfers

import (
	"CS5260_Final_Project/scheduler"
)

// Transfer represents a resource transfer between two countries
type Transfer struct {
	Name     string
	Resource scheduler.ActionResource
}

// Take will perform the transformation on the provided country
func (t Transfer) Take(country scheduler.Country, otherCountry scheduler.Country) error {
	return country.Transfer(otherCountry, t.Resource)
}

// InitializeTransfers will initialize the map of Transfers
func InitializeTransfers() map[string]scheduler.ScheduleAction {
	transferMap := make(map[string]scheduler.ScheduleAction)

	// MetallicElements Template
	transferMap["MetallicElements Transfer"] = Transfer{
		Name: "MetallicElements Transfer",
		Resource: scheduler.ActionResource{
			Name:   "MetallicElements",
			Amount: 10,
		},
	}

	// Timber Template
	transferMap["Timber Transfer"] = Transfer{
		Name: "Timber Transfer",
		Resource: scheduler.ActionResource{
			Name:   "Timber",
			Amount: 10,
		},
	}

	// MetallicAlloys Template
	transferMap["MetallicAlloys Transfer"] = Transfer{
		Name: "MetallicAlloys Transfer",
		Resource: scheduler.ActionResource{
			Name:   "MetallicAlloys",
			Amount: 10,
		},
	}

	// MetallicAlloysWaste Template
	transferMap["MetallicAlloysWaste Transfer"] = Transfer{
		Name: "MetallicAlloysWaste Transfer",
		Resource: scheduler.ActionResource{
			Name:   "MetallicAlloysWaste",
			Amount: 10,
		},
	}

	// Electronics Template
	transferMap["Electronics Transfer"] = Transfer{
		Name: "Electronics Transfer",
		Resource: scheduler.ActionResource{
			Name:   "Electronics",
			Amount: 10,
		},
	}

	// ElectronicsWaste Template
	transferMap["ElectronicsWaste Transfer"] = Transfer{
		Name: "ElectronicsWaste Transfer",
		Resource: scheduler.ActionResource{
			Name:   "ElectronicsWaste",
			Amount: 10,
		},
	}

	// Food Template
	transferMap["Food Transfer"] = Transfer{
		Name: "Food Transfer",
		Resource: scheduler.ActionResource{
			Name:   "Food",
			Amount: 10,
		},
	}

	// FoodWaste Template
	transferMap["FoodWaste Transfer"] = Transfer{
		Name: "FoodWaste Transfer",
		Resource: scheduler.ActionResource{
			Name:   "FoodWaste",
			Amount: 10,
		},
	}

	return transferMap
}
