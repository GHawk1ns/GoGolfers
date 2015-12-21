package main

import (
	"os"
	"io/ioutil"
	"net/http"
	"github.com/ghawk1ns/golf/util"
	"github.com/ghawk1ns/golf/routes"
	"github.com/ghawk1ns/golf/database"
	"github.com/ghawk1ns/golf/blah"
)

func main() {

	/**
		A basic restful API

		TODO: Handle errors better
		TODO: Clean up unused endpoints
		TODO: Add cooler, newer endpoints
		TODO: Auth
	 */

	blah.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	config := util.GetConfig()
	database.InitHBase(config.HBaseConfig)
	database.InitSQL(config.SQLConfig)
	router := routes.NewRouter()
	err := http.ListenAndServe(":"+ config.Port, router)
	// We should never arrive here in most cases
	if err != nil {
		database.TryClose()
		blah.Error.Fatal(err)
	}
}