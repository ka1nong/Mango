package  main 
package stockmanger

import "fmt"

func main() {
    stockmgr := stockmanger.NewStockMgr();
    if err :=stockmgr.start() err != nil {
    	fmt.Println("start error");
    }
    fmt.Println("start success");
}