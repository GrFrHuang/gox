package redis

import (
	_redis "github.com/gomodule/redigo/redis"
)

// Ping to redis server.
func (redis *Redis) Ping() (bool, error) {
	result, err := _redis.String(redis.conn.Do("PING"))
	if result == "PONG" && err == nil {
		return true, nil
	}
	return false, err
}

// Check redis current info.
func (redis *Redis) INFO() (string, error) {
	result, err := _redis.String(redis.conn.Do("INFO"))
	return result, err
}

// Persistent redis data for safety.
// This command creates the dump.rdb file in the redis-cli script current directory.
func (redis *Redis) Save() (error) {
	_, err := redis.conn.Do("SAVE")
	return err
}

// Asynchronous preservation of data in the background.
func (redis *Redis) BgSave() (error) {
	_, err := redis.conn.Do("BGSAVE")
	return err
}

// Close redis connection.
func (redis *Redis) CloseRedis() (error) {
	err := redis.conn.Close()
	return err
}

// Emptying the entire Redis server's data (removing all key from all databases).
func (redis *Redis) FlushAll() (error) {
	_, err := redis.conn.Do("FLUSHALL")
	return err
}

// Use the channel to resolve request and response.
func (redis *Redis) Send(command string, arg ... interface{}) (error) {
	err := redis.conn.Send(command, arg...)
	return err
}

// Get the data size in current redis db.
func (redis *Redis) DBSize() (int, error) {
	size, err := _redis.Int(redis.conn.Do("DBSIZE"))
	return size, err
}

// todo save snapshot management setting.
func (redis *Redis) Snapshot() (error) {
	err := redis.conn.Flush()
	return err
}

// If during the preparation phase,
// the redis command into the cache error queue,
// so all the orders will not be put into the cache queue.
// But if redis carried out after the EXEC command have a mistake,
// other commands will be executed any way.
// todo redis transaction has not atomicity, consistency, isolation?
// Before redis(2.6.5) multi operates is not atomic.
//...
// The MULTI, EXEC, WATCH, DISCARD command are the basic commands for the Redis transaction function.
// Redis can't insert a request to execute another client during the execution of one Redis transaction.
// When a client is performing a transaction,
// if it disconnects from the Redis server before invoking the MULTI command,
// it does not perform any action in the transaction; instead, the.
// If it disconnects from the Redis server after calling the EXEC command,
// all operations in the transaction are performed.
// There is a CAS(check and set) mechanism with transaction.
// ...
// Start redis transaction.
func (redis *Redis) Multi() (error) {
	_, err := redis.conn.Do("MULTI")
	return err
}

// Take out and execute command cluster from redis cache queue.
func (redis *Redis) Exec() (error) {
	_, err := redis.conn.Do("EXEC")
	return err
}

// Monitor one (or more) key,
// if this (or these) key is changed by other commands before the transaction is executed,
// the transaction will be interrupted.
func (redis *Redis) Watch(keys ... interface{}) (error) {
	_, err := redis.conn.Do("WATCH", keys...)
	return err
}

// Unmonitor all key by the WATCH command.
func (redis *Redis) UnWatch() (error) {
	_, err := redis.conn.Do("UNWATCH")
	return err
}

// Discard all commands in transaction block and exit transaction.
func (redis *Redis) Discard() (error) {
	_, err := redis.conn.Do("DISCARD")
	return err
}
