package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var xKey = []byte("a46c58229f9011fd21d4f78f672bef57")

func addPadding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func removePadding(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid padding")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}

func Encrypt(plainText string, key []byte) (string, error) {
	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Pad the plain text to a multiple of the block size
	paddedPlainText := addPadding([]byte(plainText), aes.BlockSize)

	// Generate a random IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Perform AES encryption
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(paddedPlainText))
	encrypter.CryptBlocks(encrypted, paddedPlainText)

	// Combine the IV and encrypted data
	result := make([]byte, len(iv)+len(encrypted))
	copy(result[:aes.BlockSize], iv)
	copy(result[aes.BlockSize:], encrypted)

	// Encode the result in base64 to make it human-readable
	encodedResult := base64.StdEncoding.EncodeToString(result)

	return encodedResult, nil
}

func Decrypt(encodedText string, key []byte) (string, error) {
	// Decode the base64-encoded input
	decoded, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}

	// Extract the IV and encrypted data
	iv := decoded[:aes.BlockSize]
	encrypted := decoded[aes.BlockSize:]

	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Perform AES decryption
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	decrypter.CryptBlocks(decrypted, encrypted)

	// Remove padding from the decrypted result
	decrypted, err = removePadding(decrypted, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
