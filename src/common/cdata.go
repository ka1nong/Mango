package common

import _ "github.com/go-sql-driver/mysql"
import "database/sql"
import "fmt"

type cData struct {
	db         *sql.DB
	stock_name string
}

//idata func-----------------------
func (data *cData) GetData(infos []string, count int) (datas []map[string]string, err error) {
	if count == 0 || len(infos) == 0 {
		return
	}
	return nil, Error("param error")
}

func (data *cData) Close() {
	data.db.Close()
}

func (data *cData) InsertData(info map[string]interface{}) error {
	//时间，开盘、最高、收盘、最低,
	date := info[DATE].(int)
	if date == 0 {
		return Error("insert fmt error")
	}
	rows, err := data.db.Query("select * from "+data.stock_name+" where date = ?", date)
	if err != nil {
		return err
	}
	defer rows.Close()
	/*
		tx, _ := data.db.Begin()
		//IF (SELECT * FROM ipstats WHERE date='192.168.0.1)' {
		//	    UPDATE ipstats SET clicks=clicks+1 WHERE ip='192.168.0.1';
		//	} else {
		//		 INSERT INTO ipstats (ip, clicks) VALUES ('192.168.0.1', 1);
		//}
		stmt, err := data.db.Prepare("INSERT INTO " + data.stock_name + "(date, open, hight, close, low) VALUES(?, ?, ?,?, ?)")
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		//for i := 0; i < len(date); i++ {
		//	_, err = stmt.Exec(date[i], open[i], hight[i], close[i], low[i])
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		//}
		err = tx.Commit() //err := Tx.Rollback()
		return nil
	*/
	return nil
}

func (data *cData) GetInfoCount() int {
	mgr := getStockMgr()
	return mgr.GetInfoCount(data)
}

//-----------------------end----------------------------------
//main func
func (data *cData) GetStockCount() int {
	rows, err := data.db.Query("select count(*) from " + MAIN_TABLE)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			return 0
		}
	}
	return id
}

func (data *cData) GetRandomMainData() (stock map[string]string, err error) {
	return nil, err
}

//------------------------------------end---------------------------------------
