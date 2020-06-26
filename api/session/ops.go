package session

import (
	"github.com/video_server/api/utils"
	"log"
	"sync"
	"github.com/video_server/api/model"
	"github.com/video_server/api/db"
	"time"
)

var sessionMap *sync.Map

func init()  {
	sessionMap = &sync.Map{}
}

func LoadSessionFromDB() {
	m, err := db.RetrieveAllSessions()
	if err != nil {
		log.Printf("Load session from db err: %s", err.Error())
		return
	}
	m.Range(func(key, value interface{}) bool {
		session := value.(*model.Session)
		sessionMap.Store(key, session)
		return true
	})
}

func GenerateNewSession(username string) string {
	sessionId, _ := utils.NewUUID()
	ctime := time.Now().UnixNano() / 1000000
	ttl := ctime +  30 * 60 * 1000 // valid period: 30min
	session := model.Session{
		Username: username,
		TTL:      ttl,
	}

	sessionMap.Store(sessionId, &session)
	db.InsertSession(sessionId, ttl, username)
	return sessionId
}

func IsSessionExpired(sessionId string) (string, bool) {
	session, ok := sessionMap.Load(sessionId)
	now := time.Now().UnixNano() / 1000000
	if ok {
		if session.(*model.Session).TTL < now {
			// 从map和db中删除过期的session
			sessionMap.Delete(sessionId)
			db.DeleteSession(sessionId)
			return "", true
		} else {
			return session.(*model.Session).Username, false
		}
	} else {
		session, err := db.RetrieveSession(sessionId)
		if err != nil || session == nil {
			return "", true
		}
		if session.TTL < now {
			db.DeleteSession(sessionId)
			return "", true
		}

		sessionMap.Store(sessionId, session)
		return session.Username, false
	}
}




