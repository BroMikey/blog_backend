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

// RandomAmount 生成随机金额（正负数都支持）
func RandomAmount(min, max int32) int32 {
	return min + int32(rng.Int31n(max-min+1))
}

// RandomPositiveAmount 生成正数金额
func RandomPositiveAmount() int32 {
	return RandomAmount(1, 10000)
}

func RandomCoinType() string {
	coin_type := []string{"penny", "nickel", "dime"}
	// 1 5 10美分
	n := len(coin_type)
	return coin_type[rand.Intn(n)]
}
