package score_models

import "math"

// Upcycler Score : estimation du CO2 (kg) economise en reutilisant un objet.
// Score = poids CO2 de reference de la categorie x (1 - coefficient de l'etat).

var categoryWeights = map[string]float64{
	"clothing":    30,
	"electronics": 80,
	"furniture":   150,
	"books":       5,
	"toys":        15,
	"appliances":  100,
	"sports":      25,
	"other":       40,
}

var conditionCoefficients = map[string]float64{
	"new":      0.0,
	"like_new": 0.7,
	"good":     0.8,
	"used":     0.9,
}

const (
	DefaultCategory  = "other"
	DefaultCondition = "new"
)

func IsValidCategory(category string) bool {
	_, ok := categoryWeights[category]
	return ok
}

func IsValidCondition(condition string) bool {
	_, ok := conditionCoefficients[condition]
	return ok
}

// Compute renvoie le score CO2 (kg). Une categorie ou un etat inconnu retombe
// sur les valeurs par defaut plutot que d'echouer.
func Compute(category, condition string) float64 {
	weight, ok := categoryWeights[category]
	if !ok {
		weight = categoryWeights[DefaultCategory]
	}
	coefficient, ok := conditionCoefficients[condition]
	if !ok {
		coefficient = conditionCoefficients[DefaultCondition]
	}
	return math.Round(weight*(1-coefficient)*10) / 10
}

// Config expose les categories et etats au frontend (formulaire d'annonce).
func Config() map[string]interface{} {
	return map[string]interface{}{
		"categories": categoryWeights,
		"conditions": conditionCoefficients,
	}
}
