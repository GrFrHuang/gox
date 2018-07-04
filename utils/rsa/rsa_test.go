package encoding

import (
	"testing"
	"fmt"
)

func init() {
	InitKey("", "")
}

func TestRSAEncryptAndDecrypt(t *testing.T) {
	cipher, err := RsaEncrypt([]byte("hello world!"))
	fmt.Println(string(cipher), err)
	plain, err := RsaDecrypt(cipher)
	fmt.Println(string(plain), err)
}
