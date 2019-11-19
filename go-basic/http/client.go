package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	flag.Parse()
	ip := flag.Arg(0)
	client := http.DefaultClient
	unixNano := time.Now().UnixNano
	fmt.Println(fmt.Sprintf("%dns  %dms", unixNano(), unixNano()/1000000))
	url := "http://" + ip + ":9095/time"
	if _, err := client.Get(url); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(fmt.Sprintf("ping %s success", ip))
	}
}
