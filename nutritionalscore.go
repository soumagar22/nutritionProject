package main

type ScoreType int

const (
	Food ScoreType = iota // Akhil had this statement as Food ScoreType= iota , after changing to Food =iota saving it then reverting it works #vsc lint issue
	Beverage
	Water
	Cheese
)

type NutritionalScore struct {
	Value     int
	Positive  int
	Negative  int
	ScoreType ScoreType
}

/*
Reference article for difference btw type and var in Go
https://medium.com/rungo/variables-and-constants-in-go-programming-c715443fa788
Basically if i do following:-
 type float float64
 This means that i am defining my own datatype float which has all the properties of float64
 it follows type newType fromType
 newType is a type derived from fromType having all properties of fromType and we can add additional properties
 on it without modifying fromType
*/

type EnergyKJ float64
type SugarGram float64
type SaturatedFattyAcidsGram float64
type SodiumMilliGram float64
type FruitsPercent float64
type FibreGram float64
type ProteinGram float64

var energyLevels = []float64{3350, 3015, 2680, 2345, 2010, 1675, 1340, 1005, 670, 335}
var sugarLevels = []float64{45, 40, 36, 31, 27, 22.5, 18, 13.5, 9, 4.5}
var saturatedFattyAcidLevels = []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
var sodiumLevels = []float64{900, 810, 720, 630, 540, 450, 360, 270, 180, 90}
var fibreLevels = []float64{4.7, 3.7, 2.8, 1.9, 0.9}
var proteinLevels = []float64{8, 6.4, 4.8, 3.2, 1.6}

var energyLevelsFromBeverage = []float64{270, 240, 210, 180, 150, 120, 90, 60, 30, 0}
var sugarLevelsFromBeverage = []float64{13.5, 12, 10.5, 9, 7.5, 6, 4.5, 3, 1.5, 0}

var scoreToLetter = []string{"A", "B", "C", "D", "E"}

func (e EnergyKJ) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(e), energyLevelsFromBeverage)
	}
	return getPointsFromRange(float64(e), energyLevels)
}

func (s SugarGram) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(s), sugarLevelsFromBeverage)
	}
	return getPointsFromRange(float64(s), sugarLevels)
}

func (s SaturatedFattyAcidsGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(s), saturatedFattyAcidLevels)
}

func (s SodiumMilliGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(s), sodiumLevels)
}

func (f FruitsPercent) GetPoints(st ScoreType) int {
	if st == Beverage {
		if f > 80 {
			return 10
		} else if f > 60 {
			return 4
		} else if f > 40 {
			return 2
		}
		return 0
	}

	if f > 80 {
		return 5
	} else if f > 60 {
		return 2
	} else if f > 40 {
		return 1
	}
	return 0
}

func (f FibreGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(f), fibreLevels)
}

func (p ProteinGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(p), proteinLevels)
}

func EnergyFromKcal(kcal float64) EnergyKJ {
	return EnergyKJ(kcal * 4.164)
}

func SodiumFromSalt(SaltMg float64) SodiumMilliGram {
	return SodiumMilliGram(SaltMg / 2.5)
}

type NutritionalData struct {
	Energy              EnergyKJ
	Sugars              SugarGram
	SaturatedFattyAcids SaturatedFattyAcidsGram
	Sodium              SodiumMilliGram
	Fruits              FruitsPercent
	Fibre               FibreGram
	Protein             ProteinGram
	IsWater             bool
}

func GetNutritionalScore(n NutritionalData, st ScoreType) NutritionalScore {
	value := 0
	positive := 0
	negative := 0

	if st != Water {
		fruitPoints := n.Fruits.GetPoints(st)
		fibrePoints := n.Fibre.GetPoints(st)

		positive = fibrePoints + n.Protein.GetPoints(st) + fruitPoints
		negative = n.Energy.GetPoints(st) + n.Sugars.GetPoints(st) + n.Sodium.GetPoints(st) + n.SaturatedFattyAcids.GetPoints(st)

		if st == Cheese {
			value = negative - positive
		} else {
			if negative >= 11 && fruitPoints < 5 {
				value = negative - positive - fruitPoints
			} else {
				value = negative - positive
			}
		}
	}

	return NutritionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		ScoreType: st,
	}
}

func (ns NutritionalScore) GetNutriScore() string {
	if ns.ScoreType == Food {
		return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{18, 10, 2, -1})]
	}
	if ns.ScoreType == Water {
		return scoreToLetter[0] // Always grade A for water
	}
	return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{9, 5, 1, -2})]
}

func getPointsFromRange(value float64, steps []float64) int {
	lenSteps := len(steps)

	for i, v := range steps {
		if value > v {
			return lenSteps - i //values in slices are in descending order
		}
	}
	return 0
}
