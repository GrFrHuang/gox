package aes

import (
	"testing"
	"fmt"
)

func TestAesEncrypt(t *testing.T) {
	cipherText, err := Encrypt([]byte("hello world"), []byte(IV))
	fmt.Println(cipherText, err)
}

func TestAesEncrypt2(t *testing.T) {
	cipherText, _ := Encrypt([]byte("hello world"), []byte(IV))
	decipherText, err := Decrypt(cipherText, []byte(IV))
	fmt.Println(string(decipherText), err)
}
