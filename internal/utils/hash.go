package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash512WithSalt(str string, salt string) string {
	var sha512Hasher = sha512.New()
	sha512Hasher.Write([]byte(str+salt))
	hashedBytes := sha512Hasher.Sum(nil)
	hashed := hex.EncodeToString(hashedBytes)
	return hashed
}
