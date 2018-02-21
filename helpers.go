package main

import (
	"crypto/rand"
	"math/big"
)

// RandInt - return a crypto random number between 0 and max
func RandInt(max int64) int64 {
	num, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0
	}
	return num.Int64()
}
