package main

import "fmt"
import "stockmanger"
import "bufio"
import "strconv"
import "os"

func main() {

	for {
		fmt.Println("//-------------------")
		fmt.Println("//1. 建库-------------")
		fmt.Println("//2. 峰值分析法--------")
		fmt.Println("//-------------------")
		fmt.Println("//0. 退出-------------")
		fmt.Println("//-------------------")
		fmt.Println("please enter key.")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadBytes('\n')
		chose, err := strconv.Atoi(string(input[0]))
		if err != nil {
			fmt.Println(err)
			return
		}
		switch chose {
		case 0:
			return
		case 1:
			fmt.Println("this operate is important，enter key 1 continue")
			reader = bufio.NewReader(os.Stdin)
			input, _ = reader.ReadBytes('\n')
			value, err := strconv.Atoi(string(input[0]))
			if err != nil {
				fmt.Println(err)
			} else if value == 1 {
				stockmgr := stockmanger.NewStockMgr()
				go stockmgr.Start()
			} else {
				fmt.Println("enter key error, return main menu")
			}

		default:
			fmt.Println("chose error!")

		}
	}

	return
}
