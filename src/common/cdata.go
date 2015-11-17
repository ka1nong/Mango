package common

import _ "github.com/go-sql-driver/mysql"
import "database/sql"
import "fmt"

type cData struct {
	db         *sql.DB
	stock_name string
}

//idata func-----------------------
func (data *cData) GetData(count int) {
}

func (data *cData) Close() {
	data.db.Close()
}

const (
	DATE  = "date"
	OPEN  = "open"
	HIGHT = "hight"
	CLOSE = "close"
	LOW   = "low"
)

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
	rows, err := data.db.Query("select count(*) from " + data.stock_name)
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

func (data *cData) InsertMainData(codes, names []string, address []int) error {
	stmt, err := data.db.Prepare("INSERT INTO " + MAIN_TABLE + "(code, name, address) VALUES(?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i := 0; i < len(codes); i++ {
		_, err = stmt.Exec(codes[i], names[i], address[i])
		if err != nil {
			continue
		}
	}
	return nil
}

//------------------------------------end---------------------------------------
