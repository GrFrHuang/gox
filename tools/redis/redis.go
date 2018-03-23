package redis

import (
	"github.com/gomodule/redigo/redis"
	"encoding/json"
	"io/ioutil"
	"gox/log"
	"time"
)

var Conn *Redis

// Redis config target is default from where read current path config.json.
type Config struct {
	Host        string `json:"host"`          // Redis server host.
	Port        string `json:"port"`          // Redis server host's port.
	Protocol    string `json:"protocol"`      // Protocol cluster, default tcp.
	TimeOut     int    `json:"time_out"`      // Connection timeout time.
	IsKeepAlive bool   `json:"is_keep_alive"` // Whether keep long connection or not.
	Valid       bool   `json:"valid"`         // Current state whether available or not.
}

// Extend config for distinct the develop and product environment.
type ExtendConfig struct {
	Dev  *Config `json:"dev"`
	Prod *Config `json:"prod"`
}

type Connection interface {
	redis.ConnWithTimeout
}

type Redis struct {
	Connection
}

// Create a redis connect pool target by config.
func NewRedis() *Redis {
	var econfig *ExtendConfig
	var config *Config
	bt, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Error("[redis]: ", err)
		return nil
	}
	err = json.Unmarshal(bt, &econfig)
	if err != nil {
		log.Error("[redis]: ", err)
		return nil
	}
	if econfig.Dev.Valid {
		config = econfig.Dev
	} else {
		config = econfig.Prod
	}
	options := redis.DialOption{

	}
	connect, err := redis.Dial(config.Protocol, config.Host+":"+config.Port, options)
	if err != nil {
		log.Error("[redis]: ", err)
		return nil
	}
	log.Info("[redis]: success to connect redis !")
	return &Redis{
		Conn: connect,
	}
}

func (conn *Redis) Do() {

}

func (conn *Redis) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {

	return
}

func (conn *Redis) ReceiveWithTimeout(timeout time.Duration) (reply interface{}, err error) {

	return
}

func (conn *Redis) Flush() {

}

func (conn *Redis) Close() {

}

func (conn *Redis) Send() {

}

func (conn *Redis) Receive() {

}

func (conn *Redis) Err() {

}

func (conn *Redis) Get() {

}

func (conn *Redis) Set() {

}

func (conn *Redis) Delete() {

}

func init() {
	Conn = NewRedis()
}
