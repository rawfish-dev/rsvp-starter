package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func GenerateRandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func GenerateRandomPhoneNumber() string {
	// For simplicity's sake just generate between 8 - 10 characters in length
	return fmt.Sprintf("%d", GenerateRandomInt(10000000, 999999999))
}
