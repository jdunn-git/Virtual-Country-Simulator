package transformations

import (
	"CS5260_Final_Project/scheduler"
	"CS5260_Final_Project/util"
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
	Type    string
	Inputs  []scheduler.ActionResource
	Outputs []scheduler.ActionResource
}

// Take will perform the transformation on the provided country
func (t Transformation) Take(country scheduler.Country, otherCountry scheduler.Country) error {
	//fmt.Printf("Action Taken: %s\n", t.Name)

	// The interface requires two countries passed in, but only one is needed for the Transformation Action
	if otherCountry != nil {
		return fmt.Errorf("two countries passed into Transformation, but only one is needed")
	}

	return country.Transform(t.Inputs, t.Outputs, t.Name)
}

// GetType returns the type of this action
func (t Transformation) GetType() string {
	return t.Type
}

// GetName returns the name of this action
func (t Transformation) GetName() string {
	return t.Name
}

// InitializeTransformations will initialize the map of Transformations
func InitializeTransformations() map[string]scheduler.ScheduleAction {
	transformationMap := make(map[string]scheduler.ScheduleAction)

	// Housing Template
	transformationMap["Housing Transformation"] = Transformation{
		Name: "Housing Transformation",
		Type: scheduler.TransformType,
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 5,
			},
			{
				Name:   "MetallicElements",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Timber",
				Amount: util.TransformMult * 5,
			},
			{
				Name:   "MetallicAlloys",
				Amount: util.TransformMult * 3,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Housing",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "HousingWaste",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Population",
				Amount: util.TransformMult * 5,
			},
		},
	}

	// Alloys Template
	transformationMap["Alloys Transformation"] = Transformation{
		Name: "Alloys Transformation",
		Type: scheduler.TransformType,
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "MetallicElements",
				Amount: util.TransformMult * 2,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "MetallicAlloys",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "MetallicAlloysWaste",
				Amount: util.TransformMult * 1,
			},
		},
	}

	// Electronics Template
	transformationMap["Electronics Transformation"] = Transformation{
		Name: "Electronics Transformation",
		Type: scheduler.TransformType,
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "MetallicElements",
				Amount: util.TransformMult * 3,
			},
			{
				Name:   "MetallicAlloys",
				Amount: util.TransformMult * 2,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Electronics",
				Amount: util.TransformMult * 2,
			},
			{
				Name:   "ElectronicsWaste",
				Amount: util.TransformMult * 1,
			},
		},
	}

	// Farm Template
	transformationMap["Farm Transformation"] = Transformation{
		Name: "Farm Transformation",
		Type: scheduler.TransformType,
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "AvailableLand",
				Amount: util.TransformMult * 5,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Farm",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "FarmWaste",
				Amount: util.TransformMult * 1,
			},
		},
	}

	// Food Template
	transformationMap["Food Transformation"] = Transformation{
		Name: "Food Transformation",
		Type: scheduler.TransformType,
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Farm",
				Amount: util.TransformMult * 3,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Food",
				Amount: util.TransformMult * 2,
			},
			{
				Name:   "FoodWaste",
				Amount: util.TransformMult * 1,
			},
		},
	}

	// Military Template
	transformationMap["Military Transformation"] = Transformation{
		Name: "Military Transformation",
		Type: scheduler.TransformType,
		Inputs: []scheduler.ActionResource{
			{
				Name:   "Population",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "Housing",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "MetallicAlloys",
				Amount: util.TransformMult * 2,
			},
			{
				Name:   "Food",
				Amount: util.TransformMult * 3,
			},
		},
		Outputs: []scheduler.ActionResource{
			{
				Name:   "Military",
				Amount: util.TransformMult * 1,
			},
			{
				Name:   "MilitaryWaste",
				Amount: util.TransformMult * 1,
			},
		},
	}

	return transformationMap
}

func (t Transformation) ToString(country scheduler.Country, otherCountry scheduler.Country) string {
	str := fmt.Sprintf("(TRANSFORM \"%s\"\n\t(INPUTS  ", country.GetName())

	for i, input := range t.Inputs {
		str = fmt.Sprintf("%s(%s %v)", str, input.Name, input.Amount)

		if i < len(t.Inputs)-1 {
			str = fmt.Sprintf("%s\n\t         ", str)
		}
	}

	str = fmt.Sprintf("%s\n\t(OUTPUTS ", str)

	for i, output := range t.Outputs {
		str = fmt.Sprintf("%s(%s %v)", str, output.Name, output.Amount)

		if i < len(t.Outputs)-1 {
			str = fmt.Sprintf("%s\n\t         ", str)
		}
	}

	str = fmt.Sprintf("%s))", str)

	return str
}
