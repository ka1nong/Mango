package peakanalysis

import "common"
import "fmt"

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
	}
	return nil
}
