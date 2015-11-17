package main

import "fmt"
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
		_, _ = reader.ReadBytes('\n')
		chose, err := strconv.Atoi(string("1"))
		if err != nil {
			fmt.Println(err)
			return
		}
		switch chose {
		case 0:
			return
		case 1:

		default:
			fmt.Println("chose error!")

		}
	}

	return
}
