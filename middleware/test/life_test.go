/*
 @Desc

 @Date 2020-04-02 13:07
 @Author yinjk
*/
package main

import (
	"fmt"
	"testing"
)

func Test1(_ *testing.T) {
	fmt.Println(red, "hello word", reset)
	fmt.Println(yellow, "hello word", reset)
	fmt.Print("\u001b[1000A")
	fmt.Print("\u001b[1000D")
	fmt.Println(green, "ni hao", reset)
}
