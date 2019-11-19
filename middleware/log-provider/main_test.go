//@Desc
//@Date 2019-10-15 10:12
//@Author yinjk
package main

import (
	"fmt"
	"os"
	"testing"
)

func TestExists(_ *testing.T) {
	fmt.Println(filePath, c, t)

	if !Exists(filePath) {
		if err := os.MkdirAll(filePath, 0666); err != nil {
			panic(err)
		}
	}

	if f, err := os.OpenFile(filePath+"apm.log", os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		panic(err)
	} else {
		if info, err := f.Stat(); err == nil {
			fmt.Println(info.Size() / 1000000)
		}
	}
}
