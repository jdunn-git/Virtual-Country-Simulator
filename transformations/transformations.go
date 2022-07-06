package transformations

import (
	"CS5260_Final_Project/scheduler"
	"fmt"
)

// TransformationResource is a resource being input or output from a transformation
type TransformationResource struct {
	Name   string
	Amount int
}

// Transformation is a resource being input or output from a transformation
type Transformation struct {
	Name    string
	Inputs  []scheduler.ActionResource
	Outputs []scheduler.ActionResource
}

// Take will perform the transformation on the provided country
func (t Transformation) Take(country scheduler.Country, otherCountry scheduler.Country) error {
	// The interface requires two countries passed in, but only one is needed for the Transformation Action
	if otherCountry != nil {
		return fmt.Errorf("two countries passed into Transformation, but only one is needed")
	}

	return country.Transform(t.Inputs, t.Outputs, t.Name)
}

// InitializeTransformations will initialize the map of Transformations
func InitializeTransformations() map[string]scheduler.ScheduleAction {
	transformationMap := make(map[string]scheduler.ScheduleAction)

	// Housing Template
	transformationMap["Housing Transformation"] = Transformation{
		Name: "Housing Transformation",
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 5,
			},
			{
				Name:   "MetallicElements",
				Amount: 1,
			},
			{
				Name:   "Timber",
				Amount: 5,
			},
			{
				Name:   "MetallicAlloys",
				Amount: 3,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Housing",
				Amount: 1,
			},
			{
				Name:   "HousingWaste",
				Amount: 1,
			},
			{
				Name:   "Population",
				Amount: 5,
			},
		},
	}

	// Alloys Template
	transformationMap["Alloys Transformation"] = Transformation{
		Name: "Alloys Transformation",
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "MetallicElements",
				Amount: 2,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "MetallicAlloys",
				Amount: 1,
			},
			{
				Name:   "MetallicAlloysWaste",
				Amount: 1,
			},
		},
	}

	// Electronics Template
	transformationMap["Electronics Transformation"] = Transformation{
		Name: "Electronics Transformation",
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "MetallicElements",
				Amount: 3,
			},
			{
				Name:   "MetallicAlloys",
				Amount: 2,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "Electronics",
				Amount: 2,
			},
			{
				Name:   "ElectronicsWaste",
				Amount: 1,
			},
		},
	}

	// Farm Template
	transformationMap["Farm Transformation"] = Transformation{
		Name: "Farm Transformation",
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "AvailableLand",
				Amount: 5,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "Farm",
				Amount: 1,
			},
			{
				Name:   "FarmWaste",
				Amount: 1,
			},
		},
	}

	// Food Template
	transformationMap["Food Transformation"] = Transformation{
		Name: "Food Transformation",
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "Farm",
				Amount: 3,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "Food",
				Amount: 2,
			},
			{
				Name:   "FoodWaste",
				Amount: 1,
			},
		},
	}

	// Military Template
	transformationMap["Military Transformation"] = Transformation{
		Name: "Military Transformation",
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: 1,
			},
			{
				Name:   "Housing",
				Amount: 1,
			},
			{
				Name:   "MetallicAlloys",
				Amount: 2,
			},
			{
				Name:   "Food",
				Amount: 3,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Military",
				Amount: 1,
			},
			{
				Name:   "MilitaryWaste",
				Amount: 1,
			},
		},
	}

	return transformationMap
}
