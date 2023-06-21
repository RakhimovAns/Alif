package pkg

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

func CalculateHMACSHA1(data map[string]string, key string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hmacKey := []byte(key)
	h := hmac.New(sha1.New, hmacKey)
	h.Write(jsonData)
	hmacSHA1 := hex.EncodeToString(h.Sum(nil))

	return hmacSHA1, nil
}
