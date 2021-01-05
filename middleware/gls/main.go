/*
 @Desc

 @Date 2020-12-30 14:03
 @Author inori
*/
package main

import (
	"fmt"
	"sync"

	"github.com/jtolds/gls"
)

var (
	mgr = gls.NewContextManager()
	key = gls.GenSym()
)

func main() {
	MyLog := func() {
		if requestId, ok := mgr.GetValue(key); ok {
			fmt.Println("My request id is:", requestId)
		} else {
			fmt.Println("No request id found")
		}
	}

	mgr.SetValues(gls.Values{key: "12345"}, func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			MyLog()
		}()
		wg.Wait()
		wg.Add(1)
		gls.Go(func() {
			defer wg.Done()
			MyLog()
		})
		wg.Wait()
	})
}
