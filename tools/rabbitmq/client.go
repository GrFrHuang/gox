// RabbitMQ is an open source AMQP implementation, written in Erlang on the server side.
// The current package encapsulates the publish/subscribe mode of RabbitMQ
// RabbitMQ's ack mechanism lets consumers return ack messages to ensure that messages are consumed

// The process:
//	1. The message sender sends the message to exchange
//	2. After exchange receives a message, it is responsible for routing it to a specific queue
//	3. Bindings is responsible for connecting exchange and queue.
//	4. The message arrives in the queue and waits to be processed by the message receiver
//	5. The message receiver processes the message

// RabbitMQ provide a default exchange machine for queue, the mode default is direct.
// todo how to judge timeout from unack msg ?
package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/GrFrHuang/gox/log"
	"errors"
	"time"
)

type RabbitMQ struct {
	options        amqp.Table
	timeoutMessage chan interface{} // use the DLX(dead-letter-exchange) to make a timeout pool
	timeout        int32            // second
	url            string
	conn           *amqp.Connection
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
		options:        make(amqp.Table),
		timeoutMessage: make(chan interface{}),
		timeout:        0,
		url:            url,
		conn:           conn,
	}, nil
}

// As a producer publish the message from RabbitMQ server node.
func (r *RabbitMQ) Publish(exchange, queue string, body []byte, headers map[string]interface{}, timeout int32) (error) {
	if queue == "" || body == nil {
		return errors.New("[RabbitMQ] param not correct")
	}
	// Create a channel like a concurrent multithreading mode.
	channel, err := r.conn.Channel()
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
	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	if timeout > 0 {
		r.timeout = timeout
		err = r.SetTimeout(exchange, channel)
		if err != nil {
			return err
		}
	}
	q, err := channel.QueueDeclare(
		queue,     // message queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive, if connection disconnect, delete queue or not
		false,     // no-wait
		r.options, // arguments
	)
	if err != nil {
		return err
	}
	routeKey := exchange + "." + q.Name
	err = channel.QueueBind(q.Name, routeKey, exchange, true, nil)
	if err != nil {
		return err
	}
	h := make(amqp.Table)
	if len(headers) > 0 {
		for k, v := range headers {
			h[k] = v
		}
	}
	err = channel.Publish(
		exchange, // exchange  默认模式下，exchange为空
		routeKey, // routing key 默认模式路由到同名队列，即是传入的queue name
		false,    // mandatory
		false,
		amqp.Publishing{
			Headers: h,
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
	var err error
	if queue == "" || exchange == "" {
		return errors.New("[RabbitMQ] param not correct")
	}
	go func(_exchange, _queue string, _handler Handler) {
		for {
			channel, err := r.conn.Channel()
			if err != nil {
				log.Error(err)
				return
			}
			// exchange有4个类型：direct\topic\fanout\header。
			err = channel.ExchangeDeclare(_exchange, "topic", true, false, false, false, nil)
			if err != nil {
				log.Error(err)
				return
			}
			q, err := channel.QueueDeclare(
				_queue,
				true,
				false,
				false,
				false,
				r.options,
			)
			if err != nil {
				log.Error(err)
				return
			}
			routeKey := exchange + "." + q.Name
			err = channel.QueueBind(q.Name, routeKey, _exchange, true, nil)
			if err != nil {
				log.Error(err)
				return
			}
			//for i := 0; i < 10; i++ {
			//	msg, ok, err := channel.Get(q.Name, true)
			//	log.Info(msg.Body, ok, err)
			//}
			//return
			// noack="no manual acks"=autoack
			// Messages is goroutine safe.
			messages, err := channel.Consume(q.Name, "consumer_tag1", false, false, false, true, nil)
			if err != nil {
				log.Error(err)
				return
			}
			for msg := range messages {
				err := _handler(msg)
				if err != nil {
					log.Error(err)
				}
				// Confirm receive this msg, multiple must be false.
				//err = msg.Ack(false)
				//if err != nil {
				//	log.Error(err)
				//}
			}
			channel.Close()
			time.Sleep(time.Second * 1)
		}
	}(exchange, queue, handler)
	return err
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
	return r.url
}

// Close the RabbitMQ connection, at the same time, channel will be closed.
func (r *RabbitMQ) Close(topic, queue string) (error) {
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	err = channel.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	return err
}

// Test current connection is available or not.
func (r *RabbitMQ) Ping() (err error) {
	channel, err := r.conn.Channel()
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

func (r *RabbitMQ) SetTimeout(exchange string, channel *amqp.Channel) (error) {
	queue := "TimeoutQueue"
	timeoutRouteKey := exchange + "." + queue
	options := make(amqp.Table)
	options["x-message-ttl"] = int64(r.timeout * 1000)
	options["x-dead-letter-exchange"] = exchange
	options["x-dead-letter-routing-key"] = timeoutRouteKey
	//options["x-expires"] = int64(r.timeout * 200)
	r.options = options
	outQueue, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.QueueBind(outQueue.Name, timeoutRouteKey, exchange, true, nil)
	return err
}

// Start listen the timeout message in exchange, route key is "exchange".timeoutQueue.
func (r *RabbitMQ) ListenTimeOut(exchange string, handler Handler) (error) {
	err := r.Receive(exchange, "TimeoutQueue", handler)
	return err
}
