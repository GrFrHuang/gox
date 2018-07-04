// You can query JSON data and store large amounts of data, but you don't support transactions.

// Redis data is stored in memory and is written to disk regularly.
// When there is insufficient memory, you can select the specified LRU algorithm to delete the data.
// All the data of mongodb is actually stored on the hard disk,
// and all the data to be operated on is mapped to an area of memory by mmap.

// If the table structure changes frequently, you don't need to store complex data structures,
// and you need to store document-type data in real time, and you need to extends The amount of data aways, you need mongoDB.\

// Bson(Binary Json) is Json's extends data format for mongoDB.

package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"github.com/GrFrHuang/gox/log"
	"strconv"
	"github.com/Unknwon/goconfig"
)

var session *mgo.Session
var Database *mgo.Database

// todo 重试机制的方法
func main() {
	configFile, err := goconfig.LoadConfigFile("./gorm/mongo/config.ini")
	if err != nil {
		panic(err)
	}
	runMode, err := configFile.GetValue("currentEnv", "mode")
	if err != nil {
		panic(err)
	}
	var conf map[string]string
	if runMode == "dev" {
		conf, err = configFile.GetSection("dev")
	} else {
		conf, err = configFile.GetSection("prod")
	}
	if err != nil {
		panic(err)
	}
	direct, err := strconv.ParseBool(conf["direct"])
	if err != nil {
		log.Error(err)
	}
	maxConnections, err := strconv.Atoi(conf["maxConnections"])
	if err != nil {
		log.Error(err)
	}
	timeout, err := time.ParseDuration(conf["timeout"])
	if err != nil {
		log.Error(err)
	}
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{conf["host"] + ":" + conf["port"]},
		Direct:    direct,
		Timeout:   timeout,
		Username:  conf["username"],
		Password:  conf["password"],
		PoolLimit: maxConnections, // Session.SetPoolLimit
		Database:  conf["dbName"],
	}
	// Create a tcp socket pool and gain it's session.
	session, err := mgo.DialWithInfo(dialInfo)
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
