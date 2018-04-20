package regexp

import (
	"testing"
	"fmt"
	"github.com/xxtea/xxtea-go/xxtea"
	"regexp"
	"encoding/json"
	"encoding/base64"
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
	str := "vySbdpW91fphkbLmXNkt/vGqo6cLMXCf+x7BoHDFR3Oh0CCXtgtIO3snL4k+EHISARmleQ=="
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

type Order struct {
	Id           int     `xorm:"not null pk autoincr INT(10)" json:"id"`
	OrderNo      string  `xorm:"default '' comment('订单号') VARCHAR(50)" json:"order_no"`
	ThirdOrderNo string  `xorm:"default '' comment('第三方订单号') VARCHAR(50)" json:"third_order_no"`
	UserId       int     `xorm:"comment('关联的用户id') INT(10)" json:"user_id"`
	GameId       int     `xorm:"comment('关联的游戏id') INT(10)" json:"game_id"`
	Amount       float64 `xorm:"default 0 comment('金额') DECIMAL(10)" json:"amount"`
	PayType      int     `xorm:"comment('1.支付宝 2.微信 3.ios') TINYINT(2)" json:"pay_type"`
	State        int     `xorm:"default 0 comment('0.下单未支付 1.支付成功 2.订单超时 3.支付失败 4.取消支付') TINYINT(2)" json:"state"`
	CreateTime   int     `xorm:"INT(11) created" json:"create_time"`
	UpdateTime   int     `xorm:"INT(11) updated" json:"update_time"`
}

type Role struct {
	Id         int    `xorm:"not null pk autoincr INT(10)" json:"id"`
	GameId     int    `xorm:"comment('关联的game_id') INT(10)" json:"game_id"`
	UserId     int    `xorm:"comment('关联的用户id') INT(10)" json:"user_id"`
	RoleName   string `xorm:"default '' comment('游戏角色名') VARCHAR(255)" json:"role_name"`
	RoleGrade  string `xorm:"default '' comment('角色等级') VARCHAR(255)" json:"role_grade"`
	GameRegion string `xorm:"default '' comment('游戏区服') VARCHAR(255)" json:"game_region"`
	CreateTime int    `xorm:"INT(11) created" json:"create_time"`
	UpdateTime int    `xorm:"INT(11) updated" json:"update_time"`
}

func TestXxteaEncrypt(t *testing.T) {
	key := "12345678"
	//ls := LoginSession{
	//	Account:  "AB12345678910",
	//	Password: "123456",
	//}
	//ls := Order{
	//	UserId:  1,
	//	GameId:  2,
	//	Amount:  6.00,
	//	PayType: 1,
	//}
	ls := &Role{
		RoleName:   "刺客1",
		RoleGrade:  "72级",
		GameRegion: "青云门",
	}
	bts, _ := json.Marshal(ls)
	data := bts
	fmt.Println(string(data))
	fmt.Println(xxtea.EncryptString(string(data), key))

	encodeString := base64.StdEncoding.EncodeToString([]byte("?noti=http://baidu.com"))

	fmt.Println("-----", string(encodeString))
	decodeBytes, err := base64.URLEncoding.DecodeString(string(encodeString))
	fmt.Println("=======", string(decodeBytes), err)
}
