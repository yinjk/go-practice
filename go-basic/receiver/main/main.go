package main

import (
	"fmt"
	"go-practice/go-basic/receiver"
)

func main() {
	//p:=  receiver.NewPerson("hello")
	//fmt.Println(p.Name())
	//fmt.Println(p.Child().Name())

	var rec *receiver.Person
	fmt.Println(rec.Name())
}
