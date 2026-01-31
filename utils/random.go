package utils

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(b)
}

func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

func RandomEmail() string {
	return RandomString(int(RandomInt(1, 10))) + "@" + RandomString(int(RandomInt(3, 8))) + ".com"
}

func RandomUsername() string {
	return RandomString(int(RandomInt(1, 6)))
}

func RandomPassword() string {
	return RandomString(12)
}
