package util

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)


// MakeIp256 returns a sha256 hash of the given ip.
func MakeIp256(ip string) string {
	s := sha256.New()
	s.Write([]byte(ip))
	return fmt.Sprintf("%x", s.Sum(nil))
}

// KeyGen generates a random verification key.
// Keys are 32 characters long, and contain only alphanumeric characters.
func KeyGen() string {
	Runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var key string
	for i := 0; i < 32; i++ {
		key += string(Runes[rand.Intn(len(Runes))])
	}
	return key
}

func Clean(input string, length int) string {
	return input
}

func CleanName(input string) string {
	return input
}

func Anticheat(sprite string) string {
	return sprite
}
