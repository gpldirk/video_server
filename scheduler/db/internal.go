package db

import (
	"log"
	_ "github.com/go-sql-driver/mysql"
)

// 读取video 删除记录到data chan
func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmt, err := dbConn.Prepare("select video_id from video_del_rec limit ?")
	var videoIds []string
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return videoIds, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(count)
	if err != nil {
		log.Printf("Db query err: %s", err.Error())
		return videoIds, err
	}

	for rows.Next() {
		var videoId string
		if err = rows.Scan(&videoId); err != nil {
			log.Printf("Db scan err: %s", err.Error())
			return videoIds, err
		}

		videoIds = append(videoIds, videoId)
	}

	return videoIds, nil
}

// 删除video成功之后需要删除videoId在DB中的记录
func DelVideoDeletionRecord(videoId string) error {
	stmt, err := dbConn.Prepare("delete from video_del_rec where video_id=?")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
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
		log.Printf("Db delete err: %s", err.Error())
		return err
	}
}


