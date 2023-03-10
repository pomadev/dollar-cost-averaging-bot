package common

import "math"

func CalcAmount(price, yen, base int64) float64 {
	amount := float64(yen) / float64(price)
	fbase := float64(base)
	return math.Round(amount*fbase) / fbase
}
