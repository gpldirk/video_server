package db

import (
	"log"
	"strconv"
	"sync"
	"github.com/video_server/api/model"
)

func InsertSession(sessionId string, ttl int64, username string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmt, err := dbConn.Prepare("insert into sessions (session_id, TTL, login_name) values (?, ?, ?)")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(sessionId, ttlstr, username)
	if err != nil {
		log.Printf("Db execution err: %s", err.Error())
		return err
	}
	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("Insert session failed")
		return err
	}
}

func RetrieveSession(sessionId string) (*model.Session, error) {
	stmt, err := dbConn.Prepare("select TTL, login_name from sesisons where session_id=?")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var ttlstr, username string
	_ = stmt.QueryRow(sessionId).Scan(&ttlstr, &username)

	var session model.Session
	if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err == nil {
		session.Username = username
		session.TTL = ttl
	} else {
		log.Printf("TTL parse err: %s", err.Error())
		return nil, err
	}

	return &session, nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmt, err := dbConn.Prepare("select * from sessions")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return m, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Printf("DB query err: %s", err.Error())
		return m, err
	}

	var id, ttlstr, username string
	for rows.Next() {
		if err := rows.Scan(&id, &ttlstr, &username); err != nil {
			log.Printf("retrieve session err: %s", err.Error())
			return m, err
		}
		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err != nil {
			log.Printf("TTL parse err: %s", err.Error())
			return m, err
		} else {
			m.Store(id, &model.Session{
				Username: username,
				TTL:      ttl,
			})
		}
	}

	return m, nil
}

func DeleteSession(sessionId string) error {
	stmt, err := dbConn.Prepare("delete from sessions where session_id=?")
	if err != nil {
		log.Printf("Db connection err: %s", err.Error())
		return err
	}
	defer stmt.Close()

	rf, err := stmt.Exec(sessionId)
	if err != nil {
		log.Printf("Db execution err: %s", err.Error())
		return err
	}
	if rows, err := rf.RowsAffected(); err == nil && rows > 0 {
		return nil
	} else {
		log.Printf("Delete session failed")
		return err
	}
}


