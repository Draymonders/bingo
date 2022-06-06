package common

import "math/rand"

func GenerateIntSlice(size int, maxVal int) []int {
	res := make([]int, 0, size)
	for i := 0; i < size; i++ {
		res = append(res, rand.Intn(maxVal))
	}
	return res
}
