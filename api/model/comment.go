package model

// request
type Comment struct {
	Id string `json:"comment_id"`
	VideoId string `json:"video_id"`
	AuthorName string `json:"author_name"`
	Content string `json:"content"`
}

// response
type NewComment struct {
	AuthorId int `json:"author_id"`
	Content string `json:"content"`
}

// model
type Comments struct {
	Comments []*Comment `json:"comments"`
}
