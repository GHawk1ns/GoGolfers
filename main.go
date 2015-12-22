package main

import (
	"os"
	"io/ioutil"
	"net/http"
	"github.com/ghawk1ns/golf/util"
	"github.com/ghawk1ns/golf/routes"
	"github.com/ghawk1ns/golf/database"
	"github.com/ghawk1ns/golf/logger"
)

func main() {
	logger.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	config := util.GetConfig()
	database.InitHBase(config.HBaseConfig)
	database.InitSQL(config.SQLConfig)
	router := routes.NewRouter()
	err := http.ListenAndServe(":"+ config.Port, router)
	// We should never arrive here in most cases
	if err != nil {
		database.TryClose()
		logger.Error.Fatal(err)
	}
}