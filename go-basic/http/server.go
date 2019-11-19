package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// w表示response对象，返回给客户端的内容都在对象里处理
// r表示客户端请求对象，包含了请求头，请求参数等等
func times(w http.ResponseWriter, r *http.Request) {
	unixNano := time.Now().UnixNano
	fmt.Println(fmt.Sprintf("%dns  %dms", unixNano(), unixNano()/1000000))
}

func main() {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/time", times)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9095", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
