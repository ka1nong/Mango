package main

import "fmt"
import "stockmanger"
import "bufio"
import "strconv"
import "os"

func main() {

	for {
		fmt.Println("//-------------------")
		fmt.Println("//1. 更新库-------------")
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
			stockmgr := stockmanger.Instance()
			go stockmgr.Start()

		default:
			fmt.Println("chose error!")

		}
	}

	return
}
