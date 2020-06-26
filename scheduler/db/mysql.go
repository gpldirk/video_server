package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	dbConn *sql.DB
	err error
)

func init()  {
	dbConn, err = openConn()
}

func openConn() (*sql.DB, error) {
	dbConn, err := sql.Open("mysql", "root:14121314@tcp(localhost:3306)/video_server?charset=utf8")
	if err != nil {
		log.Printf("Connect to db err: %s", err.Error())
		return nil, err
	}
	err = dbConn.Ping()
	if err != nil {
		log.Printf("Ping db err: %s", err.Error())
		return nil, err
	}

	return dbConn, nil
}




