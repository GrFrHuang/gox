// RabbitMQ is an open source AMQP implementation, written in Erlang on the server side.
// The current package encapsulates the publish/subscribe mode of RabbitMQ
// RabbitMQ's ack mechanism lets consumers return ack messages to ensure that messages are consumed

// The process:
//	1. The message sender sends the message to exchange
//	2. After exchange receives a message, it is responsible for routing it to a specific queue
//	3. Bindings is responsible for connecting exchange and queue.
//	4. The message arrives in the queue and waits to be processed by the message receiver
//	5. The message receiver processes the message

// 中文：
//（1）客户端连接到消息队列服务器，打开一个channel。
//（2）客户端声明一个exchange，并设置相关属性。
//（3）客户端声明一个queue，并设置相关属性。
//（4）客户端使用routing key，在exchange和queue之间建立好绑定关系。
//（5）客户端投递消息到exchange。

// RabbitMQ provide a default exchange machine for queue, the mode default is direct.
// todo how to judge timeout from unack msg ?
package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/GrFrHuang/gox/log"
	"errors"
	"fmt"
)

type RabbitMQ struct {
	options        amqp.Table
	timeoutMessage chan interface{} // use the DLX(dead-letter-exchange) to make a timeout pool
	timeout        int64            // second
	url            string
	conn           *amqp.Connection
}

type Handler func(amqp.Delivery) error

// Use RabbitMQ server node's url initialize the RabbitMQ client by given.
// url format: amqp://user:password@host:port
func NewRabbit(url string, timeout int64) (*RabbitMQ, error) {
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
		timeout:        timeout,
		url:            url,
		conn:           conn,
	}, nil
}

// As a producer publish the message from RabbitMQ server node.
func (r *RabbitMQ) Publish(exchange, queue string, body []byte, headers map[string]interface{}, isConfirm bool) (error) {
	if exchange == "" || queue == "" || body == nil {
		return errors.New("[RabbitMQ] param not correct")
	}
	// Create a channel like a concurrent multithreading mode.
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	// Declare the channel's exchange machine.
	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	// Set message deadline and send timeout's message to TimeoutQueue.
	if r.timeout > 0 {
		err = r.AllowTimeout(exchange, channel)
		if err != nil {
			return err
		}
	}
	q, err := channel.QueueDeclare(
		queue, // message queue name
		true,  // durable
		true,  // Auto delete when no consumer listen it.
		false, // 1.Exclusive, if current connection disconnect, delete queue or not.
		// 2.Whether this queue is current connection private or not.
		false,     // async create queue and not wait result.
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
	if isConfirm == true {
		// Open rabbitMQ's message confirm mode.
		err = channel.Confirm(false)
		if err != nil {
			return err
		}
		err = channel.Publish(
			exchange, // exchange  默认模式下，exchange为空，这里方法强制传入exchange名字
			routeKey, // routing key 默认模式路由到同名队列，即是传入的queue name
			false,    // 当mandatory标志位设置为true时，如果exchange根据自身类型和消息routeKey无法找到一个符合条件的queue，那么会调用basic.return方法将消息返还给生产者；当mandatory设为false时，出现上述情形broker会直接将消息扔掉
			false,    // 当immediate标志位设置为true时，如果exchange在将消息route到queue(s)时发现对应的queue上没有消费者，那么这条消息不会放入队列中。当与消息routeKey关联的所有queue(一个或多个)都没有消费者时，该消息会通过basic.return方法返还给生产者
			amqp.Publishing{
				Headers: h,
				// persistence, because the queue is declared lasting, news must add this (probably not),
				// but the message or may be lost, such as message to the cache but MQ hang up too late to persistence.
				//DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				Body:        body,
			},
		)
		ack := make(chan uint64)
		nack := make(chan uint64)
		channel.NotifyConfirm(ack, nack)
		select {
		case <-nack:
			err = fmt.Errorf("msg not arrive: %s", string(body))
		default:
			return nil
		}
		return err
	} else {
		err = channel.Publish(exchange, routeKey, false, false, amqp.Publishing{Headers: h, ContentType: "text/plain", Body: body,})
	}
	return err
}

// As a consumer receive the message from RabbitMQ server node.
func (r *RabbitMQ) Receive(exchange, queue string, handler Handler) (error) {
	var err error
	if exchange == "" || queue == "" {
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
			if r.timeout > 0 {
				err = r.AllowTimeout(exchange, channel)
				if err != nil {
					return
				}
			}
			q, err := channel.QueueDeclare(
				_queue,
				true,
				true,
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
					// The message requeue.
					msg.Reject(true)
				} else {
					// Confirm receive this msg, multiple must be false.
					msg.Ack(false)
				}
			}
			channel.Close()
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

// If timeout = 0, all messages never be overdue.
func (r *RabbitMQ) SetTimeout(timeout int64) {
	if timeout < 0 {
		r.timeout = 0
	}
	r.timeout = timeout
}

// Set all messages's deadline of queue.
func (r *RabbitMQ) AllowTimeout(exchange string, channel *amqp.Channel) (error) {
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
		true,
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
