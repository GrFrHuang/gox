package redis

import (
	_redis "github.com/gomodule/redigo/redis"
	"gox/log"
	"fmt"
	"errors"
	"sync"
	"time"
)

type Message _redis.Message
type Subscription _redis.Subscription

var MessageChannel = make(chan Message)
var SubscriptionChannel = make(chan Subscription)

// Sends information to a specified channel.
func (redis *Redis) Publish(channel, value interface{}) (error) {
	result, err := _redis.Int(redis.conn.Do("PUBLISH", channel, value))
	if result == 0 && err == nil {
		err = errors.New("[redis]: publish fail")
		return err
	}
	return err
}

// Subscribe to information about a given channel or more channels.
// Current goroutine will be blocked for get message.
func (redis *Redis) Subscribe(channel ... interface{}) (error) {
	ps := _redis.PubSubConn{Conn: redis.conn}
	err := ps.Subscribe(channel...)
	if err != nil {
		log.Error(err)
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ps _redis.PubSubConn, wg *sync.WaitGroup, err error) {
		defer func() {
			wg.Done()
		}()
		for {
			switch result := ps.Receive().(type) {
			case Message:
				fmt.Printf("%s: message: %s\n", result.Channel, result.Data)
				MessageChannel <- result
			case Subscription:
				fmt.Printf("%s: %s %d\n", result.Channel, result.Kind, result.Count)
				SubscriptionChannel <- result
			case error:
				err = result
				return
			}
		}
	}(ps, wg, err)
	wg.Wait()
	return err
}

// Execute this method asynchronously to get the message for the subscription.
func (redis *Redis) GetSubscribeMessage(src []string, timeOut int64) (interface{}) {
	var after = make(<-chan time.Time)
	if timeOut > 0 {
		after = time.After(time.Duration(timeOut))
	} else {
		after = time.After(-1)
	}
	for {
		select {

		case <-MessageChannel:
			return <-MessageChannel

		case <-SubscriptionChannel:
			return <-SubscriptionChannel

		case <-after:
			goto A
		}
	}
A:
	return src
}

// Unsubscribe a given channel or more channels.
func (redis *Redis) UnSubscribe(channel ... interface{}) (error) {
	_, err := redis.conn.Do("UNSUBSCRIBE", channel...)
	return err
}
