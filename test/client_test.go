package test

import (
	"fmt"
	"github.com/kuhufu/nsqgo"
	"github.com/nsqio/go-nsq"
	"log"
	"strconv"
	"testing"
	"time"
)

type handler struct{}

func (handler) HandleMessage(m *nsq.Message) error {
	log.Printf("receive message: %s\n", m.Body)
	return nil
}

func TestServer(t *testing.T) {
	conf := nsq.NewConfig()
	conf.MaxInFlight = 10
	conf.LookupdPollInterval = time.Second * 10

	c, err := nsq.NewConsumer("fruit", "team_a", conf)
	if err != nil {
		panic(err)
	}
	c.AddConcurrentHandlers(handler{}, 2)

	//err = c.ConnectToNSQDs([]string{"localhost:4150", "localhost:4250"})
	err = c.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		panic(err)
	}

	select {}
}

func TestPublish(t *testing.T) {
	p, err := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	err = p.Publish("fruit", []byte("apple"))
	if err != nil {
		log.Println(err)
	}
	p.Stop()
}

func TestPublish2(t *testing.T) {

	p, err := nsq.NewProducer("localhost:4250", nsq.NewConfig())

	if err != nil {
		panic(err)
	}
	err = p.Publish("fruit", []byte("apple"))

	if err != nil {
		log.Println(err)
	}
	p.Stop()
}

func TestPublish3(t *testing.T) {
	producer, err := nsqgo.NewProducer(
		nsqgo.WithNSQLookupd("localhost:4161"),
		nsqgo.WithNSQConfig(nsq.NewConfig()),
		nsqgo.WithReconnectInterval(time.Second*5),
		nsqgo.WithSelector(nsqgo.RoundRobin()),
	)

	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	for i := 0; i < 10000; i++ {
		err = producer.Publish("fruit", []byte("apple"+strconv.Itoa(i)))
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second)
	}
}

func Test(t *testing.T) {
	fmt.Println(0 % 2)
}
