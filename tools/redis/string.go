package redis

import (
	_redis "github.com/gomodule/redigo/redis"
)

// Set the value of a given key to value and return the old value of key.
func (redis *Redis) GetSet(key, value interface{}) (string, error) {
	oldValue, err := _redis.String(redis.conn.Do("GETSET", key, value))
	if err == _redis.ErrNil {
		return "", nil
	}
	return oldValue, err
}

// Gets all (one or more) values of a given key.
func (redis *Redis) MGet(key ... interface{}) ([]string, error) {
	result, err := _redis.Strings(redis.conn.Do("MGET", key...))
	return result, err
}

// Set one or more key-value pairs at the same time.
func (redis *Redis) MSet(kvs []struct {
	key   interface{}
	value interface{}
}) (error) {
	var sets []interface{}
	for _, v := range kvs {
		sets = append(sets, v.key, v.value)
	}
	_, err := redis.conn.Do("MSET", sets...)
	return err
}

// If key already exists and is a string,
// the Append function appends the specified value to the end of the modified key value.
func (redis *Redis) AppendValueForKey(key, value interface{}) (int, error) {
	valueLength, err := _redis.Int(redis.conn.Do("APPEND", key, value))
	return valueLength, err
}

// Overwrites the string value stored by a given key with the value parameter,
// starting with the offset and offset is string index - 1.
func (redis *Redis) SetRangeFromOffset(key, offset, value interface{}) (error) {
	_, err := redis.conn.Do("SETRANGE", key, offset, value)
	return err
}
