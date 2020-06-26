package model

type NewVideo struct {
	AuthorId int `json:"author_id"`
	Name string `json:"name"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type VideoInfo struct {
	Id string `json:"id"`
	AuthorId int `json:"author_id"`
	Name string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}



