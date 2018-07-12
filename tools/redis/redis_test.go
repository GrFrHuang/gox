package redis

import (
	"testing"
	"fmt"
	"strconv"
)

var conn *Redis

func init() {
	config := Config{
		Protocol:    "tcp",
		Host:        "127.0.0.1",
		Port:        "6379",
		IsKeepAlive: true,
		TimeOut:     0,
	}
	conn = NewRedisPoolByConfig(&config).GetRedisConnection()
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

func TestHash(t *testing.T) {
	u, err := conn.HVals("human")
	defer func() {
		conn.CloseRedis()
	}()
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

func TestMember(t *testing.T) {
	fmt.Println(conn.SisMember("hehe", "GrFrHuang"))
}

func TestManage(t *testing.T) {
	conn.Ping()
}

func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		//err := conn.LPush("list2", "hello", "world", "hi", "GrFrHuang")
		fmt.Println(i)
	}
}
