package main

import "fmt"
import "bufio"
import "strconv"
import "os"
import "historyStock/data"

func main() {

	a := true
	for {
		fmt.Println("//------------------")
		fmt.Println("//1.load historyStock------")
		fmt.Println("//------------------")
		fmt.Println("//------------------")
		fmt.Println("//0. 退出-----------")
		fmt.Println("//------------------")
		fmt.Println("please enter key.")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadBytes('\n')
		chose, err := strconv.Atoi(string("1"))
		if err != nil {
			fmt.Println(err)
			return
		}
		if a == false {
			return 
		}
		a = false
		switch chose {
		case 0:
			return
		case 1:
			err = historyStock.StartLoadData()
			if err != nil {
				fmt.Println(err)
			}
			break
		default:
			fmt.Println("chose error!")

		}
	}

	return
}
