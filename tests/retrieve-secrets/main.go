// Command mock-conjur-env simulates the conjur-env production binary
// for testing by generating random environment variable settings.
package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	// Constructor bytes
	letters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	digits   = "0123456789"
	specials = "~=+%^*/()[]{}/!@#$?|`\""
	all      = letters + digits + specials

	// The first letters of a Bash env variable can only be letters and `_`
	// If the first character is anything else, an error will be thrown.
	validKeyFirstBytes     = letters + "_"
	validKeyRemainingBytes = letters + "_" + digits
	// Any character is valid for a secret
	validSecretBytes = all

	// Key constraints
	keyLength       = 10
	keyCount        = 20
	maxSecretLength = 30
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomValue(n int, runes string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = runes[rand.Int63()%int64(len(runes))]
	}
	return string(b)
}

func main() {
	var settings []string

	for x := 0; x < keyCount; x++ {
		keyPrefix := randomValue(1, validKeyFirstBytes)
		keySuffix := randomValue(keyLength-1, validKeyRemainingBytes)
		key := keyPrefix + keySuffix
		secret := randomValue(rand.Intn(maxSecretLength), validSecretBytes)

		encodedValue := base64.StdEncoding.EncodeToString([]byte(secret))
		// Create a setting of the form "<variable>: <base64-encoded-value>"
		setting := fmt.Sprintf("%s: %s", key, encodedValue)
		settings = append(settings, setting)
	}

	fmt.Print(strings.Join(settings, "\n"))
}
