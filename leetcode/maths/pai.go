package maths

import (
	"math/rand"
	"time"
)

// GetPi 求π
func GetPi() float64 {
	rand.Seed(time.Now().Unix())
	N := 1e7

	sum := 0
	for i := 0; i < int(N); i++ {
		var x, y = float64(rand.Intn(100)) / 100, float64(rand.Intn(100)) / 100

		if x*x+y*y <= 1 {
			sum++
		}
	}
	return float64(sum) / N * 4
}
