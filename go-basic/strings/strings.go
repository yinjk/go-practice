/*
 @Desc

 @Date 2020-03-18 15:34
 @Author yinjk
*/
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	useBuffer()
	useBuilder()
}
func useBuilder() {
	var builder strings.Builder
	builder.WriteString("hello")
	builder.WriteString(" ")
	builder.WriteString("world")
	builder.WriteString("!")
	fmt.Println(builder.String())
}
func useBuffer() {
	var buffer bytes.Buffer
	buffer.WriteString("hello")
	buffer.WriteString(" ")
	buffer.WriteString("world")
	buffer.WriteString("!")
	fmt.Println(buffer.String())
}
