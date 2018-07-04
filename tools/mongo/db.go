// You can query JSON data and store large amounts of data, but you don't support transactions.

// Redis data is stored in memory and is written to disk regularly.
// When there is insufficient memory, you can select the specified LRU algorithm to delete the data.
// All the data of mongodb is actually stored on the hard disk,
// and all the data to be operated on is mapped to an area of memory by mmap.

// If the table structure changes frequently, you don't need to store complex data structures,
// and you need to store document-type data in real time, and you need to extends The amount of data aways, you need mongoDB.\

// Bson(Binary Json) is Json's extends data format for mongoDB.

package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/GrFrHuang/gox/log"
	"time"
	"strconv"
)

var session *mgo.Session
var err error

type Config struct {
	Host           string
	Port           string
	Direct         string // 主节点发生故障时, 是否与集群中其他被选举的节点建立连接
	Timeout        string
	Username       string
	Password       string
	MaxConnections string // Session.SetPoolLimit
	Database       string
}

func Init(config *Config) {
	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		log.Panic(err)
	}
	maxConnections, err := strconv.Atoi(config.MaxConnections)
	if err != nil {
		log.Panic(err)
	}
	isDirect, err := strconv.ParseBool(config.Direct)
	if err != nil {
		log.Panic(err)
	}
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{config.Host + ":" + config.Port},
		Direct:    isDirect,
		Timeout:   timeout,
		Username:  config.Username,
		Password:  config.Password,
		PoolLimit: maxConnections,
		Database:  config.Database,
	}
	// Create a tcp socket pool and gain it's session.
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Panic("[mongdb] ", err)
	}
	session.SetMode(mgo.Monotonic, true)
}

func GetSession() *mgo.Session {
	return session.Clone()
}

// Get default database by config file.
func GetDefaultDatabase() *mgo.Database {
	return session.DB("")
}

func GetDatabase(database string) *mgo.Database {
	return session.DB(database)
}
