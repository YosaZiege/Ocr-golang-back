package util

import (
	crand "crypto/rand" // for secure bytes
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // seed once per call (simple usage)
	result := make([]byte, length)
	for i := range result {
		index := rand.Intn(len(alphabet))
		result[i] = alphabet[index]
	}
	return string(result)
}

// RandomInit generates a random integer between min and max
func RandomInit(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomUsername returns a random username + 4-digit number
func RandomUsername() string {
	usernames := []string{
		"ninja", "shadow", "warrior", "falcon", "panther",
		"samurai", "ghost", "hunter", "viper", "phoenix",
	}
	base := usernames[rand.Intn(len(usernames))]
	suffix := rand.Intn(10000)
	return fmt.Sprintf("%s%04d", base, suffix)
}

// RandomEmail returns an email like ninja2341@gmail.com
func RandomEmail() string {
	name := RandomUsername()
	domains := []string{"gmail.com", "yahoo.com", "protonmail.com", "hotmail.com"}
	domain := domains[rand.Intn(len(domains))]
	return fmt.Sprintf("%s@%s", strings.ToLower(name), domain)
}

// RandomPasswordHash returns a random hex string (simulating hashed passwords)
func RandomPasswordHash() string {
	bytes := make([]byte, 16) // 128-bit hash
	_, _ = crand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// RandomProvider returns a random auth provider
func RandomProvider() string {
	providers := []string{"google", "github", "none", "facebook"}
	return providers[rand.Intn(len(providers))]
}
func RandomFilename() string {
	extensions := []string{"pdf", "docx", "png", "jpg"}
	name := RandomUsername()
	ext := extensions[rand.Intn(len(extensions))]
	return fmt.Sprintf("%s.%s", strings.ToLower(name), ext)
}

func RandomContent() string {
	snippets := []string{
		"Lorem ipsum dolor sit amet.",
		"Sample OCR text result from scan.",
		"This is a test extracted paragraph.",
		"Data recognized from image file.",
	}
	return snippets[rand.Intn(len(snippets))]
}
