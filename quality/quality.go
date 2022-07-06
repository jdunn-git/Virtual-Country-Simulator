package quality

import (
	"CS5260_Final_Project/countries"
	"fmt"
	"math"
)

// PerformQualityCalculation will generate a quality rating for a country
func PerformQualityCalculation(country *countries.Country) float64 {

	// Get the raw materials
	population := float64(country.GetResourceAmount("Population"))
	//metallicElements := country.GetResourceAmount("MetallicElements")
	//timber := country.GetResourceAmount("Timber")
	availableLand := float64(country.GetResourceAmount("AvailableLand"))

	// Get the derived materials and weights
	metallicAlloys, metallicAlloysWeight := country.GetResourceAmountAndWeight("MetallicAlloys")
	metallicAlloysWaste, metallicAlloysWasteWeight := country.GetResourceAmountAndWeight("MetallicAlloysWaste")
	electronics, electronicsWeight := country.GetResourceAmountAndWeight("Electronics")
	electronicsWaste, electronicsWasteWeight := country.GetResourceAmountAndWeight("ElectronicsWaste")
	housing, housingWeight := country.GetResourceAmountAndWeight("Housing")
	housingWaste, housingWasteWeight := country.GetResourceAmountAndWeight("HousingWaste")
	farm, farmWeight := country.GetResourceAmountAndWeight("Farm")
	farmWaste, farmWasteWeight := country.GetResourceAmountAndWeight("FarmWaste")
	food, foodWeight := country.GetResourceAmountAndWeight("Food")
	foodWaste, foodWasteWeight := country.GetResourceAmountAndWeight("FoodWaste")
	military, militaryWeight := country.GetResourceAmountAndWeight("Military")
	militaryWaste, militaryWasteWeight := country.GetResourceAmountAndWeight("MilitaryWaste")

	// Partial Quality - Food
	foodQuality := ((food * foodWeight) + (foodWaste * foodWasteWeight)) / population

	// Partial Quality - Housing
	housingQuality := ((housing * housingWeight) + (housingWaste * housingWasteWeight)) / (population * 0.4)

	// Partial Quality - Electronics / Technology
	electronicsQuality := (electronics * electronicsWeight) + (electronicsWaste * electronicsWasteWeight)

	// Partial Quality - Metals Materials
	metalsQuality := (metallicAlloys * metallicAlloysWeight) + (metallicAlloysWaste * metallicAlloysWasteWeight)

	// Partial Quality - Land
	landQuality := (math.Pow((farm*farmWeight), 2) + (farmWaste * farmWasteWeight)) / availableLand

	// Partial Quality - Military
	militaryQuality := ((military * militaryWeight) + (militaryWaste * militaryWasteWeight)) / (population * 0.2)

	///*
	fmt.Printf("foodQuality: %v\n", foodQuality)
	fmt.Printf("housingQuality: %v\n", housingQuality)
	fmt.Printf("electronicsQuality: %v\n", electronicsQuality)
	fmt.Printf("metalsQuality: %v\n", metalsQuality)
	fmt.Printf("landQuality: %v\n", landQuality)
	fmt.Printf("militaryQuality: %v\n", militaryQuality)
	//*/

	// Final Quality Calculation
	return (foodQuality * 10) + (housingQuality * 3) + electronicsQuality + metalsQuality +
		landQuality + (militaryQuality * 5)
}
