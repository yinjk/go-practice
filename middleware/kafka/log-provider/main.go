package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	_topicId = "log-trace"
)

var (
	n int
	c int
	t string
)

func init() {
	flag.IntVar(&n, "n", 2000, "the number of log written per seconds")
	flag.IntVar(&c, "c", 8, "the concurrency of service invocation for test")
	flag.StringVar(&t, "t", "1m", "the time for test program run")
	flag.Parse()
}

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o logToKafka .
func main() {
	completed := make(chan int)
	for i := 1; i <= c; i++ {
		go writtenLog(i, completed)
	}
	count := 0
	for {
		select {
		case <-completed:
			count++
			if count == c {
				fmt.Println("time wait 10 seconds for test ...")
				time.Sleep(time.Second * 10)
				return
			}
		}
	}
}

func createKafkaProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	config.Producer.Flush.Messages = 10000
	config.Producer.Flush.Frequency = time.Millisecond * 50
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewAsyncProducer([]string{"10.10.108.34:9092", "10.10.108.35:9092", "10.10.108.36:9092", "10.10.108.42:9092", "10.10.108.43:9092"}, config)
	if err != nil {
		panic(err)
	}
	return producer
}

func writtenLog(id int, completed chan int) {
	producer := createKafkaProducer()
	duration, err := time.ParseDuration(t)
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(0)
	shutdown := time.NewTimer(duration)
	start := time.Now()
	execCount := 0
	for {
		select {
		case <-timer.C:
			timer.Reset(time.Second)
			now := time.Now()
			for i := 0; i < n; i++ {
				for _, logStr := range provideLog() {
					msg := &sarama.ProducerMessage{
						Topic:     _topicId,
						Partition: int32(-1),
						Value:     sarama.ByteEncoder(logStr),
					}
					producer.Input() <- msg
				}
			}
			execCount++
			nanoseconds := time.Now().Sub(now).Nanoseconds()
			fmt.Println(fmt.Sprintf("[id: %d] - [%d] this write times: %d ms %d us", id, execCount, nanoseconds/1000000, nanoseconds/1000))
		case <-shutdown.C: //结束信号
			log.Printf("exec count: %d, time duration is: %v", execCount, time.Now().Sub(start))
			fmt.Println("time wait 10 seconds for test ...")
			time.Sleep(time.Second * 10)
			completed <- 1
			return

		}
	}
}

// traceId: 0DCN001TRANSFER001566174000000000
func provideLog() []string {
	var examleLog = `
s|++++--++++533816|1518|+++traceId+++|1DCN001TRANSFER001571632753000005|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|dts001|INS001|DTS_AGENT_REGISTER|0||
s|++++--++++545491|1320|+++traceId+++|1DCN001WITHDRAW001571632753000002|1DCN001TRANSFER001571632753000006|ORG001|AZ0001|DCN001|NODE00000001|dts001|INS001|DTS_AGENT_ENLIST|0||
s|++++--++++561596|765|+++traceId+++|1DCN001DEPT0100001571632753000002|1DCN001TRANSFER001571632753000007|ORG001|AZ0001|DCN001|NODE00000001|dts001|INS001|DTS_AGENT_ENLIST|0||
c|++++--++++576091|6561|+++traceId+++|100000000000000001571632753000003|1DCN001TRANSFER001571632753000008|ORG001|AZ0001|DCN001|NODE00000001|dts001|INS001|DTS_CLIENT_WITHDRAW_CONFIRM|0||
c|++++--++++575986|6704|+++traceId+++|100000000000000001571632753000004|1DCN001TRANSFER001571632753000008|ORG001|AZ0001|DCN001|NODE00000001|dts001|INS001|DTS_CLIENT_DEPOSIT_CONFIRM|0||
s|++++--++++573600|9852|+++traceId+++|1DCN001TRANSFER001571632753000008|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|dts001|INS001|DTS_AGENT_TRY_RESULT_REPORT|0||
c|++++--++++532237|4341|+++traceId+++|1DCN001TRANSFER001571632753000005|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|transfer|TSF001|DTS_AGENT_REGISTER|0||
c|++++--++++537598|15287|+++traceId+++|1DCN001TRANSFER001571632753000006|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|transfer|TSF001|T_WITHDRAW|0||
c|++++--++++553744|14173|+++traceId+++|1DCN001TRANSFER001571632753000007|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|transfer|TSF001|T_DEPOSIT|0||
c|++++--++++569042|15894|+++traceId+++|1DCN001TRANSFER001571632753000008|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|transfer|TSF001|DTS_AGENT_TRY_RESULT_REPORT|0||
c|++++--++++543893|4710|+++traceId+++|1DCN001WITHDRAW001571632753000002|1DCN001TRANSFER001571632753000006|ORG001|AZ0001|DCN001|NODE00000001|withdraw|WTD001|DTS_AGENT_ENLIST|0||
s|++++--++++540691|10667|+++traceId+++|1DCN001TRANSFER001571632753000006|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|withdraw|WTD001|T_WITHDRAW|0||
s|++++--++++577749|3020|+++traceId+++|100000000000000001571632753000003|1DCN001TRANSFER001571632753000008|ORG001|AZ0001|DCN001|NODE00000001|withdraw|WTD001|DTS_CLIENT_WITHDRAW_CONFIRM|0||
c|++++--++++560024|4450|+++traceId+++|1DCN001DEPT0100001571632753000002|1DCN001TRANSFER001571632753000007|ORG001|AZ0001|DCN001|NODE00000001|deposit|DPT001|DTS_AGENT_ENLIST|0||
s|++++--++++556531|10059|+++traceId+++|1DCN001TRANSFER001571632753000007|+++traceId+++|ORG001|AZ0001|DCN001|NODE00000001|deposit|DPT001|T_DEPOSIT|0||
s|++++--++++577727|3312|+++traceId+++|100000000000000001571632753000004|1DCN001TRANSFER001571632753000008|ORG001|AZ0001|DCN001|NODE00000001|deposit|DPT001|DTS_CLIENT_DEPOSIT_CONFIRM|0||`

	result := make([]string, 0)
	now := time.Now().UnixNano() / 1000000000
	examleLog = strings.ReplaceAll(examleLog, "++++--++++", strconv.Itoa(int(now)))
	generateUUID, _ := uuid.GenerateUUID()
	examleLog = strings.ReplaceAll(examleLog, "+++traceId+++", generateUUID)
	//now := time.Now()
	for _, v := range strings.Split(examleLog, "\n") {
		if v == "" {
			continue
		}
		result = append(result, v)
	}
	return result
}
