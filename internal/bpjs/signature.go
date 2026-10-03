package bpjs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func generateSignature(consumerID string, consumerSecret string, timestamp string) string {
	message := consumerID + "&" + timestamp

	mac := hmac.New(sha256.New, []byte(consumerSecret))
	mac.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
