package peakanalysis

import "common"

type PeakAnalysis struct {
}

func (analysis *PeakAnalysis) Start() error {
	iMainData, err := common.GetIMainData()
	if err != nil {
		return err
	}
	defer iMainData.Close()
	return nil
}
