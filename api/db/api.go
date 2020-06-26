package db

import (
	"database/sql"
	"github.com/video_server/api/model"
	"github.com/video_server/api/utils"
	"log"
	"time"
)

// user CRUD operation
func AddUserCredential(username, password string) error {
	// 将sql命令部分和数据部分分离开来，保证安全
	stmt, err := dbConn.Prepare("insert into users (login_name, pwd) values (?, ?)")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(username, password)
	if err != nil {
		log.Printf("Db execuation err: %s", err.Error())
		return err
	}
	if rows, err:= rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("Add user failed")
		return err
	}
}

func GetUserCredential(username string) (string, error) {
	stmt, err := dbConn.Prepare("select pwd from users where login_name=?")
	if err != nil {
		log.Printf("db connection err: %s", err.Error())
		return "", err
	}
	defer stmt.Close()

	var password string
	err = stmt.QueryRow(username).Scan(&password)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return password, nil
}

func DeleteUserCredential(username, password string) error {
	stmt, err := dbConn.Prepare("delete from users where login_name=? and pwd=?")
	if err != nil {
		log.Printf("db connection err: %s", err.Error())
		return err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(username, password)
	if err != nil {
		log.Printf("db execution err: %s", err.Error())
	}
	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("delete user failed")
		return err
	}
}

func GetUser(username string) (*model.User, error) {
	stmt, err := dbConn.Prepare("select id, pwd from users where login_name=?")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	user.Username = username
	err = stmt.QueryRow(username).Scan(&user.Id, &user.Password)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Db query err: %s", err.Error())
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}


// video CRUD operation
func AddNewVideo(authorId int, name string) (*model.VideoInfo, error) {
	videoId, err := utils.NewUUID()
	if err != nil {
		log.Printf("Generate uuid err: %s", err.Error())
		return nil, err
	}

	displayTime := time.Now().Format("Jan 02 2006, 15:04:05")
	stmt, err := dbConn.Prepare("insert into video_info (id, author_id, name, display_ctime) values (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(videoId, authorId, name, displayTime)
	if err != nil {
		log.Printf("Db execution err: %s", err.Error())
		return nil, err
	}

	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return &model.VideoInfo{Id: videoId, AuthorId: authorId, Name: name, DisplayCtime: displayTime}, nil
	} else {
		log.Printf("Add video failed")
		return nil, err
	}
}

func GetVideoInfo(id string) (*model.VideoInfo, error) {
	stmt, err := dbConn.Prepare("select id, author_id, name, display_ctime from video_info where id=?")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var video model.VideoInfo
	err = stmt.QueryRow(id).Scan(&video.Id, &video.AuthorId, &video.Name, &video.DisplayCtime)
	if err != nil {
		log.Printf("Query video err: %s", err.Error())
		return nil, err
	}
	return &video, nil
}

func ListVideoInfo(username string, from, to int) ([]*model.VideoInfo, error) {
	stmt, err := dbConn.Prepare(`select video_info.id, video_info.author_id, video_info.name, video_info.display_ctime 
			from video_info inner join users on video_info.author_id=users.id where users.login_name=?
			and video_info.create_time > FROM_UNIXTIME(?) and video_info.create_time <= FROM_UNIXTIME(?) 
			order by video_info.create_time DESC`)
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, from, to)
	if err != nil {
		log.Printf("Db query err: %s", err.Error())
		return nil, err
	}

	var videosInfo []*model.VideoInfo
	var id, name, displayCtime string
	var authorId int
	for rows.Next() {
		if err = rows.Scan(&id, &authorId, &name, &displayCtime); err != nil {
			log.Printf("Db next err:%s", err.Error())
			return nil, err
		}
		videosInfo = append(videosInfo, &model.VideoInfo{
			Id:           id,
			AuthorId:     authorId,
			Name:         name,
			DisplayCtime: displayCtime,
		})
	}
	return videosInfo, nil
}

func DeleteVideoInfo(id string) error {
	stmt, err := dbConn.Prepare("delete from video_info where id=?")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Delete video err: %s", err.Error())
		return err
	}
	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("Delete video failed")
		return nil
	}
}

// comments CREATE/GET operation
func AddNewComment(videoId string, authorId int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		log.Printf("Generate uuid err: %s", err.Error())
		return err
	}

	stmt, err := dbConn.Prepare("insert into comments (id, video_id, author_id, content) values (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(id, videoId, authorId, content)
	if err != nil {
		log.Printf("Db execution err: %s", err.Error())
		return err
	}
	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("Add new comment failed")
		return err
	}
}

func ListComments(videoId string, from, to int) ([]*model.Comment, error)  {
	// inner join comments and users, specifying the video_id and time period
	stmt, err := dbConn.Prepare(`select comments.id, comments.content, users.login_name from 
		comments inner join users on comments.author_id = users.id where comments.video_id = ? and 
		comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)`)
	var comments []*model.Comment
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return comments, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(videoId, from, to)
	if err != nil {
		log.Printf("Db query err: %s", err.Error())
		return comments, err
	}

	for rows.Next() {
		var id , content, authorName string
		if err := rows.Scan(&id, &content, &authorName); err != nil {
			log.Printf("Db next err: %s", err.Error())
			return comments, err
		}
		comments = append(comments, &model.Comment{Id: id, VideoId: videoId, AuthorName: authorName, Content: content})
	}

	return comments, nil
}