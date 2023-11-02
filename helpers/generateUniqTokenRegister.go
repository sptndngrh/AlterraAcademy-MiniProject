package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateUniqueToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)

}
