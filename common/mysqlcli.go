package common

import (
	"database/sql"
	//"fmt"
	"../config"
	"../logs"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
	err  error
)

func DB() (*sql.DB, error) {
	Refresh()
	return db, err
}

func Refresh() {
	once.Do(func() {
		db, err = sql.Open("mysql", config.MYSQL_CONN_STR+"?charset=utf8")
		if err != nil {
			//logs.Log.Error("Mysql：%v\n", err)
			return
		}
		db.SetMaxOpenConns(config.MYSQL_CONN_CAP)
		db.SetMaxIdleConns(config.MYSQL_CONN_CAP)
	})
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		//logs.Log.Error("Mysql：%v\n", err)
		logs.Log.Debug("url: %v", err)
	}
}
