package sha

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSha256Hash(text string) string {
	// Create a new SHA256 hash
	hash := sha256.New()

	// Write the data to the hash
	hash.Write([]byte(text))

	// Get the resulting hash as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
