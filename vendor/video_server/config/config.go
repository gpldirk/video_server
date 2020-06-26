package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	LBAddr string `json:"lb_addr"`
	OSSAddr string `json:"OSS_addr"`
}

var configuration *Configuration

func init() {
	file, _ := os.Open("./conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = &Configuration{}
	err := decoder.Decode(configuration)
	if err != nil {
		panic(err)
	}
}

func GetLBAddr() string {
	return configuration.LBAddr
}

func GetOSSAddr() string {
	return configuration.OSSAddr
}



