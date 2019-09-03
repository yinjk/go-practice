package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

var (
	wg       sync.WaitGroup
	_topicId = "topic-test"
)

func main() {
	//consumer()
	consumeGroup()
}

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d message:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}

func consumeGroup() {
	config := sarama.NewConfig()
	// Init config, specify appropriate version
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_1_0_0
	//consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "consumer-group-0", config)
	//if err != nil {
	//	panic(err)
	//}
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
