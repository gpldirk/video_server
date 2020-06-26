package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func AddVideoDeletionRecord(videoId string) error {
	stmt, err := dbConn.Prepare("insert into video_del_rec (video_id) values(?)")
	if err != nil {
		log.Printf("Db Connection err: %s", err.Error())
		return err
	}
	stmt.Close()

	rf, err := stmt.Exec(videoId)
	if err != nil {
		log.Printf("Db execution err: %s", err.Error())
		return err
	}
	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("Add video deletion record failed")
		return err
	}
}