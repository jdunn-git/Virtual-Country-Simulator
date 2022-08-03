package quality

import (
	"CS5260_Final_Project/scheduler"
	"math"
)

// Gamma is used for decay of discounted rewards over times
var Gamma float64

// CountryQualityWeights are the weights each country declares for the quality metric categories
type CountryQualityWeights struct {
	FoodQualityWeight  float64
	HousingQuality     float64
	ElectronicsQuality float64
	MetalsQuality      float64
	LandQuality        float64
	MilitaryQuality    float64
}

// CountryQualityWeightsMap is a map of country name to its quality weights
var CountryQualityWeightsMap map[string]CountryQualityWeights

// PerformQualityCalculation will generate a quality rating for a country
func PerformQualityCalculation(country scheduler.Country) float64 {

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
	electronicsQuality := ((electronics * electronicsWeight) + (electronicsWaste * electronicsWasteWeight)) / population

	// Partial Quality - Metals Materials
	metalsQuality := ((metallicAlloys * metallicAlloysWeight) + (metallicAlloysWaste * metallicAlloysWasteWeight)) / population

	// Partial Quality - Land
	landQuality := (math.Pow(farm*farmWeight, 2) + (farmWaste * farmWasteWeight)) / (availableLand * population)

	// Partial Quality - Military
	militaryQuality := ((military * militaryWeight) + (militaryWaste * militaryWasteWeight)) / (population * 0.05)

	// Check for no remaining land or population
	if availableLand == 0 {
		landQuality = 0
	}
	if population == 0 {
		foodQuality = 0
		housingQuality = 0
		electronicsQuality = 0
		metalsQuality = 0
		landQuality = 0
		militaryQuality = 0
	}

	// Multiple each by 100 since they're all low typically - higher values make for a more interesting simulation
	foodQuality *= 100
	housingQuality *= 100
	electronicsQuality *= 100
	metalsQuality *= 100
	landQuality *= 100
	militaryQuality *= 100

	/*
		fmt.Printf("foodQuality: %v\n", foodQuality)
		fmt.Printf("housingQuality: %v\n", housingQuality)
		fmt.Printf("electronicsQuality: %v\n", electronicsQuality)
		fmt.Printf("metalsQuality: %v\n", metalsQuality)
		fmt.Printf("landQuality: %v\n", landQuality)
		fmt.Printf("militaryQuality: %v\n", militaryQuality)
	*/

	// Final Quality Calculation
	qualityValue := (foodQuality * CountryQualityWeightsMap[country.GetName()].FoodQualityWeight) +
		(housingQuality * CountryQualityWeightsMap[country.GetName()].HousingQuality) +
		(electronicsQuality * CountryQualityWeightsMap[country.GetName()].ElectronicsQuality) +
		(metalsQuality * CountryQualityWeightsMap[country.GetName()].MetalsQuality) +
		(landQuality * CountryQualityWeightsMap[country.GetName()].LandQuality) +
		(militaryQuality * CountryQualityWeightsMap[country.GetName()].MilitaryQuality)
	if math.IsInf(qualityValue, 1) {
		_ = qualityValue
	}
	return qualityValue
	//return (foodQuality * 1.00) + (housingQuality * .3) + (electronicsQuality * .15) + (metalsQuality * .15) +
	//	(landQuality * .20) + (militaryQuality * .50)
}

// CalculateUndiscountedReward will calculate the undiscounted reward for a Country. This is equivalent to:
//   newQuality - originalQuality
func CalculateUndiscountedReward(countriesMap map[string]scheduler.Country, countryNew scheduler.Country) float64 {
	newQuality := countryNew.GetQualityRating()
	originalQuality := countriesMap[countryNew.GetName()].GetQualityRating()
	undiscountedReward := newQuality - originalQuality
	if newQuality != originalQuality {
		_ = undiscountedReward
	}
	if math.IsNaN(undiscountedReward) {
		_ = undiscountedReward
	}
	return undiscountedReward
}

// CalculateDiscountedReward will calculate the discounted reward for a Country. This is equivalent to:
//   undiscountedReward * Gamma^n
func CalculateDiscountedReward(country scheduler.Country, num int) float64 {
	undiscountedReward := country.GetUndiscountedReward()
	discountedReward := undiscountedReward * math.Pow(Gamma, float64(num))
	if math.IsNaN(discountedReward) {
		_ = discountedReward
	}
	return discountedReward
}
