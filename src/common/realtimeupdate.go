package common

import "download"
import "io/ioutil"

func startRealTimeUpdate() error {
	mgr := getStockMgr()
	for i := 0; i < 50; i++ {
		cdata, err := mgr.open(MAIN_TABLE)
		if err != nil {
			return err
		}
		defer cdata.Close()
		v := NewRealTime()
		go v.realTimeUpdate(cdata)
	}
	return nil
}

type realTime struct {
}

func NewRealTime() (re *realTime) {
	return &realTime{}
}
func (mgr *realTime) isStockUpdateTime() bool {
	return true
}

//todo:闲下来的时候进行数据补齐
func (mgr *realTime) realTimeUpdate(cdata *cData) {
	if mgr.isStockUpdateTime() {
		dwMgr := download.Instance()
		localPath, err := dwMgr.Download( /*mgr.stockUrl +*/ string("sh") + "12")
		_, err = ioutil.ReadFile(localPath)
		if err != nil {
			//todo: retry self
		}
		//parse insert
		//addrecord self time
	}
}
