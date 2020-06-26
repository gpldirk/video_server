package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// init(connect db, truncate tables) -> run tests -> truncate tables
func TestMain(m *testing.M)  {
	truncateTables()
	m.Run()
	truncateTables()
}

func truncateTables()  {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate videos")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

// test users
func TestUserWorkFlow(t *testing.T) {
	truncateTables()
	t.Run("Add user", testAddUserCredential)
	t.Run("Get user", testGetUserCredential)
	t.Run("Delete user", testDeleteUserCredential)
	t.Run("ReGet user", testReGetUserCredential)
	truncateTables()
}

func testAddUserCredential(t *testing.T) {
	err := AddUserCredential("GepeiLu", "lu13120515928")
	if err != nil {
		t.Errorf("Add user err: %s", err.Error())
	}
}

func testGetUserCredential(t *testing.T) {
	password, err := GetUserCredential("GepeiLu")
	if err != nil {
		t.Errorf("Get user err: %s", err.Error())
	}
	if password != "lu13120515928" {
		t.Errorf("Get user password err")
	}
}

func testDeleteUserCredential(t *testing.T) {
	err := DeleteUserCredential("GepeiLu", "lu13120515928")
	if err != nil {
		t.Errorf("Delete user err: %s", err.Error())
	}
}

func testReGetUserCredential(t *testing.T) {
	password, err := GetUserCredential("GepeiLu")
	if err != nil {
		t.Errorf("ReGet user err: %s", err.Error())
	}
	if password != "" {
		t.Errorf("Delete user failed")
	}
}

// test videos
var id string

func TestVideoWorkFlow(t *testing.T)  {
	truncateTables()
	t.Run("Add new user", testAddUserCredential)
	t.Run("Add new video", testAddNewVideo)
	t.Run("Get video info", testGetVideoInfo)
	t.Run("Delete video", testDeleteVideo)
	t.Run("ReGet video info", testReGetVideoInfo)
	truncateTables()
}

func testAddNewVideo(t *testing.T) {
	video, err := AddNewVideo(1, "Happy")
	if err != nil {
		t.Errorf("Add new video err: %s", err.Error())
	}
	id = video.Id
}

func testGetVideoInfo(t *testing.T) {
	curVideo, err := GetVideoInfo(id)
	if err != nil {
		t.Errorf("Get video info err: %s", err.Error())
	}
	if curVideo == nil {
		t.Errorf("Add new video failed")
	}
}

func testDeleteVideo(t *testing.T) {
	err := DeleteVideoInfo(id)
	if err != nil {
		t.Errorf("Delete video err: %s", err.Error())
	}
}

func testReGetVideoInfo(t *testing.T) {
	curVideo, err := GetVideoInfo(id)
	if err != nil && err != sql.ErrNoRows{
		t.Errorf("ReGet video info err: %s", err.Error())
	}
	if err != sql.ErrNoRows || curVideo != nil {
		t.Errorf("Delete video failed")
	}
}

// test comments
func TestCommentsWorkFlow(t *testing.T) {
	truncateTables()
	t.Run("Add new user", testAddUserCredential)
	t.Run("Add new comment", testAddNewComment)
	t.Run("List all comments", testListComments)
	truncateTables()
}

func testAddNewComment(t *testing.T) {
	videoId := "12345"
	authorId := 1
	content := "I love this video"
	err := AddNewComment(videoId, authorId, content)
	if err != nil {
		t.Errorf("Add new comment err: %s", err.Error())
	}
}

func testListComments(t *testing.T) {
	videoId := "12345"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano() / 1000000000, 10))
	comments, err := ListComments(videoId, from, to)
	if err != nil {
		t.Errorf("List comments err:%s", err.Error())
	}

	for i, comment := range comments {
		fmt.Printf("Comment %d: %v\n", i, comment)
	}
}

