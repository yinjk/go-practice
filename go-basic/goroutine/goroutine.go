package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var queen = make(chan string)

func main() {
	go provider(queen)
	var wait sync.WaitGroup
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			for msg := range queen {
				handMsg(msg)
			}
			wait.Done()
		}()
	}
	wait.Wait()
}

func handMsg(msg string) {
	time.Sleep(time.Millisecond * 100)
	fmt.Println(msg)
}

func provider(in chan string) {
	for i := 0; i < 10000; i++ {
		in <- strconv.Itoa(i)
	}
	close(queen)
}
