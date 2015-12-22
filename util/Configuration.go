package util

import (
	"os"
	"encoding/json"
	"github.com/ghawk1ns/golf/logger"
)

var configPointer *Configuration

// Config corresponding to config.json go in here
type Configuration struct {
	Secret      string 		`json:"secret"`
	Port        string 		`json:"port"`
	SQLConfig   SQLConfig 	`json:"sql"`
	HBaseConfig HBaseConfig `json:"hbase"`
}

type SQLConfig struct {
	User 	 string `json:"user"`
	Password string `json:"password"`
	Host 	 string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type HBaseConfig struct {
	Host 	 string `json:"host"`
	Root 	 string `json:"root"`
	Table 	 string `json:"table"`
}

// http://stackoverflow.com/a/16466189
func GetConfig() Configuration {
	if (configPointer == nil) {
		logger.Info.Println("Creating new config")
		file, _ := os.Open("conf.json")
		decoder := json.NewDecoder(file)
		configPointer = &Configuration{}
		err := decoder.Decode(configPointer)
		if err != nil {
			logger.Error.Println("error:", err)
		}
	}
	return *configPointer
}

func GetSecret() string {
	return GetConfig().Secret
}