package main

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	plainText := ""
	encryptedText, err := Encrypt(plainText, xKey)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(encryptedText)
}

func TestDecrypt(t *testing.T) {
	encryptedText := "3xkarmSuNsZFHzgRcKyj2YO2zEQE/mSqEuB0ob5CvMH71p51egAdvAFIQif+WC79mzGBnUos64nWAJn1uLHxDQ=="
	decryptedText, err := Decrypt(encryptedText, xKey)
	if err != nil {
		t.Error(err)
	}

	t.Log(decryptedText)
}
