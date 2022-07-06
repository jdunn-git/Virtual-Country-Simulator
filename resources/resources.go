package resources

import "fmt"

const (
	raw     string = "Raw"
	derived        = "Derived"
)

// Resource is a struct to maintain all information of the resource definition
type Resource struct {
	Name         string
	Weight       float32
	ResourceType string
	Renewable    bool
}

// NewResource is a constructor for Resource
func NewResource(name string, weight float32, renewable bool) *Resource {
	rawResources := []string{"Population", "MetallicElements", "Timber", "AvailableLand"}

	resourceType := derived
	for _, resName := range rawResources {
		if resName == name {
			resourceType = raw
			break
		}
	}

	return &Resource{
		Name:         name,
		Weight:       weight,
		ResourceType: resourceType,
		Renewable:    renewable,
	}
}

// Print will print the Resource
func (r *Resource) Print() {
	fmt.Printf("Resource: %s\n\tWeight: %v\n\tResourceType: %s\n\tand Renewable: %v\n", r.Name, r.Weight, r.ResourceType, r.Renewable)
}

// CountryResource will maintin informtaion about an individual country's resource
type CountryResource struct {
	Resource *Resource
	Amount   int
}

// NewCountryResource is a constructor for Resource
func NewCountryResource(resource *Resource, amount int) *CountryResource {
	return &CountryResource{
		Resource: resource,
		Amount:   amount,
	}
}

// Print will print the CountryResource
func (cr *CountryResource) Print() {
	fmt.Println(cr.String())
}

// String returns a string representation of CountryResource
func (cr *CountryResource) String() string {
	return fmt.Sprintf("Resource Name: %s, Amount: %v", cr.GetName(), cr.Amount)
}

// GetName returns the name of this resource
func (cr *CountryResource) GetName() string {
	return cr.Resource.Name
}

// GetAmount returns the amount of this resource
func (cr *CountryResource) GetAmount() int {
	return cr.Amount
}

// SetAmount will set a new amount for the resource
func (cr *CountryResource) SetAmount(newAmount int) {
	cr.Amount = newAmount
}

// GetWeight returns the amount of this resource
func (cr *CountryResource) GetWeight() float32 {
	return cr.Resource.Weight
}

// Duplicate will return an identical copy to this CountryResource
func (cr *CountryResource) Duplicate() *CountryResource {
	return &CountryResource{
		Resource: cr.Resource,
		Amount:   cr.Amount,
	}
}
