package regexp

import (
	"testing"
	"fmt"
	"github.com/xxtea/xxtea-go/xxtea"
	"regexp"
	"encoding/json"
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
	str := "IXGXQbB1ExS9i3ptQ/ajWwgjAzCJw2Hz1pM0Sg=="
	fmt.Println(xxtea.DecryptString(str, key))
	account := "18690702021"
	fmt.Println(regexp.MatchString(`^\d{3,}$`, account))
}

type LoginSession struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	SmsCode  string `json:"sms_code" description:"短信验证码"`
	Type     int    `json:"type" description:"1.普通 2.手机验证码 3.三方授权"`
}

func TestXxteaEncrypt(t *testing.T) {
	key := "12345678"
	ls := LoginSession{
		Account:  "AB12345678910",
		Password: "123456",
	}
	bts, _ := json.Marshal(ls)
	data := bts
	fmt.Println(string(data))
	fmt.Println(xxtea.EncryptString(string(data), key))
}
