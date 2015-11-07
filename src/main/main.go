package main

import "fmt"
import "stockmanger"

func main() {

	fmt.Println("main start")
	stockmgr := stockmanger.NewStockMgr()
	err := stockmgr.Start()
	if err != nil {
		fmt.Println("stock mgr return  error!")
	}

	fmt.Println("main end!")
}
