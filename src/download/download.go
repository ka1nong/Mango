package download

import "net/http"
import "os"
import "io/ioutil"
import "fmt"
import "strconv"

type Error string

func (e Error) Error() string {
	return string(e)
}

type DownloadMgr struct {
}

var instance *DownloadMgr

func Instance() *DownloadMgr {
	if instance == nil {
		instance = new(DownloadMgr)
		instance.init()
	}
	return instance
}

func (mgr *DownloadMgr) Download(url string) (localPath string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return localPath, Error("download error")
	}
	defer resp.Body.Close()

	localPath, _ = mgr.getOnlyName()

	fout, err := os.Create(localPath)

	if err != nil {
		fmt.Println(err.Error())
		return localPath, Error("creat localPath error")
	}
	defer fout.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fout.WriteString(string(body))
	return localPath, nil
}

func (mgr *DownloadMgr) getOnlyName() (localPathName string, err error) {
	filePath := "../"
	count := 1
	for {
		if _, err := os.Stat(filePath + strconv.Itoa(count)); err == nil {
			//存在该文件
			count++
		} else {
			break
		}
	}
	return filePath + strconv.Itoa(count), nil
}

func (mgr *DownloadMgr) init() {
	filePath := "../"
	count := 1
	for {
		if _, err := os.Stat(filePath + strconv.Itoa(count)); err == nil {
			os.Remove(filePath + strconv.Itoa(count))
			count++
		} else {
			break
		}
	}
	return
}
