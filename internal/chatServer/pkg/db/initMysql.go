package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	config:="root:123456@tcp(127.0.0.1:3306)/todolist?parseTime=true"
	db,err:=sql.Open("mysql",config)
	if err!= nil {
		return err
	}
	db.SetConnMaxLifetime(60)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	err=db.Ping()
	if err!= nil {
		return err
	}
	DB=db
	return nil
}
