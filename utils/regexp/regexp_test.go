package regexp

import (
	"testing"
	"fmt"
	"github.com/xxtea/xxtea-go/xxtea"
	"regexp"
)

func TestEmail(t *testing.T) {
	fmt.Println(Email("grfrhuang@163.com"))
}

func TestMobilePhone(t *testing.T) {
	fmt.Println(MobilePhone("15100110011"))
}

func TestTellPhone(t *testing.T) {
	fmt.Println(TellPhone("11011011"))
}

func TestXxteaDecrypt(t *testing.T) {
	key := "12345678"
	str := "a5wFeSsjapPjnsmL5dwEx2PrAf8zDK9q0ApLiMiWdBhJ0P2JTaXLNZqQMyGAdtAAeG85+Q=="
	fmt.Println(xxtea.DecryptString(str, key))
	account := "AB12345678910"
	fmt.Println(regexp.MatchString(`^AB[0-9]{11}$`, account))
}

func TestXxteaEncrypt(t *testing.T) {
	key := "12345678"
	data := []byte{}
	fmt.Println(string(xxtea.Encrypt(data, []byte(key))))
}