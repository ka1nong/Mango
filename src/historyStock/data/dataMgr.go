package historyStock

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}

func walkDir(dirPth, suffix string) (files []string, err error) {
	if isDirExists(dirPth) == false {
		return nil, err
	}
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

func createDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", USER+":"+PASSWORD+"@/"+DATABASE)
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		db.Close()
		return nil, err
	}
	return db, nil
}

func parseFileFromCSV(stocksText string, filename string) error {
	db, err := createDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	//什么鬼，语法不能插入change字段，会报错。我操。只能把change改成hange
	stock_database := filename + "(id INTEGER PRIMARY KEY AUTO_INCREMENT, date DATE, open DOUBLE, hight DOUBLE, low DOUBLE, close DOUBLE,hange DOUBLE,volume DOUBLE,money DOUBLE,traded_market_value DOUBLE,market_value DOUBLE, turnover DOUBLE, adjust_price DOUBLE, report_type DATETIME,report_date DATETIME, PE_TTM DOUBLE,PS_TTM DOUBLE,PC_TTM DOUBLE,PB DOUBLE)"
	_, err = db.Exec("create table if not exists " + stock_database)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO " + filename + "(date, open, hight, low, close,hange,volume,money,traded_market_value,market_value, turnover, adjust_price, report_type,report_date, PE_TTM,PS_TTM,PC_TTM,PB) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	stockInfos := strings.Split(stocksText, "\n")
	stockInfos = stockInfos[1:]

	for _, v := range stockInfos {
		infos := strings.Split(v, ",")
		if len(infos) != 19 {
			return err
		}
		param1 := infos[1]
		param2, err := strconv.ParseFloat(infos[2], 64)
		if err != nil {
			return err
		}
		param3, err := strconv.ParseFloat(infos[3], 64)
		if err != nil {
			return err
		}
		param4, err := strconv.ParseFloat(infos[4], 64)
		if err != nil {
			return err
		}
		param5, err := strconv.ParseFloat(infos[5], 64)
		if err != nil {
			return err
		}
		param6, err := strconv.ParseFloat(infos[6], 64)
		if err != nil {
			return err
		}
		param7, err := strconv.ParseFloat(infos[7], 64)
		if err != nil {
			return err
		}
		param8, err := strconv.ParseFloat(infos[8], 64)
		if err != nil {
			return err
		}
		param9, err := strconv.ParseFloat(infos[9], 64)
		if err != nil {
			return err
		}
		param10, err := strconv.ParseFloat(infos[10], 64)
		if err != nil {
			return err
		}
		param11, err := strconv.ParseFloat(infos[11], 64)
		if err != nil {
			return err
		}
		param12, err := strconv.ParseFloat(infos[12], 64)
		if err != nil {
			return err
		}
		param13 := infos[13]
		param14 := infos[14]
		param15, err := strconv.ParseFloat(infos[15], 64)
		if err != nil {
			return err
		}
		param16, err := strconv.ParseFloat(infos[16], 64)
		if err != nil {
			return err
		}
		param17, err := strconv.ParseFloat(infos[17], 64)
		if err != nil {
			return err
		}
		param18, err := strconv.ParseFloat(infos[18], 64)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15, param16, param17, param18)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}
	return nil
}

func parseDapanFileFromCSV(stocksText string, filename string) error {
	db, err := createDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	//什么鬼，语法不能插入change字段，会报错。我操。只能把change改成hange
	stock_database := filename + "(id INTEGER PRIMARY KEY AUTO_INCREMENT, date DATE, open DOUBLE, close DOUBLE, low DOUBLE, hight DOUBLE,volume DOUBLE,money DOUBLE,hange DOUBLE)"
	_, err = db.Exec("create table if not exists " + stock_database)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO " + filename + "(date, open, close,low,hight,volume,money,hange) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	stockInfos := strings.Split(stocksText, "\n")
	stockInfos = stockInfos[1:]

	for _, v := range stockInfos {
		infos := strings.Split(v, ",")
		if len(infos) != 19 {
			return err
		}
		param1 := infos[1]
		param2, err := strconv.ParseFloat(infos[2], 64)
		if err != nil {
			return err
		}
		param3, err := strconv.ParseFloat(infos[3], 64)
		if err != nil {
			return err
		}
		param4, err := strconv.ParseFloat(infos[4], 64)
		if err != nil {
			return err
		}
		param5, err := strconv.ParseFloat(infos[5], 64)
		if err != nil {
			return err
		}
		param6, err := strconv.ParseFloat(infos[6], 64)
		if err != nil {
			return err
		}
		param7, err := strconv.ParseFloat(infos[7], 64)
		if err != nil {
			return err
		}
		param8, err := strconv.ParseFloat(infos[8], 64)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(param1, param2, param3, param4, param5, param6, param7, param8)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}
	return nil
}

func parseFiles(stockfiles []string, isDaPan bool) error {
	db, err := createDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	stockList := string("stockList")
	stock_database := stockList + "(code  INTEGER  PRIMARY KEY, address VARCHAR(4))"
	_, err = db.Exec("create table if not exists " + stock_database)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO " + stockList + "(code, address) VALUES(?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	for _, fileNamePath := range stockfiles {
		fileNameList := strings.Split(fileNamePath, "/")
		fileName := fileNameList[len(fileNameList)-1]
		address := fileName[:2]
		codeStr := fileName[2 : len(fileName)-4]
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = stmt.Exec(code, address)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buf, err := ioutil.ReadFile(fileNamePath)
		if err != nil {
			continue
		}
		if isDaPan {
			err = parseDapanFileFromCSV(string(buf), fileName[:len(fileName)-4])
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			err = parseFileFromCSV(string(buf), fileName[:len(fileName)-4])
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	return nil
}

func StartLoadData() error {
	go func() {
		stockfiles, err := walkDir("/home/huangchen/stocks/stock data", ".csv")
		if err != nil {
			fmt.Println(err)
		}
		err = parseFiles(stockfiles, true)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		stockfiles, err := walkDir("/mnt/stocks/stock data", ".csv")
		if err != nil {
			fmt.Println(err)
		}
		err = parseFiles(stockfiles, false)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}
