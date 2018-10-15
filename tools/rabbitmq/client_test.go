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
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672", 30)
	if err != nil {
		log.Error(err)
		return
	}
	for i := 0; i < 101; i++ {
		h := make(amqp.Table)
		h["id_"+strconv.Itoa(i)] = strconv.Itoa(i + 100)
		body := []byte("hello GrFrHuang " + strconv.Itoa(i))
		err = rabbit.Publish("ex1", "test1", body, h, true)
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
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672", 30)
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
	queue := "TimeoutQueue"
	timeoutRouteKey := "ex1" + "." + queue
	options := make(amqp.Table)
	options["x-message-ttl"] = int64(30 * 1000)
	options["x-dead-letter-exchange"] = "ex1"
	options["x-dead-letter-routing-key"] = timeoutRouteKey
	rabbit.options = options
	err = rabbit.Receive("ex1", "test1", handler)
	if err != nil {
		log.Error(err)
	}
	<-forever
}

func TestRabbitMQ_Receive3(t *testing.T) {
	var err error
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672", 30)
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
	queue := "TimeoutQueue"
	timeoutRouteKey := "ex1" + "." + queue
	options := make(amqp.Table)
	options["x-message-ttl"] = int64(30 * 1000)
	// 超时就进入死信队列
	options["x-dead-letter-exchange"] = "ex1"
	options["x-dead-letter-routing-key"] = timeoutRouteKey
	rabbit.options = options
	err = rabbit.Receive("ex1", "test1", handler)
	if err != nil {
		log.Error(err)
	}
	<-forever
}

func TestRabbitMQ_Receive2(t *testing.T) {
	var err error
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672", 30)
	if err != nil {
		log.Error(err)
		return
	}
	queue := "TimeoutQueue"
	timeoutRouteKey := "ex1" + "." + queue
	options := make(amqp.Table)
	options["x-message-ttl"] = int64(30 * 1000)
	options["x-dead-letter-exchange"] = "ex1"
	options["x-dead-letter-routing-key"] = timeoutRouteKey
	rabbit.options = options
	// 主协程优先跑完,其他的goroutine就没有执行的机会了,后面声明了forever让主协程一直等待
	forever := make(chan bool)
	var handler = func(msg amqp.Delivery) error {
		log.Info("消费超时了: ", string(msg.Body))
		return nil
	}
	//err = rabbit.ListenTimeOut("ex1", handler)
	//if err != nil {
	//	log.Error(err)
	//}
	err = rabbit.Receive("ex1", "test1", handler)
	if err != nil {
		log.Error(err)
	}
	<-forever
}

//var TimeoutPool = make(map[int]int64)
//
//func ListenTimeout() {
//	nowtime := time.Now().Unix()
//	rw := sync.RWMutex{}
//	for k, v := range TimeoutPool {
//		rw.Lock()
//		if nowtime-v >= 10 {
//			if "缓存还在"
//			go func() {
//				refund(k)
//			}()
//			delete(TimeoutPool, k)
//			delete("删除缓存")
//		}
//		time.Sleep(time.Second * 1)
//		rw.Unlock()
//	}
//}
