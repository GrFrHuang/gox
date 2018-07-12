// RabbitMQ is an open source AMQP implementation, written in Erlang on the server side.
// The current package encapsulates the publish/subscribe mode of RabbitMQ
// RabbitMQ's ack mechanism lets consumers return ack messages to ensure that messages are consumed

// The process:
//	1. The message sender sends the message to exchange
//	2. After exchange receives a message, it is responsible for routing it to a specific queue
//	3. Bindings is responsible for connecting exchange and queue.
//	4. The message arrives in the queue and waits to be processed by the message receiver
//	5. The message receiver processes the message

package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/GrFrHuang/gox/log"
	"errors"
	"time"
	"bytes"
)

type RabbitMQ struct {
	Error   chan error
	TimeOut chan int
	Url     string
	Conn    *amqp.Connection
}

type Handler func(amqp.Delivery) error

// Use RabbitMQ server node's url initialize the RabbitMQ client by given.
// url format: amqp://user:password@host:port
func NewRabbit(url string) (*RabbitMQ, error) {
	if url == "" {
		return nil, errors.New("[RabbitMQ] url not correct")
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	log.Info("[RabbitMQ] success to connect server: ", url)
	return &RabbitMQ{
		Error:   make(chan error),
		TimeOut: make(chan int),
		Url:     url,
		Conn:    conn,
	}, nil
}

// As a producer publish the message from RabbitMQ server node.
func (r *RabbitMQ) Publish(exchange, queue string, body []byte) (error) {
	if queue == "" || body == nil {
		return errors.New("[RabbitMQ] param not correct")
	}
	// Create a channel like a concurrent multithreading mode.
	channel, err := r.Conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_err := channel.Close()
		if _err != nil {
			log.Error(_err)
		}
	}()
	// Declare the channel's exchange machine.
	err = channel.ExchangeDeclare(exchange, "topic", false, false, false, true, nil)
	if err != nil {
		return err
	}
	q, err := channel.QueueDeclare(
		queue, // message queue name
		false, // durable
		false, // delete when unused
		false, // exclusive, if connection disconnect, delete queue or not
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	routeKey := exchange + "." + q.Name
	err = channel.QueueBind(q.Name, routeKey, exchange, true, nil)
	if err != nil {
		return err
	}
	err = channel.Publish(
		"",     // exchange 默认模式，exchange为空
		q.Name, // routing key 默认模式路由到同名队列，即是传入的queue
		false,  // mandatory
		false,
		amqp.Publishing{
			// persistence, because the queue is declared lasting, news must add this (probably not),
			// but the message or may be lost, such as message to the cache but MQ hang up too late to persistence.
			//DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        body,
		},
	)
	return err
}

// As a consumer receive the message from RabbitMQ server node.
func (r *RabbitMQ) Receive(exchange, queue string, handler Handler) (error) {
	if queue == "" || exchange == "" {
		return errors.New("[RabbitMQ] param not correct")
	}
	channel, err := r.Conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_err := channel.Close()
		if _err != nil {
			log.Error(_err)
		}
	}()
	// exchange有4个类型：direct\topic\fanout\header。
	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, true, nil)
	if err != nil {
		return err
	}
	q, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		return err
	}
	routeKey := exchange + "." + q.Name
	err = channel.QueueBind(q.Name, routeKey, exchange, true, nil)
	if err != nil {
		return err
	}
	messages, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	// Messages is goroutine safe.
	go func() {
		for msg := range messages {
			err := handler(msg)
			if err != nil {
				log.Error(err)
			}
			// Confirm receive this msg, multiple must be false.
			err = msg.Ack(false)
			if err != nil {
				log.Error(err)
			}
		}
	}()
	<-forever
	return nil
}

// Get RabbitMQ server node's url.
func (r *RabbitMQ) GetUrl() (string) {
	defer func() {
		recover()
	}()
	if r == nil {
		log.Panic("[RabbitMQ] client is not initialize")
		return ""
	}
	return r.Url
}

// Close the RabbitMQ connection, at the same time, channel will be closed.
func (r *RabbitMQ) Close(topic, queue string) (error) {
	channel, err := r.Conn.Channel()
	if err != nil {
		return err
	}
	err = channel.Close()
	if err != nil {
		return err
	}
	err = r.Conn.Close()
	return err
}

// Test current connection is available or not.
func (r *RabbitMQ) Ping() (err error) {
	channel, err := r.Conn.Channel()
	if err != nil {
		return
	}
	defer channel.Close()
	err = channel.ExchangeDeclare("ping.ping", "topic", false, true, false, true, nil)
	if err != nil {
		return err
	}
	msg := "hello GrFrHuang !"
	err = channel.Publish("ping.ping", "ping.ping", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		return err
	}
	// Delete Exchange after complete ping test.
	err = channel.ExchangeDelete("ping.ping", false, false)
	return err
}

func (r *RabbitMQ) ListenMessageTimeOut(routeKey string, id int, body []byte, duration time.Duration) () {
	ticker := time.NewTicker(duration)
	go func(rabbit *RabbitMQ) {
		for range ticker.C {
			ch, err := rabbit.Conn.Channel()
			if err != nil {
				log.Error(err)
				r.Error <- err
			}
			msgs, ok, err := ch.Get(routeKey, true)
			if err != nil || !ok {
				log.Error("error: ", err)
				r.Error <- err
			}
			if bytes.Equal(msgs.Body, body) {
				r.TimeOut <- id
			}
		}
	}(r)
}
