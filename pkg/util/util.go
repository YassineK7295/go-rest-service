package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Creates a json string with one key-value pair
func MessageJson(key string, msg string) string {
	return fmt.Sprintf("{\"%s\":\"%s\"}\n", key, strings.TrimSpace(msg))
}

// Creates a random string
// credit: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return string(b)
}
