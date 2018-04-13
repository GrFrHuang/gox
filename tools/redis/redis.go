package redis

import (
	_redis "github.com/gomodule/redigo/redis"
	"encoding/json"
	"io/ioutil"
	"github.com/GrFrHuang/gox/log"
	"time"
	"strconv"
)

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

type Redis struct {
	conn _redis.Conn
}

type Pool struct {
	pool _redis.Pool
}

// Create a redis connect pool target by config.
func newRedis() *Redis {
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
	//options := redis.DialOption{
	//
	//}
	connect, err := _redis.Dial(config.Protocol, config.Host+":"+config.Port)
	if err != nil {
		log.Error("[redis]: ", err)
		return nil
	}
	log.Info("[redis]: success to connect redis server !")
	return &Redis{
		conn: connect,
	}
}

func GetRedisConnection(p *Pool) *Redis {
	return &Redis{
		conn: p.pool.Get(),
	}
}

// Create a redis connect pool target by config.
func NewRedisPoolByConfig(config *Config) *Pool {
	// Default connection time out is 60 second.
	var timeOut time.Duration = 60
	var err error
	if config != nil && config.TimeOut > 0 {
		timeOut, err = time.ParseDuration(strconv.Itoa(config.TimeOut))
		if err != nil {
			log.Panic("[redis]: parse connection time out time error ", err)
		}
	}
	//options := redis.DialOption{
	//
	//}
	p := _redis.Pool{
		MaxIdle:     500,
		MaxActive:   10000,
		IdleTimeout: time.Second * timeOut,
		Dial: func() (_redis.Conn, error) {
			connect, err := _redis.Dial(config.Protocol, config.Host+":"+config.Port)
			if err != nil {
				log.Panic("[redis]: ", err)
			}
			return connect, nil
		},
	}
	log.Info("[redis]: success to connect redis server !")
	return &Pool{
		pool: p,
	}
}
