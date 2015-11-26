package onlineStock

import "download"
import "fmt"
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"mahonia"
	"strconv"
	"strings"
	"time"
)

type stockInfo struct {
	name    string
	url     string
	address int // 0 深圳 1 上海
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

func updateMainDatabase() (stocks map[string]stockInfo, err error) {
	db, err := createDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	stocks, err = loadStockMap()
	if err != nil {
		return nil, err
	}

	stockList := string("stockList")
	stock_database := stockList + "(code VARCHAR(40)  PRIMARY KEY ,name VARCHAR(40) NOT NULL, address INTEGER)"
	_, err = db.Exec("create table if not exists " + stock_database)
	if err != nil {
		return nil, err
	}
	//建立本地库
	codes := make([]string, len(stocks))
	names := make([]string, len(stocks))
	address := make([]int, len(stocks))
	i := 0
	for key, info := range stocks {
		codes[i] = key
		names[i] = info.name
		address[i] = info.address
		if len(codes[i]) != 0 && len(info.name) != 0 {
			i++
		}
	}
	_, err = InsertMainData(db, codes[:i], names[:i], address[:i])
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func InsertMainData(db *sql.DB, codes, names []string, address []int) (count int, err error) {
	stmt, err := db.Prepare("INSERT INTO " + MAIN_TABLE + "(code, name, address) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return count, err
	}
	defer stmt.Close()
	for i := 0; i < len(codes); i++ {
		_, err = stmt.Exec(codes[i], names[i], address[i])
		if err != nil {
			fmt.Println(err)
			fmt.Println("insert main data error: %s, %s, %d", codes[i], names[i], address[i])
			continue
		}
		count++
	}
	return count, nil
}

func loadStockMap() (map[string]stockInfo, error) {
	dwMgr := download.Instance()
	localPath, err := dwMgr.Download("http://quote.eastmoney.com/stocklist.html#sz")
	if err != nil {
		fmt.Println("load remote stock map error")
		return nil, err
	}
	remoteStockMap, err := getStockMapByParseFile(localPath)
	if err != nil {
		fmt.Println("parse remote stock map error")
		return nil, err
	}
	return remoteStockMap, nil
}

func getStockMapByParseFile(localPath string) (stockMap map[string]stockInfo, err error) {
	buf, err := ioutil.ReadFile(localPath)
	if err != nil {
		return stockMap, Error("read remote stock map error")
	}

	enc := mahonia.NewEncoder("gbk")
	remoteStocksText := string(buf)
	//remoteStocksText = enc.ConvertString(remoteStocksText) //转换utf-8编码
	index := strings.Index(remoteStocksText, enc.ConvertString("股票代码查询一览表："))
	if index == -1 {
		return stockMap, Error("parse remote stock map error")
	}
	remoteStocksText = remoteStocksText[index:]
	index = strings.Index(remoteStocksText, string("<ul>"))
	if index == -1 {
		return stockMap, Error("parse remote stock map error")
	}
	shangHaiEndIndex := strings.Index(remoteStocksText, string("</ul>"))
	if shangHaiEndIndex == -1 || shangHaiEndIndex < index {
		return stockMap, Error("parse remote stock map error")
	}
	shangHaiStocksText := remoteStocksText[index:shangHaiEndIndex]

	remoteStocksText = remoteStocksText[shangHaiEndIndex+1:] //从加1开始查深圳股票，不然shenZhenEndIndex就是0
	index = strings.Index(remoteStocksText, string("<ul>"))
	if index == -1 {
		return stockMap, Error("parse remote stock map error")
	}
	shenZhenEndIndex := strings.Index(remoteStocksText, string("</ul>"))
	if shenZhenEndIndex == -1 || shenZhenEndIndex < index {
		return stockMap, Error("parse remote stock map error")
	}
	shenZhenStocksText := remoteStocksText[index:shenZhenEndIndex]

	shanghaiStocks, err := parseRemoteStockMapByUL(shangHaiStocksText)
	if err != nil {
		return stockMap, err
	}

	shenZhenStocks, err := parseRemoteStockMapByUL(shenZhenStocksText)
	if err != nil {
		return stockMap, err
	}

	stockMap = make(map[string]stockInfo)
	for key, info := range shanghaiStocks {
		info.address = 1
		stockMap[key] = info
	}
	//充分利用深圳股票的枚举是0
	for key, info := range shenZhenStocks {
		stockMap[key] = info
	}
	return stockMap, err
}

func parseRemoteStockMapByUL(stocksText string) (stockMap map[string]stockInfo, err error) {

	stockMap = make(map[string]stockInfo)

	parseli := func(litext string) (err error) {
		index := strings.Index(litext, string("href="))
		if index == -1 {
			return Error("parse remote li error")
		}
		litext = litext[index+len(string("href=")):]
		index = strings.Index(litext, string(">"))
		if index == -1 {
			return Error("parse remote li error")
		}
		url := litext[1 : index-1] //去除两边的分号
		litext = litext[index+1:]
		index = strings.Index(litext, string("("))
		if index == -1 {
			return Error("parse remote li error")
		}
		name := litext[:index]
		litext = litext[index+1:]
		index = strings.Index(litext, string(")"))
		if index == -1 {
			return Error("parse remote li error")
		}
		stockid := litext[:index]
		if len(url) == 0 || len(name) == 0 || len(stockid) == 0 {
			return Error("parse remote li error")
		}

		stockMap[stockid] = stockInfo{name, url, 0}
		return err
	}

	for {
		liBeginIndex := strings.Index(stocksText, string("<li>"))
		liEndIndex := strings.Index(stocksText, string("</li>"))
		if liBeginIndex == -1 || liEndIndex == -1 || liEndIndex < liBeginIndex {
			break
		}
		text := stocksText[liBeginIndex:liEndIndex]
		err := parseli(text)
		if err != nil {
			return stockMap, err
		}
		stocksText = stocksText[liEndIndex+len(string("</li>")):]
	}
	return stockMap, err
}

func isStockUpdateTime() bool {
	dwMgr := download.Instance()

	contextFunc := func() (text string, err error) {
		count := 0
		for { //随便拿个股票测试
			localPath, err := dwMgr.Download(string("http://hq.sinajs.cn/list=sh601006"))
			if count == 5000 {
				return string(""), Error("error")
			}
			if err != nil {
				count++
				continue
			}
			buf, err := ioutil.ReadFile(localPath)
			if err != nil {
				dwMgr.RemoveFile(localPath)
				return string(""), Error("error")
			}
			dwMgr.RemoveFile(localPath)
			return string(buf), nil
		}
		return string(""), nil
	}

	originText, err := contextFunc()
	if err != nil {
		//send email
		return false
	}
	time.Sleep(time.Second * 1) //------------------------------
	text, err := contextFunc()
	if err != nil {
		//send email
		return false
	}
	if strings.EqualFold(originText, text) {
		return false
	}
	return true
}

func CreateStockDataBase(stockStr string, db *sql.DB) error {
	stock_database := stockStr + "(code  INTEGER  PRIMARY KEY, open DOUBLE, yesterdayClose DOUBLE, currentPrice DOUBLE,  todayHight DOUBLE, todayLow DOUBLE,buyOne DOUBLE, sellOne DOUBLE,dealCount DOUBLE,dealPrice DOUBLE,buyOneCount DOUBLE,buyOnePrice DOUBLE,buyTwoCount DOUBLE,buyTwoPrice DOUBLE,buyThreeCount DOUBLE, buyThreePrice DOUBLE, buyFourCount DOUBLE,buyFourPrice DOUBLE,buyFiveCount DOUBLE, buyFivePrice DOUBLE, sellOneCount DOUBLE, sellOnePrice DOUBLE, sellTwoCount DOUBLE, sellTwoPrice DOUBLE,sellThreeCount DOUBLE, sellThreePrice DOUBLE, sellFourCount DOUBLE, sellFourPrice DOUBLE, sellFiveCount DOUBLE, sellFivePrice DOUBLE, currentDate DATE, currentTime  TIME)"
	_, err := db.Exec("create table if not exists " + stock_database)
	if err != nil {
		return err
	}
	return nil
}

func realTimeUpdate(code []string, address []int) {
	db, err := createDatabase()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var stockStr string
	for i, value := range code {
		if address[i] == 0 {
			stockStr = "sz" + value
		} else {
			stockStr = "sh" + value
		}
		err := CreateStockDataBase(stockStr, db)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	dwMgr := download.Instance()
	for {
		if isStockUpdateTime() {
			for i := 0; i < len(code); i++ {
				localPath, err := dwMgr.Download(string("http://hq.sinajs.cn/list=") + string("sh") + code[i])
				buf, err := ioutil.ReadFile(localPath)
				if err != nil {
					dwMgr.RemoveFile(localPath)
					fmt.Println(err)
				}
				ParseStockByUrl(stockStr, string(buf), db)
				dwMgr.RemoveFile(localPath)
			}
		}
	}
}

func ParseStockByUrl(stockStr string, data string, db *sql.DB) error {
	index := strings.Index(data, string("="))
	if index == -1 || index > len(data)-1 {
		return Error("parse stock by url error")
	}
	data = data[index+len(string("="))+1 : len(data)-2]
	stockInfos := strings.Split(data, string(","))
	infos := stockInfos[1:]

	if len(infos) < 32 {
		return Error("parse param error")
	}

	errHandle := func(param *float64, str string) error {
		v, err := strconv.Atoi(str)
		if err != nil {
			return nil
		}
		*param = float64(v)
		return nil
	}

	var params [29]float64
	for i := 0; i < 29; i++ {
		param, err := strconv.ParseFloat(infos[i], 64)
		if err != nil {
			err = errHandle(&param, infos[i])
			if err != nil {
				return err
			}
		}
		params[i] = param
	}

	stmt, err := db.Prepare("INSERT INTO " + stockStr + "( open, yesterdayClose, currentPrice, todayHight, todayLow,buyOne, sellOne,dealCount,dealPrice,buyOneCount,buyOnePrice,buyTwoCount,buyTwoPrice,buyThreeCount, buyThreePrice, buyFourCount,buyFourPrice,buyFiveCount, buyFivePrice, sellOneCount, sellOnePrice, sellTwoCount, sellTwoPrice,sellThreeCount, sellThreePrice, sellFourCount, sellFourPrice, sellFiveCount, sellFivePrice, currentDate, currentTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7], params[8], params[9], params[10], params[11], params[12], params[13], params[14], params[15], params[16], params[17], params[18], params[19], params[20], params[21], params[22], params[23], params[24], params[25], params[26], params[27], params[28], infos[29], infos[30])
	if err != nil {
		return err
	}
	return nil
}

func StartLoadOnlineData() error {
	stocks, err := updateMainDatabase()
	if err != nil {
		fmt.Println("updata main data base error")
		fmt.Println(err)
		return err
	}
	fmt.Println("updata main data success")
	codes := make([]string, len(stocks))
	address := make([]int, len(stocks))
	i := 0
	for key, info := range stocks {
		codes[i] = key
		address[i] = info.address
		if len(codes[i]) != 0 && len(info.name) != 0 {
			i++
		}
	}
	syncCount := i / 50
	j := 0
	for ; j < syncCount; j++ {
		go realTimeUpdate(codes[j*50:j*50+50], address[j*50:j*50+50])
	}
	if (i - j*50) != 0 {
		go realTimeUpdate(codes[j*50:i], address[j*50:i])
	}
	return nil
}
