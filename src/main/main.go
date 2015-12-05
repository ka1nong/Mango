package main

import "fmt"
import "historyStock/data"
import "onlineStock/data"
import "flag"

func main() {
	//	mycontent := " my dear令"
	//	email := NewEmail("546958900@qq.com;546958900@qq.com;", "test golang email", mycontent)
	//	effdrr := SendEmail(email)
	//	if effdrr != nil {
	//		fmt.Println(effdrr)
	//	}
	//for {
	/*	fmt.Println("//-----------------------------------------------")
		fmt.Println("//1.nohup historyStock-----------------")
		fmt.Println("//2.nohup onlineStock------------------")
		fmt.Println("//-----------------------------------------------")
		fmt.Println("//0. 退出-------------------------------------")
		fmt.Println("//----------------------------------------------")
		fmt.Println("please enter key.")
		reader := bufio.NewReader(os.Stdin)
		ch, _ := reader.ReadBytes('\n')
		chose, err := strconv.Atoi(string(ch[0]))
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	id := flag.Int("id", 0, "id")
	flag.Parse()
	*id = 2
	switch *id {
	case 0:
		return
	case 1:
		err := historyStock.StartLoadData()
		if err != nil {
			fmt.Println(err)
		}
		return
	case 2:
		err := onlineStock.StartLoadOnlineData()
		if err != nil {
			fmt.Println(err)
		}
		return
	default:
		fmt.Println("chose error!")

	}
	//}

	return
}
