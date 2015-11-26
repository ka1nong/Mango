package main

import "fmt"
import "bufio"
import "strconv"
import "os"
import _ "historyStock/data"
import "onlineStock/data"

func main() {
	//	mycontent := " my dear令"
	//	email := NewEmail("546958900@qq.com;546958900@qq.com;", "test golang email", mycontent)
	//	effdrr := SendEmail(email)
	//	if effdrr != nil {
	//		fmt.Println(effdrr)
	//	}
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
		a = true
		switch chose {
		case 0:
			return
		case 1:
			err = onlineStock.StartLoadOnlineData()
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
