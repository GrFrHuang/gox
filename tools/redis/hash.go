// Redis hash is a string type field and value mapping table
// that is especially suitable for storing objects.
// Each hash in Redis can store 232-1 key-value pairs or more than 4 billion lbs.

package redis

import (
	_redis "github.com/gomodule/redigo/redis"
	"gox/log"
	"encoding/json"
)

// Get field, value from hash table by table and field name.
func (redis *Redis) HGet(table, field interface{}) (string, error) {
	result, err := _redis.String(redis.conn.Do("HGET", table, field))
	return result, err
}

// Get all value by field list.
func (redis *Redis) HmGet(table interface{}, field ... interface{}) ([]string, error) {
	array := make([]interface{}, 1)
	array[0] = table
	array = append(array, field...)
	results, err := _redis.Strings(redis.conn.Do("HMGET", array...))
	return results, err
}

// Get all fields, values from hash table by table name.
func (redis *Redis) HGetAll(table interface{}) ([]string, error) {
	results, err := _redis.Strings(redis.conn.Do("HGETALL", table))
	return results, err
}

// Get all fields, values from hash table, return json object.
func (redis *Redis) HGetAllByJson(table, obj interface{}) (interface{}, error) {
	results, err := _redis.Strings(redis.conn.Do("HGETALL", table))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	keyValues := ""
	jsonStr := ""
	for k, v := range results {
		if k%2 == 0 {
			keyValues += `"` + v + `":`
			continue
		}
		if k == len(results)-1 {
			keyValues += v
			continue
		}
		keyValues += `"` + v + `"` + ","
	}
	jsonStr = "{" + keyValues + "}"

	bts, err := json.Marshal(jsonStr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = json.Unmarshal(bts, &obj)

	return obj, err
}

// Hash table set field.
func (redis *Redis) HSet() (error) {
	_, err := redis.conn.Do("HSET")
	return err
}

// Hash table set multiple fields and values.
func (redis *Redis) HmSet() (error) {
	_, err := redis.conn.Do("HMSET")
	return err
}

// Delete one or more hash fields.
func (redis *Redis) HDelete(table interface{}, field ...interface{}) (error) {
	array := make([]interface{}, 1)
	array[0] = table
	array = append(array, field...)
	_, err := redis.conn.Do("HDEL", array...)
	return err
}

// Get all fields from hash table.
func (redis *Redis) HKeys(table interface{}) ([]string, error) {
	results, err := _redis.Strings(redis.conn.Do("HKEYS", table))
	return results, err
}

func (redis *Redis) HExists(table, field interface{}) (bool, error) {
	results, err := _redis.Int(redis.conn.Do("HKEYS", table))
	if results != 0 && err == nil {
		return true, err
	}
	return false, err
}
