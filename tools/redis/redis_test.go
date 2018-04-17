package redis

import (
	"testing"
	"fmt"
	"strconv"
)

var conn *Redis

func init() {
	conn = newRedis()
}

func TestRedisKey(t *testing.T) {
	var err error
	//err = conn.Set("GrFr", "4", 30, EX)
	err = conn.Delete(24)
	fmt.Println(err)
	//err = conn.RenameNXKey("hello3", "hello4")
	//yes, err := conn.Set("hello4")
	//fmt.Println(yes, err)
	//result, err := r.KeysByRegexp("h*")
	//fmt.Println(result, err)
}

func TestRedisString(t *testing.T) {
	//v, err := conn.GetSet("hua4", "hello world")
	//v, err := conn.MGet("hu", "hua9", "hua3")
	//v, err := conn.AppendValueForKey("hua", 66)
	//err := conn.SetRangeFromOffset("hua", 2, "meiyou")
	var s []struct {
		key   interface{}
		value interface{}
	}
	for i := 0; i < 3; i++ {
		s = append(s, struct {
			key   interface{}
			value interface{}
		}{key: "h" + strconv.Itoa(i), value: "2"})
	}
	err := conn.MSet(s)
	fmt.Println(err)
}

type Human struct {
	Sex  string `json:"sex"`
	User Us     `json:"user"`
}

type Us struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestHash(t *testing.T) {
	h := Human{}
	//user.Name = "huang"
	//user.Age = 24
	//bts, err := json.Marshal(user)
	//jsonStr := string(bts)
	//fmt.Println(jsonStr, err)
	u, err := conn.HGetAllByJson("human", h)
	//var s string
	//err := json.Unmarshal([]byte("hello world"), &s)

	fmt.Println(err, u)
}

func TestList(t *testing.T) {
	err := conn.LPush("list2", "hello", "world", "hi", "GrFrHuang")
	fmt.Println(err)
	err = conn.BlPop("20", "list2")
	fmt.Println(err)
	result, err := conn.LIndex("list2", "0")
	fmt.Println(result, err)
	results, err := conn.LRange("list2", "0", "1")
	fmt.Println(results, err)
	len, err := conn.LLen("list2")
	fmt.Println(len, err)
}

func TestManage(t *testing.T) {
	conn.Ping()
}
