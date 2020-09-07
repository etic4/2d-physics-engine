package physique

import (
	"crypto/rand"
	"fmt"
)

//Abs valeur absolue de n
func Abs(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}

//Min retourne le minimum
func Min(nums ...float64) float64 {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

//Max retourne le maximum
func Max(nums ...float64) float64 {
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

//Sign retourne le signe du float64 passé
func Sign(n float64) float64 {
	if n < 0 {
		return -1
	}
	return 1
}

//Clamp clamp restreint une value à l'interval min - max
func Clamp(d float64, min float64, max float64) float64 {
	return Max(min, Min(max, d))
}

func stringListContains(stringList1 []string, str string) bool {
	for _, t1 := range stringList1 {
		if t1 == str {
			return true
		}
	}
	return false
}

//UUID Retourne un uuid maison de 16 bytes
func UUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
