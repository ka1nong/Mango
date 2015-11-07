package download

type DownloadMgr struct {
}

var instance *DownloadMgr

func Instance() *DownloadMgr {
	if instance == nil {
		instance = new(DownloadMgr)
	}
	return instance
}

func (mgr *DownloadMgr) Download(url string) (localPath string, err error) {
	return localPath, nil
}
