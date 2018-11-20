package util

import (
	"math/rand"
	"time"
)

// RandomBool random bool value
func RandomBool() bool {
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)

	return r.Intn(2) != 0
}

// RandomNumber random number value
func RandomNumber(length int) int {
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)

	return r.Intn(length)
}
