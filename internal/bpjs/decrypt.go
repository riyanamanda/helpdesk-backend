package bpjs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	lzstring "github.com/daku10/go-lz-string"
)

// DecryptResponse decrypts the AES-256-CBC payload and decompresses the LZ-String data returned by BPJS.
func DecryptResponse(key string, encryptedData string) (string, error) {
	// 1. Generate key and IV (SHA-256 of ConsID + ConsSecret + Timestamp)
	hash := sha256.Sum256([]byte(key))
	keyBytes := hash[:]
	ivBytes := hash[:16]

	// 2. Base64 decode encrypted payload
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(ciphertext) == 0 {
		return "", errors.New("empty ciphertext")
	}

	// 3. AES-256-CBC Decryption
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext length is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// Remove PKCS7 padding
	decrypted, err = pkcs7Unpad(decrypted)
	if err != nil {
		return "", err
	}

	// 4. Decompress using go-lz-string
	decompressed, err := lzstring.DecompressFromEncodedURIComponent(string(decrypted))
	if err != nil {
		return "", fmt.Errorf("lz-string decompression failed: %w", err)
	}

	return decompressed, nil
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, errors.New("invalid padding size")
	}
	return data[:length-padding], nil
}
