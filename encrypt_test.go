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
	encryptedText := "i/LOatue4KyPz9MRDB61XW9BIez/ZMyRD2/EbR0oOPWt7dVA1Jg5R5UKy02vEJotBbb4p6ohzEVjf0AD+SFhrS4RWldSzpH3dlABnVzKpBNtDpvCPKl/4/fTP2sKlyOFTEUUV74vgaab8FjJsKwXeV4PJOhSIoJFreB3hLSIQZRNBE75mM1oLvGTsWrm8Ll9"
	decryptedText, err := Decrypt(encryptedText, xKey)
	if err != nil {
		t.Error(err)
	}

	t.Log(decryptedText)
}
