package aes

import (
	"testing"
	"fmt"
	"github.com/xxtea/xxtea-go/xxtea"
	"regexp"
	"encoding/base64"
	"strings"
	"strconv"
	"time"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

func TestXxteaDecrypt(t *testing.T) {
	key := "12345678"
	str := "U0aEgtHaJbpIH11Gh3RRPp0UAXrnZTa3vda1kYHBkp/qMzePaCRM8A=="
	fmt.Println(xxtea.DecryptString(str, key))
	account := "18690702021"
	fmt.Println(regexp.MatchString(`^\d{3,}$`, account))
}

func TestXxteaEncrypt(t *testing.T) {
	key := "12345678"
	ls := ""
	bts, _ := json.Marshal(ls)
	data := bts
	fmt.Println(string(data))
	fmt.Println(xxtea.EncryptString(string(data), key))
	encodeString := base64.StdEncoding.EncodeToString([]byte("?noti=http://baidu.com"))

	fmt.Println("-----", string(encodeString))
	decodeBytes, err := base64.URLEncoding.DecodeString(string(encodeString))
	fmt.Println("=======", string(decodeBytes), err)

	fmt.Println(strings.Split("/v1/verificationCode", "?")[0])
	fmt.Println(strconv.FormatFloat(1.10, 'f', 2, 64))
	fmt.Println(strconv.Atoi(strconv.Itoa(time.Now().Year())[2:] + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Hour()) + strconv.Itoa(time.Now().Minute())))
	reg, _ := regexp.Compile(`\d+\d?`)
	fmt.Println("regexp : ", reg.FindAllString("ad4f6bdaf0f3c816957f0ce1a8c45aa7", -2))
	var str = "CHW"
	for _, v := range reg.FindAllString("ad4f6bdaf0f3c816957f0ce1a8c45aa7", -2) {
		str += v
	}
	fmt.Println(str)

	//fmt.Println(strconv.Itoa(time.Now().Year())[2:] + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Hour()) + strconv.Itoa(time.Now().Minute()))
	h := md5.New()
	h.Write([]byte("123456")) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	fmt.Printf("%x\n", hex.EncodeToString(cipherStr))
}
