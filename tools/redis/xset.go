// The Set of Redis is an unordered collection of String types.
// The collection member is unique, which means that duplicate data cannot appear in the collection.
// The collections in Redis are implemented through the hash table, so add, delete, and find the complexity of O(1).
// The largest number of members in the collection is 232-1 (4294967295, each of which can store over 4 billion members).

package redis

import (
	"fmt"
	_redis "github.com/gomodule/redigo/redis"
)

type SortedSet struct {
	Member string
	Score  float64
}

// Get all fields, values from hash table by table name.
func (redis *Redis) SAdd(set interface{}, member ...interface{}) (error) {
	array := make([]interface{}, 1)
	array[0] = set
	array = append(array, member...)
	result, err := redis.conn.Do("SADD", array...)
	if result == 0 && err == nil {
		return err
	}
	return err
}

// Check a member element whether is a member of a collection (set).
func (redis *Redis) SisMember(set, member interface{}) (bool, error) {
	num, err := _redis.Int(redis.conn.Do("SISMEMBER", set, member))
	if num != 0 {
		return true, nil
	}
	if err == _redis.ErrNil {
		return false, nil
	}
	return false, err
}

// Get all fields, values from hash table by table name.
func (redis *Redis) SMembers(set interface{}, member ...interface{}) ([]string, error) {
	results, err := _redis.Strings(redis.conn.Do("SMEMBERS", set))
	return results, err
}

// Move the member element from the srcSet to the distSrc.
func (redis *Redis) SMove(srcSet, distSet, member interface{}) (error) {
	_, err := redis.conn.Do("SMOVE", srcSet, distSet, member)
	return err
}

//...
// The Redis ordered collection is also a collection of string type elements, and does not allow repeated members.
// The difference is that each element is associated with a double type score.
// Redis is the ranking of members in a collection through scores.
// The members of an ordered set are unique, but scores can be repeated.
//...
// Adds one or more members to an ordered collection, or updates the score of an existing member.
func (redis *Redis) ZAdd(set interface{}, member ...*SortedSet) (error) {
	var err error
	if len(member) < 1 {
		err = fmt.Errorf("params error: %v", member)
		return err
	}
	for _, v := range member {
		_, err = redis.conn.Do("ZADD", set, v.Score, v.Member)
	}
	return err
}

// Get the number of members from sorted set.
func (redis *Redis) ZCard(set interface{}, member ...interface{}) ([]string, error) {
	results, err := _redis.Strings(redis.conn.Do("ZCARD", set))
	return results, err
}

// Returns an ordered set of members within a specified interval.
// The location of the members is sorted by increasing the score (from small to large).
// Members with the same score value are arranged in lexicographical order.
// The start or stop default value is 0.
func (redis *Redis) ZRange(set, start, stop interface{}, showScore bool) ([]string, error) {
	var withScores string
	if showScore {
		withScores = "WITHSCORES"
	}
	results, err := _redis.Strings(redis.conn.Do("ZRANGE", set, start, stop, withScores))
	return results, err
}

func (redis *Redis) ZRevRange(set, start, stop interface{}, showScore bool) ([]string, error) {
	var withScores string
	if showScore {
		withScores = "WITHSCORES"
	}
	results, err := _redis.Strings(redis.conn.Do("ZREVRANGE", set, start, stop, withScores))
	return results, err
}

// Increment the number of a specified member in an ordered set.
func (redis *Redis) ZIncrby(set, member interface{}, number float64) (error) {
	_, err := redis.conn.Do("ZINCRBY", set, number, member)
	return err
}

// Return the index of member.
func (redis *Redis) ZRank(set, member interface{}) (int, error) {
	index, err := _redis.Int(redis.conn.Do("ZRANK", set, member))
	return index, err
}

// Returns the ranking of the specified members in an ordered set, sorted by decreasing value (from large to small).
func (redis *Redis) ZRevrank(set, member interface{}) (int, error) {
	place, err := _redis.Int(redis.conn.Do("ZREVRANK", set, member))
	return place, err
}
