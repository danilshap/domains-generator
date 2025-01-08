package utils

import (
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

var random *rand.Rand

const (
	laters  = "qwertyuiopasdfghjklzxcvbnm"
	numbers = "0123456789"
	charset = laters + numbers

	length_of_random_owner = 6

	min_amount_value = int64(0)
	max_amount_value = int64(10000)
)

func init() {
	random = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
}

// generate a random int berween min and max values
func RandomInt(min, max int64) int64 {
	return min + random.Int63n(max-min+1)
}

// generate a randomg string with n length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(laters)
	for i := 0; i < n; i++ {
		c := laters[random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(5)
}

func RandomProvider() string {
	return RandomString(10)
}

func RandomEmail() string {
	username := RandomString(10)
	return username + "@test.com"
}

// Generate random string with specified length using alphanumeric characters
func RandomStringWithCharset(length int, chars string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[random.Intn(len(chars))]
	}
	return string(b)
}

// Generate random alphanumeric string
func RandomAlphanumeric(length int) string {
	return RandomStringWithCharset(length, charset)
}
