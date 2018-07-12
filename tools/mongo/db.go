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
	"gopkg.in/mgo.v2/bson"
	"github.com/GrFrHuang/gox/log"
	"time"
	"strconv"
	"reflect"
	"strings"
	"errors"
	"encoding/json"
)

var (
	engine = new(Engine)
	err    error
)

var (
	ErrNotFound = mgo.ErrNotFound
	ErrCursor   = mgo.ErrCursor
)

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

type Engine struct {
	*mgo.Session
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
	engine.Session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Panic("[mongdb] ", err)
	}
	engine.Session.SetMode(mgo.Monotonic, true)
	//	Print sql expression.
	mgo.SetDebug(true)
}

func GetEngine() *Engine {
	err = engine.Session.Ping()
	if err != nil {
		log.Panic("[mongdb] ", err)
	}
	return &Engine{
		Session: engine.Session.Clone(),
	}
}

// Get default database by config file.
func (e *Engine) GetDefaultDatabase() *mgo.Database {
	return e.Session.DB("")
}

func (e *Engine) GetDatabase(database string) *mgo.Database {
	return e.Session.DB(database)
}

func NewBsonFromJson(v interface{}, ignoreFields ... string) (bson.M, error) {
	var value = reflect.ValueOf(v)
	var elem = value.Type().Elem()
	var doc = bson.M{}
	var err error
	if elem.Kind() != reflect.Struct {
		err = errors.New("Type not is reflect.Struct ! ")
		log.Error(err)
		return nil, err
	}
	// Protobuf object have three fields always.
	if elem.NumField() <= 3 {
		err = errors.New("The lack of Field ! ")
		log.Error(err)
		return nil, err
	}
	m := make(map[string]bool)
	for _, v := range ignoreFields {
		m[v] = true
	}
	for i := 0; i < elem.NumField(); i++ {
		jsonTag := elem.Field(i).Tag.Get("json")
		array := strings.Split(jsonTag, ",")
		if strings.ToLower(array[0]) == "id" || strings.ToLower(elem.Field(i).Name) == "id" {
			// Change the field id to _id.
			doc["_id"] = value.Elem().Field(i).Interface()
			continue
		}
		if array[0] != "-" || !m[array[0]] {
			doc[array[0]] = value.Elem().Field(i).Interface()
		}
	}
	return doc, nil
}

func NewJsonFromBson(m bson.M, v interface{}, ignoreFields ... string) (error) {
	var null interface{}
	if len(ignoreFields) > 0 {
		for _, v := range ignoreFields {
			m[v] = null
		}
	}
	// Change the field _id to id.
	if v, ok := m["_id"]; ok {
		m["id"] = v
		delete(m, "_id")
	}
	bts, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bts, &v)
	if err != nil {
		return err
	}
	return nil
}
