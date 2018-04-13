package redis

import (
	_redis "github.com/gomodule/redigo/redis"
)

const (
	EX = "EX" // EX seconds-specifies expiration time in seconds.
	PX = "PX" // PX milliseconds-specifies expiration time in milliseconds.
)

// Get Value by key string.
func (redis *Redis) Get(key interface{}) (string, error) {
	result, err := _redis.String(redis.conn.Do("GET", key))
	return result, err
}

// Get key list by regular expression.
func (redis *Redis) KeysByRegexp(pattern string) ([]string, error) {
	result, err := _redis.Strings(redis.conn.Do("KEYS", pattern))
	return result, err
}

// Set a key-value in redis server by specify expire type and time.
func (redis *Redis) Set(key string, value interface{}, expireTime int, expireType string) (error) {
	var args []interface{}
	args = append(args, key, value)
	switch expireType {
	case EX:
		args = append(args, EX)
	case PX:
		args = append(args, PX)
	}
	if expireTime != 0 {
		args = append(args, expireTime)
	}
	_, err := redis.conn.Do("SET", args...)
	return err
}

// NX-sets the key value only if the key value does not exist.
func (redis *Redis) SetNX(key, value interface{}) (error) {
	_, err := redis.conn.Do("SET", key, value, "NX")
	return err
}

// XX-sets the key value only when the key value exists.
func (redis *Redis) SetXX(key, value interface{}) (error) {
	_, err := redis.conn.Do("SET", key, value, "XX")
	return err
}

// Delete key, if key not exist return nil error.
func (redis *Redis) Delete(key interface{}) (error) {
	_, err := redis.conn.Do("DEL", key)
	return err
}

// Rename old key to new key.
func (redis *Redis) RenameKey(oldKey, newKey interface{}) (error) {
	_, err := redis.conn.Do("RENAME", oldKey, newKey)
	return err
}

// Rename old key to new key only if new key does not exist.
func (redis *Redis) RenameNXKey(oldKey, newKey interface{}) (error) {
	_, err := redis.conn.Do("RENAMENX", oldKey, newKey)
	return err
}

// Check whether exist or not.
func (redis *Redis) Exist(key interface{}) (bool, error) {
	result, err := _redis.Int64(redis.conn.Do("EXISTS", key))
	if result != 0 && err == nil {
		return true, err
	}
	return false, err
}

// Set expire time for existent key, unit is sec.
func (redis *Redis) Expire(key interface{}, expireTime int) (error) {
	_, err := redis.conn.Do("EXPIRE", key, expireTime)
	return err
}

// Set expire time for existent key, unit is millisecond.
func (redis *Redis) PExpire(key interface{}, expireTime int) (error) {
	_, err := redis.conn.Do("PEXPIRE", key, expireTime)
	return err
}
