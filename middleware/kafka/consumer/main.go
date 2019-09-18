package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	wg       sync.WaitGroup
	_topicId = "topic-test"
)

func main() {
	//consumer()
	//consumeGroup()
	consume1()
}

type exampleConsumerGroupHandler struct {
	Id     int
	handle func(msg *sarama.ConsumerMessage)
	queen  chan *sarama.ConsumerMessage
}

func NewGroupHandler(handler func(msg *sarama.ConsumerMessage)) (h *exampleConsumerGroupHandler) {
	h = &exampleConsumerGroupHandler{queen: make(chan *sarama.ConsumerMessage), handle: handler}
	h.start()
	return
}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.queen <- msg
		//fmt.Printf("Id:%d Message topic:%q partition:%d offset:%d message:%s\n", h.Id, msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h exampleConsumerGroupHandler) start() {
	fmt.Println(h)
	//启动100个协程去处理消息
	fmt.Println("start 100 goroutine to hand msg")
	for i := 0; i < 1; i++ {
		go func() {
			for msg := range h.queen {
				time.Sleep(time.Millisecond * 10000) //模拟处理耗时
				h.handle(msg)
			}
		}()
	}
}

func consume1() {
	config := sarama.NewConfig()
	// Init config, specify appropriate version
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //初次启动程序时，从队列起点开始消费
	config.Version = sarama.V2_2_0_0

	//消费者数量
	for i := 0; i < 1; i++ {
		go func(id int) {
			// kafka consumer client
			client, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "consumer-group-0", config)
			if err != nil {
				panic(err)
			}
			defer func() {
				err = client.Close()
				if err != nil {
					panic(err)
				}
			}()

			topics := []string{_topicId}

			//groupHandler := exampleConsumerGroupHandler{Id: id, once: sync.Once{}, queen: make(chan *sarama.ConsumerMessage)}
			groupHandler := NewGroupHandler(func(msg *sarama.ConsumerMessage) {
				fmt.Printf("Message topic:%q partition:%d offset:%d message:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
			})
			for {
				err := client.Consume(context.Background(), topics, groupHandler)
				if err != nil {
					log.Printf("[main] client.Consume error=[%s]", err.Error())
					// 5秒后重试
					time.Sleep(time.Second * 5)
				}
			}
		}(i)
	}

	// os signal
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm
}

//通过consumer group消费消息
func consumeGroup() {
	config := sarama.NewConfig()
	// Init config, specify appropriate version
	config.Consumer.Return.Errors = true
	// 首次消费从最开始的位置消费
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_1_0_0
	// Start with a client
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient("consumer-group-0", client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{_topicId}
		handler := exampleConsumerGroupHandler{}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

func consumer() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}
	partitions, err := consumer.Partitions(_topicId)
	if err != nil {
		panic(err)
	}
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(_topicId, partition, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for message := range pc.Messages() {
				fmt.Println(string(message.Value))
			}
		}(pc)
		wg.Wait()
		_ = consumer.Close()
	}
}
