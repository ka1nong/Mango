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

func walkDir(dirPth, suffix string) (files []string, err error) {
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
	return db, err
}

func parseFileFromCSV(filename string) error {
	db, err := createDatabase()
	if err != nil {
		return err
	}
	db.Close()

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	stock_database := filename + "(id INTEGER PRIMARY KEY AUTOINCREMENT, date DATE, open DOUBLE, hight DOUBLE, low DOUBLE, close DOUBLE, change DOUBLE, volume DOUBLE, money DOUBLE, traded_market_value DOUBLE, market_value DOUBLE, turnover DOUBLE, adjust_price DOUBLE, report_type DOUBLE, report_date DOUBLE, PE_TTM DOUBLE, PS_TTM DOUBLE, PC_TTM DOUBLE, PB DATE)"
	_, err = db.Exec("create table if not exists " + stock_database)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO " + filename + "(date, open, hight, low, close,change,volume,money,traded_market_value,market_value, turnover, adjust_price, report_type,report_date, PE_TTM,PS_TTM,PC_TTM,PB) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	stocksText := string(buf)
	stockInfos := strings.Split(stocksText, "\n")
	stockInfos = stockInfos[1:]

	for _, v := range stockInfos {
		infos := strings.Split(v, ",")
		_, err = stmt.Exec(code, address)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}
	return nil
}

func parseFiles(stockfiles []string) error {
	db, err := createDatabase()
	if err != nil {
		return err
	}
	db.Close()

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

	for _, fileName := range stockfiles {
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
		err = parseFileFromCSV(fileName)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}

func parseDapanFileFromCSV(db *sql.DB, fileName string) error {
	return nil
}

func StartLoadData() error {
	go func() {
		stockfiles, err := walkDir("../stocks/stock data", ".csv")
		if err != nil {
			fmt.Println(err)
		}
		err = parseFiles(stockfiles)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		stockfiles, err := walkDir("../stocks/index data", ".csv")
		if err != nil {
			fmt.Println(err)
		}
		db, err := createDatabase()
		if err != nil {
			fmt.Println(err)
		}
		db.Close()
		for _, v := range stockfiles {
			err = parseDapanFileFromCSV(db, v)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	return nil
}
