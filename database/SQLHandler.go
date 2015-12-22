package database

import (
	"github.com/ghawk1ns/golf/model"
	"github.com/ghawk1ns/golf/util"
	"database/sql"
	"github.com/ghawk1ns/golf/logger"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitSQL(dbConfig util.SQLConfig) {
	logger.Info.Println("initializing database")
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dataSourceName := dbConfig.User + ":" + dbConfig.Password + "@tcp("+ dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.Database
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	} else {
		logger.Info.Println("sql connection opened")
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	} else {
		logger.Info.Println("sql connection established")
	}
}

func TryClose() {
	if db != nil {
		db.Close()
	}
}

func GetGolfers() ([]model.Golfer, error) {
	var golfers []model.Golfer
	rows, err := db.Query("SELECT * FROM golfers")
	if err != nil {
		logger.Error.Println("fuck not good: ", err)
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var golferId string
		var name string
		// imageUrl is nullable, so it must be set to []byte
		var imageUrl []byte
		err = rows.Scan(&golferId, &name, &imageUrl)
		if err != nil {
			logger.Error.Println("bad row bro: ", err)
		} else {
			golfers = append(golfers, model.Golfer{golferId, name, string(imageUrl)})
		}
	}
	return golfers, nil
}

// For Testing
func GetGolferById(golferId string) (model.Golfer, error) {
	return getGolfer("golferId", golferId)
}

// For Testing
func getGolferByName(name string) (model.Golfer, error) {
	return getGolfer("name", name)
}

// For Testing
func getGolfer(field string, value string) (model.Golfer, error) {

	logger.Info.Printf("GetGolfer: field -> %s, value -> %s\n", field, value)

	// Execute the query
	rows, err := db.Query("SELECT * FROM golfers WHERE " + field + "=" + value)
	if err != nil {
		logger.Error.Println("fuck not good: ", err)
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var golferId string
		var name string
		var imageUrl []byte
		err = rows.Scan(&golferId, &name, &imageUrl)
		if err != nil {
			logger.Error.Println("bad row bro: ", err)
		} else {
			return model.Golfer{golferId, name, string(imageUrl)}, err
		}
	}
	return model.Golfer{}, err;
}