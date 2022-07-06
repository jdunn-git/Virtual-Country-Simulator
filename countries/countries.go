package countries

import (
	"CS5260_Final_Project/resources"
	"CS5260_Final_Project/scheduler"
	"fmt"
)

// Country struct is used to store all information about a country
type Country struct {
	Name         string
	ResourcesMap map[string]*resources.CountryResource
	Quality      float64
	Self         bool
}

// NewCountry will construct a new Country pointer
func NewCountry(name string, countryResources []*resources.CountryResource, self bool) *Country {
	c := &Country{
		Name:         name,
		ResourcesMap: make(map[string]*resources.CountryResource),
		Self:         self,
	}

	for _, res := range countryResources {
		if res == nil {
			panic("No resource found")
		}
		c.ResourcesMap[res.GetName()] = res
	}

	return c
}

// GetSelf can be used to get the "self" country
func GetSelf(countriesMap map[string]*Country) (*Country, error) {
	for _, country := range countriesMap {
		if country.Self {
			return country, nil
		}
	}

	return &Country{}, fmt.Errorf("there was no \"self\" country")
}

// GetName will return the Name of this country
func (c *Country) GetName() string {
	return c.Name
}

// GetQualityRating will return the Quality of this country
func (c *Country) GetQualityRating() float64 {
	return c.Quality
}

// SetQualityRating will set the Quality of this country
func (c *Country) SetQualityRating(quality float64) {
	c.Quality = quality
}

// GetResourceAmount will return the amount this country has of the passed in resource
func (c *Country) GetResourceAmount(resourceName string) int {
	return c.ResourcesMap[resourceName].GetAmount()
}

// GetResourceWeight will return the weight of the passed in resource as a float64
func (c *Country) GetResourceWeight(resourceName string) float64 {
	return float64(c.ResourcesMap[resourceName].GetWeight())
}

// GetResourceAmountAndWeight will return the amount and weight of the passed in resource as a float64
func (c *Country) GetResourceAmountAndWeight(resourceName string) (float64, float64) {
	return float64(c.GetResourceAmount(resourceName)), c.GetResourceWeight(resourceName)
}

// Print will print the Country
func (c *Country) Print() {
	fmt.Printf("Country name: %s\n\tQuality: %v\n", c.Name, c.Quality)
	for _, res := range c.ResourcesMap {
		fmt.Printf("\t%s\n", res.String())
	}
}

// AdjustResource will increase or decrease the resource by the amount provided
func (c *Country) AdjustResource(resourceName string, adjustment int) {
	c.ResourcesMap[resourceName].SetAmount(c.ResourcesMap[resourceName].GetAmount() + adjustment)
}

// Duplicate will return an identical copy to this Country
func (c *Country) Duplicate() *Country {
	duplicatedResourcesMap := make(map[string]*resources.CountryResource)

	for key, resource := range c.ResourcesMap {
		duplicatedResourcesMap[key] = resource.Duplicate()
	}

	return &Country{
		Name:         c.Name,
		ResourcesMap: duplicatedResourcesMap,
		Quality:      c.Quality,
	}
}

// Transform will transform this country's resources according to the inputs and outputs
func (c *Country) Transform(inputs []scheduler.ActionResource,
	outputs []scheduler.ActionResource,
	transformationName string) error {
	//fmt.Printf("Performing %s on %s\n", template.Name, c.Name)

	// Validate that there is enough of each resource
	for _, transformation := range inputs {
		if c.GetResourceAmount(transformation.Name) < transformation.Amount {
			return fmt.Errorf("%s did not have enough %s to perform %s. Had %v, Needs %v",
				c.Name, transformation.Name, transformationName, c.GetResourceAmount(transformation.Name), transformation.Amount)
		}
	}

	// Perform Transformation Inputs
	for _, transformation := range inputs {
		// Since this is an input to the transformation, multiply the amount -1 to reduce the country's resources
		c.AdjustResource(transformation.Name, -1*transformation.Amount)
	}
	// Perform Transformation Outputs
	for _, transformation := range outputs {
		c.AdjustResource(transformation.Name, transformation.Amount)
	}

	return nil
}

// Transfer will transfer a resource between two countries
func (c *Country) Transfer(otherCountry scheduler.Country, resource scheduler.ActionResource) error {
	//fmt.Printf("Tranferring (%v) %s from %s to %s\n", amount, resourceName, c.Name, otherCountry.Name)

	// Validate that there is enough of the resource
	if c.GetResourceAmount(resource.Name) < resource.Amount {
		return fmt.Errorf("%s did not have enough %s to transfer for %s. Had %v, Needs %v",
			c.Name, resource.Name, otherCountry.GetName(), c.GetResourceAmount(resource.Name), resource.Amount)
	}

	// Multiple amount by negative one so that we reduce the amount for this country
	c.AdjustResource(resource.Name, -1*resource.Amount)

	otherCountry.AdjustResource(resource.Name, resource.Amount)

	return nil
}
