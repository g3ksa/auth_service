package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(target string, salt []byte) string {

	var targetBytes = []byte(target)
	var sha256Hasher = sha256.New()
	targetBytes = append(targetBytes, salt...)

	sha256Hasher.Write(targetBytes)

	var hashedTargetBytes = sha256Hasher.Sum(nil)
	var hashedTargetHex = hex.EncodeToString(hashedTargetBytes)

	return hashedTargetHex
}

func DoHashMatch(hashedString, currString string, salt []byte) bool {

	var currPasswordHash = HashString(currString, salt)

	return hashedString == currPasswordHash

}
