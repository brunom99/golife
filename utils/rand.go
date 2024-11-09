package utils

import (
	"math/rand"
)

func RandInt(min int, max int, r ...*rand.Rand) int {
	if len(r) > 0 {
		return r[0].Intn(max-min) + min
	}
	return rand.Intn(max-min) + min
}

func RandFloat(min float64, max float64, r ...*rand.Rand) float64 {
	if len(r) > 0 {
		return min + r[0].Float64()*(max-min)
	}
	return min + rand.Float64()*(max-min)
}
