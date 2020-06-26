package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"video_server/config"
	"log"
)

var (
	EndPoint string
	AccessKey string
	SecretKey string
)

func init()  {
	EndPoint = config.GetOSSAddr()
	AccessKey = "LTAI4G8FyHb4n2trrZVwfwBU"
	SecretKey = "O2J6LnFPZfKP1QsV3P5JUqQW10P8Pe"
}

func UploadToOSS(fileName, path, bucketName string) bool {
	// 获取client
	client, err := oss.New(EndPoint, AccessKey, SecretKey)
	if err != nil {
		log.Printf("OSS service err: %s", err.Error())
		return false
	}

	// 获取bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Printf("Get bucket err: %s", err.Error())
		return false
	}

	// 向指定bucket并发上传文件
	err = bucket.UploadFile(fileName, path, 500 * 1024, oss.Routines(3))
	if err != nil {
		log.Printf("File upload err: %s", err.Error())
		return false
	}

	return true
}


