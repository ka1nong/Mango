package main

import "fmt"
import "bufio"
import "strconv"
import "os"

func main() {

	for {
		fmt.Println("//------------------")
		fmt.Println("//1.峰值分析法------")
		fmt.Println("//------------------")
		fmt.Println("//------------------")
		fmt.Println("//0. 退出-----------")
		fmt.Println("//------------------")
		fmt.Println("please enter key.")
		reader := bufio.NewReader(os.Stdin)
		ch, _ := reader.ReadBytes('\n')
		chose, err := strconv.Atoi(string(ch))
		if err != nil {
			fmt.Println(err)
			return
		}
		switch chose {
		case 0:
			return
		case 1:

			break;
		default:
			fmt.Println("chose error!")

		}
	}

	return
}
