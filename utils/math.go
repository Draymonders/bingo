package utils

import (
	"math"
)

// FormatFloat64 小数 四舍五入，bit为保留几位小数
func FormatFloat64(x float64, bit int) float64 {
	isNegative := false
	if x < 0 {
		isNegative = true
		x = -x
	}
	// 保证传来的bit >= 1
	mul := 1
	for i := 1; i <= bit; i++ {
		mul *= 10
	}
	res := math.Trunc(x*float64(mul)+0.5) / float64(mul)
	if isNegative {
		res = -res
	}
	return res
}
