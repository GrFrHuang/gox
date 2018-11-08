// The Redis list is a simple string list,
// You can add an element to the head (left) or tail (right) of the list.
// A list can contain up to 232-1 elements 4294967295,
// more than 4 billion elements per list.

package redis

import (
	_redis "github.com/gomodule/redigo/redis"
)

// Remove and get the first element of the list.
func (redis *Redis) LPop(list interface{}) (error) {
	_, err := redis.conn.Do("LPOP", list)
	return err
}

// Remove and get the last element of the list.
func (redis *Redis) RPop(list interface{}) (error) {
	_, err := redis.conn.Do("RPOP", list)
	return err
}

// Remove and get the first element of the list,
// which blocks the list if there are no elements in the list
// until it waits for a timeout or finds an element that can pop up.
func (redis *Redis) BlPop(timeOut interface{}, keys ... interface{}) (error) {
	keys = append(keys, timeOut)
	_, err := redis.conn.Do("BLPOP", keys...)
	return err
}

// Remove and get the last element of the list,
// which blocks the list if no element is present
// until it waits for a timeout or finds an element that can pop up.
func (redis *Redis) BrPop(timeOut interface{}, keys ... interface{}) (error) {
	keys = append(keys, timeOut)
	_, err := redis.conn.Do("BRPOP", keys...)
	return err
}

// Insert one or more values into the list header.
func (redis *Redis) LPush(list interface{}, values ... interface{}) (error) {
	array := make([]interface{}, 1)
	array[0] = list
	array = append(array, values...)
	_, err := redis.conn.Do("LPUSH", array...)
	return err
}

// Get element by list index, the index is list index - 1.
func (redis *Redis) LIndex(list, index interface{}) (string, error) {
	result, err := _redis.String(redis.conn.Do("LINDEX", list, index))
	return result, err
}

// update list by index the value of a list element .
func (redis *Redis) LSet(list, index, value interface{}) (error) {
	_, err := redis.conn.Do("LSET", list, index, value)
	return err
}

// Get the elements in the list start to stop scope.
func (redis *Redis) LRange(list, start, stop interface{}) ([]string, error) {
	result, err := _redis.Strings(redis.conn.Do("LRANGE", list, start, stop))
	return result, err
}

// Get list length.
func (redis *Redis) LLen(list interface{}) (int, error) {
	length, err := _redis.Int(redis.conn.Do("LLEN", list))
	return length, err
}
