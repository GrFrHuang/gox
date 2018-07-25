package rabbitmq

import (
	"testing"
	"github.com/GrFrHuang/gox/log"
	"strconv"
	"github.com/streadway/amqp"
)

var rabbit *RabbitMQ

func TestRabbitMQ_Publish(t *testing.T) {
	var err error
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
		log.Error(err)
		return
	}
	for i := 0; i < 100; i++ {
		h := make(amqp.Table)
		h["id_"+strconv.Itoa(i)] = strconv.Itoa(i + 100)
		body := []byte("hello GrFrHuang " + strconv.Itoa(i))
		err = rabbit.Publish("ex1", "test1", body, h)
		if err != nil {
			log.Error(err)
			return
		}
	}
	ch := make(<-chan int)
	<-ch
}

func TestRabbitMQ_Receive(t *testing.T) {
	var err error
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
		log.Error(err)
		return
	}
	// 主协程优先跑完,其他的goroutine就没有执行的机会了,后面声明了forever让主协程一直等待
	forever := make(chan bool)
	var handler = func(msg amqp.Delivery) error {
		log.Info(string(msg.Body))
		return nil
	}
	err = rabbit.Receive("ex1", "test1", handler)
	if err != nil {
		log.Error(err)
	}
	<-forever
}


func TestRabbitMQ_Receive2(t *testing.T) {
	var err error
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
		log.Error(err)
		return
	}
	// 主协程优先跑完,其他的goroutine就没有执行的机会了,后面声明了forever让主协程一直等待
	forever := make(chan bool)
	var handler = func(msg amqp.Delivery) error {
		log.Info(string(msg.Body))
		return nil
	}
	err = rabbit.Receive("ex1", "test1", handler)
	if err != nil {
		log.Error(err)
	}
	<-forever
}