// Redis hash is a string type field and value mapping table
// that is especially suitable for storing objects.
// Each hash in Redis can store 232-1 key-value pairs or more than 4 billion lbs.

package redis

import (
	_redis "github.com/gomodule/redigo/redis"
)

// map[field]value
type HmSetBean map[interface{}]interface{}

// Get field, value from hash table by table and field name.
func (redis *Redis) HGet(table, field interface{}) (string, error) {
	result, err := _redis.String(redis.conn.Do("HGET", table, field))
	if err == _redis.ErrNil {
		return "", nil
	}
	return result, err
}

// Get all values by field list.
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

// Hash table set field and value.
func (redis *Redis) HSet(table, field, value interface{}) (error) {
	_, err := redis.conn.Do("HSET", table, field, value)
	return err
}

// Hash table set field and value if key not exist.
func (redis *Redis) HSetNx(table, field, value interface{}) (error) {
	_, err := redis.conn.Do("HSETNX", table, field, value)
	return err
}

// Hash table set multiple fields and values.
func (redis *Redis) HmSet(table interface{}, beans []HmSetBean) (error) {
	var params = []interface{}{table}
	for _, bean := range beans {
		for k, v := range bean {
			params = append(params, k, v)
		}
	}
	_, err := redis.conn.Do("HMSET", params...)
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

// Get all values from hash table, not need fields.
func (redis *Redis) HVals(table interface{}) ([]string, error) {
	results, err := _redis.Strings(redis.conn.Do("HVALS", table))
	return results, err
}

// Check field whether exist or not.
func (redis *Redis) HExists(table, field interface{}) (bool, error) {
	results, err := _redis.Int(redis.conn.Do("HKEYS", table))
	if results != 0 && err == nil {
		return true, err
	}
	return false, err
}
