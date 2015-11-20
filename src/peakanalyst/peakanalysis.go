package peakanalysis

import "common"
import "fmt"
import "time"

/*
  k := time.Now()
  //一天之前
  d, _ := time.ParseDuration("-24h")
  fmt.Println(k.Add(d))
  //一周之前
  fmt.Println(k.Add(d * 7))
  //一月之前
  fmt.Println(k.Add(d * 30))
*/

type PeakAnalysis struct {
}

func (analysis *PeakAnalysis) Start() error {
	iMainData, err := common.GetIMainData()
	if err != nil {
		return err
	}
	defer iMainData.Close()
	stocks, err := iMainData.GetAllStockInfo()
	if err != nil {
		return err
	}
	for _, stock := range stocks {
		code := stock[common.STOCK_CODE]
		iData, err := common.GetIData(code.(int))
		if err != nil {
			fmt.Println("open stock error,the code is%d", code)
			continue
		}
		defer iData.Close()
		analysis.analysisInfoByDay(iData)

	}
	return nil
}

func (analysis *PeakAnalysis) analysisInfoByDay(iData common.IData) error {
	infos := []string{common.STOCKINFO_DATE, common.STOCKINFO_CLOSE}
	beginTime := time.Now()
	d, err := time.ParseDuration("d * 31")
	if err != nil {
		fmt.Println(err)
		return err
	}
	endTime := beginTime.Add(d)
	density := string("-24h")
	datas, err := iData.GetData(infos, beginTime.Unix(), endTime.Unix(), density)
	if err != nil {
		fmt.Println("get stock data error: %s", err)
		return err
	}

	init := true
	var maxPrice float32
	var minPrice float32
	var currentPrice float32
	for _, stock := range datas {
		value := stock[common.STOCKINFO_CLOSE]
		price := value.(float32)
		if init {
			init = false
			maxPrice = price
			minPrice = price
		}
		if price < minPrice {
			minPrice = price
		}
		if maxPrice < price {
			maxPrice = price
		}
		//todo:
		fmt.Println(currentPrice)

	}
	return nil
}
