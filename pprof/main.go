/*
 @Desc

 @Date 2020-12-08 16:54
 @Author inori
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	tz *time.Location
)

func main() {
	go func() {
		for {
			LocalTz()

			doSomething([]byte(`{"a": 1, "b": 2, "c": 3}`))
		}
	}()

	fmt.Println("start api server...")
	panic(http.ListenAndServe(":8080", nil))
}

func doSomething(s []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(s, &m)
	if err != nil {
		panic(err)
	}

	s1 := make([]string, 100)
	s2 := ""
	//var buffer bytes.Buffer
	for i := 0; i < 100; i++ {
		s1[i] = string(s)
		s2 += string(s)
		//buffer.Write(s)
	}
}

func LocalTz() *time.Location {
	if tz == nil {
		tz, _ = time.LoadLocation("Asia/Shanghai")
	}
	return tz
}
