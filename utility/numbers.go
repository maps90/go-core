package utility

import (
	"math/rand"
	"time"
	"math"
)

func RandomNumber(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const number = "1234567890"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = number[rand.Intn(len(number))]
	}
	return string(result)
}

func RoundFloat(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}
