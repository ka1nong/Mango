package common

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 2500)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func (mgr *stockMgr) parseFileFromCSV(v string) (infos []map[string]string, err error) {
	return nil, err
}

func (mgr *stockMgr) initFromLocalFiles() {
	files, err := WalkDir("../stocks", ".csv")
	if err != nil {
		return
	}
	cdata, err := mgr.open(MAIN_TABLE)
	for _, v := range files {
		mgr.parseFileFromCSV(v)
	}
}
