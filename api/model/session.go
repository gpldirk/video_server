package model

type Session struct {
	Username string
	TTL int64
}

type UserSession struct {
	Username string `json:"username"`
	SessionId string `json:"session_id"`
}
