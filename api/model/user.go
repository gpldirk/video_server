package model

type UserCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id int
	Username string
	Password string
}

type UserInfo struct {
	Id int `json:"user_id"`
}

type SignUp struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}

type SignIn struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}


