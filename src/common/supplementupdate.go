package common

func isNeedSupplement(cdata *cData, stockCode int) (isNedd bool, err error) {
	return false, err
}

type supplement struct {
}

func Newsupplement() (data *supplement) {
	return &supplement{}
}

func (mgr *supplement) start() error {
	return nil
}
