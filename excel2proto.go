package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("")
	if nil != err {
		fmt.Println("this is fail")
		return
	}
	f.Close()

	fmt.Println("--------------------")
}
