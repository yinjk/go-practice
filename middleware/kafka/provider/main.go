package main

import (
	"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var (
	_topicId = "topic-test"
)

func main() {
	//provider()
	autoProvider()
}

func provider() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic:     _topicId,
		Key:       sarama.StringEncoder("key"),
		Partition: int32(-1),
	}
	var value string
	for {
		inputReader := bufio.NewReader(os.Stdin)
		value, err = inputReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		value = strings.Replace(value, "\n", "", -1)
		msg.Value = sarama.ByteEncoder(value)
		partition, offset, err := producer.SendMessage(msg)

		if err != nil {
			fmt.Println("Send Message Fail")
		}

		fmt.Printf("Partion = %d, offset = %d\n", partition, offset)
	}
}

func autoProvider() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.ChannelBufferSize = 500000

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	for i := 0; i < 100; i++ {
		go func() {
			var a int64 = 0
			for {
				msg := &sarama.ProducerMessage{
					Topic:     _topicId,
					Key:       sarama.StringEncoder("key"),
					Partition: int32(-1),
				}
				a++
				msg.Value = sarama.ByteEncoder(strconv.FormatInt(a, 10))
				_, _, err := producer.SendMessage(msg)

				if err != nil {
					fmt.Println("Send Message Fail")
				}

				//fmt.Printf("Partion = %d, offset = %d\n", partition, offset)
			}
		}()
	}
	// 通过信号量挂起程序
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//for {
	//	a++
	//	value = strings.Replace(value, "\n", "", -1)
	//	msg.Value = sarama.ByteEncoder(strconv.FormatInt(a, 10))
	//	_, _, err := producer.SendMessage(msg)
	//
	//	if err != nil {
	//		fmt.Println("Send Message Fail")
	//	}
	//
	//	//fmt.Printf("Partion = %d, offset = %d\n", partition, offset)
	//}
}
