package taskrunner

import (
	"errors"
	"log"
	"sync"
	"github.com/video_server/scheduler/db"
	"github.com/video_server/scheduler/oss"
	"os"
)

// api -> video_id -> mysql -> video_del_rec
// scheduler -> dispatcher -> video_del_rec -> data chan
// executor -> data chan -> video_id -> delete video

func deleteVideo(videoId string) error {
	err := os.Remove("./videos/" + videoId)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Delete video err: %s", err.Error())
		return err
	}

	OSSfn := "videos/" + videoId
	bucketName := "video-server-gepeilu"
	if success := oss.DeleteObject(OSSfn, bucketName); !success {
		return errors.New("Delete video error")
	}

	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := db.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Db read err: %s", err.Error())
		return err
	}

	if len(res) == 0 {
		return errors.New("All tasks finished")
	}

	for _, videoId := range res {
		dc <- videoId
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
	forloop:
		for {
			select {
			case videoId := <- dc:
				// d: 1 2 3
				// e: 1 2 3 -> 1
				// d: 2 3
				// e: 2 3 ...
				// 重复读取和删除的问题， waitGroup解决
				go func(id interface{}) {
					// 删除db记录
					if err := db.DelVideoDeletionRecord(id.(string)); err != nil {
						errMap.Store(id, err)
						return
					}
					// 删除video文件
					if err := deleteVideo(id.(string)); err != nil {
						errMap.Store(id, err)
						return
					}
				}(videoId) // 闭包实现状态保存和异步删除
			default:
				break forloop
			}
		}

	errMap.Range(func(key, value interface{}) bool {
		err := value.(error)
		if err != nil {
			return false
		}
		return true
	})
	
	return err
}


