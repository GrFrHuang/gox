package rabbitmq

import (
	"testing"
	"github.com/GrFrHuang/gox/log"
	"strconv"
)

var rabbit *RabbitMQ

func TestRabbit(t *testing.T) {
	var err error
	rabbit, err = NewRabbit("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
		log.Error(err)
		return
	}
	for i := 0; i < 15; i++ {
		body := []byte("hello GrFrHuang " + strconv.Itoa(i))
		err = rabbit.Publish("ex1", "test1", body)
		if err != nil {
			log.Error(err)
			return
		}
	}
	ch := make(<-chan int)
	<-ch
}
