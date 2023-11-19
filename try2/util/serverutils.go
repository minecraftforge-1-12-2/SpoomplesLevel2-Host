package util

import (
	"TogetherForever/config"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"unicode/utf8"
)

var BadWords = config.ParseLsf("badwords.lsf", config.BADWORDS)

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
	input = strings.TrimSpace(input)
	if len(input) > length && length != -1 {
		input = input[:length]
	}
	for _, w := range BadWords {
		re := regexp.MustCompile("(?i)" + strings.TrimSpace(w))
		input = re.ReplaceAllString(input, strings.Repeat("*", utf8.RuneCountInString(w)))
	}
	return input
}

func CleanName(input string) string {
	Runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-!@#$%^&*()+=[]{}~ ")

	input = strings.TrimSpace(input)
	if len(input) > 16 {
		input = input[:16]
	}

	for _, r := range input {
		if !utf8.ValidRune(r) {
			input = strings.Replace(input, string(r), "", -1)
		} else if !strings.ContainsRune(string(Runes), r) {
			input = strings.Replace(input, string(r), "", -1)
		}
	}

	input = strings.Replace(input, " ", "-", -1)

	if input == "" {
		input = "Player"
	}

	return input
}

func Anticheat(sprite string) string {
	if !strings.HasPrefix(sprite, "spr_player") && !strings.HasPrefix(sprite, "spr_knight") && !strings.HasPrefix(sprite, "spr_shotgun") && !strings.HasPrefix(sprite, "spr_ratmount") && !strings.HasPrefix(sprite, "spr_lone") && sprite != "spr_noise_vulnerable2" && sprite != "spr_noise_crusherfall" {
		sprite = "spr_player_idle"
	}
	return sprite
}
