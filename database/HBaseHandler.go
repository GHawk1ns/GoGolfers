package database
import (
	"fmt"
	"bytes"
	"errors"
	"strings"
	"github.com/ghawk1ns/golf/util"
	"github.com/lazyshot/go-hbase"
	"github.com/ghawk1ns/golf/model"
	"github.com/ghawk1ns/golf/blah"
)

var tableName string

// The hbase table must define these column families
var cfScores = "scores"
var cfStats = "stats"

var colRounds = "rds"
var colAverage = "avg"

var rowIdPrefix = "golfer:%s"
var colWinsPrefix = "wins.golfer.%s"
var client *hbase.Client

func InitHBase(hbaseConfig util.HBaseConfig) {
	blah.Info.Println("Initiallizing hbase")
	client = hbase.NewClient([]string{hbaseConfig.Host}, hbaseConfig.Root)
	tableName = hbaseConfig.Table
	blah.Info.Println("hbase ready -> ", tableName, cfScores, cfStats)
}

func PutRound(round model.Round) error {
	date := round.Date
	var puts []*hbase.Put
	for _, score := range round.Scores {
		golferId := score.GolferId
		score := score.Score
		rowId := getRowId(golferId)

		put := hbase.CreateNewPut([]byte(rowId))
		put.AddStringValue(cfScores, getColQualifier(date), score)

		puts = append(puts, put)
		blah.Info.Printf("%s: putting %s into scores:%s", rowId, score, cfScores)
	}

	res, err := client.Puts(tableName, puts)

	if err != nil {
		blah.Error.Printf(err.Error())
		return err
	}

	if !res {
		blah.Error.Printf("No Results saved")
		return errors.New("No results saved")
	}

	blah.Info.Println("Completed put")
	return nil;
}

func GetScoresForGolfer(golferId string) (map[string]string, error) {
	scores := make(map[string]string)
	rowId := getRowId(golferId)
	blah.Info.Printf("%s: getting scores", rowId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringFamily(cfScores)
	result, err := client.Get(tableName, get)

	if err != nil {
		blah.Error.Printf(err.Error())
		return scores, err
	}

	for columnName, scoreColumn := range result.Columns {
		blah.Info.Printf("%s: columnName: %s", rowId, columnName)
		datePlayed := getColumnId(columnName)
		blah.Info.Printf("%s: datePlayed: %s", rowId, datePlayed)
		encodedScore := scoreColumn.Value
		scores[datePlayed] = encodedScore.String()
		blah.Info.Printf("%s: %s:%s\n", rowId, datePlayed, encodedScore.String())
	}

	return scores, nil
}


func SetGolferNumRounds(golferId string, rounds string) error {
	rowId := getRowId(golferId)
	put := hbase.CreateNewPut([]byte(rowId))
	put.AddStringValue(cfStats, colRounds, rounds)
	res, err := client.Put(tableName, put)

	if err != nil {
		return err
	} else if !res {
		return errors.New("No results saved")
	}

	return nil;
}

func GetGolferNumRounds(golferId string) (string, error) {
	rowId := getRowId(golferId)

	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colRounds)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return "", err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colRounds)]
	var result string
	if rowColResult == nil {
		// If the value isn't stored, then this golfer hasn't played any rounds yet
		result = "0"
	} else {
		result = rowColResult.Value.String()
	}
	return result, nil
}


func SetGolferAverage(golferId string, average string) error {
	rowId := getRowId(golferId)
	put := hbase.CreateNewPut([]byte(rowId))
	put.AddStringValue(cfStats, colAverage, average)
	res, err := client.Put(tableName, put)

	if err != nil {
		return err
	} else if !res {
		return errors.New("No results saved")
	}

	return nil;
}

func GetGolferAverage(golferId string) (string, error) {
	rowId := getRowId(golferId)

	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colAverage)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return "", err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colAverage)]
	var result string
	if rowColResult == nil {
		// If the value isn't stored, then this golfer hasn't played any rounds yet
		result = "0"
	} else {
		result = rowColResult.Value.String()
	}
	return result, nil
}

// cf:columnid -> [cf, columnId] -> columnId
func getColumnId(columnName string) string {
	columnSlice := strings.SplitN(columnName, ":", 2)
	return columnSlice[1]
}

func getFullColumnName(colFamily string, colQualifier string) string {
	return fmt.Sprintf("%s:%s", colFamily, colQualifier)
}

func getRowId(golferId string) string {
	return fmt.Sprintf(rowIdPrefix, golferId)
}

func getColQualifier(date string) string {
	return fmt.Sprintf("%s.%s", date, util.MakeTimestamp())
}

func getWinColumn(opponentId string) string {
	return fmt.Sprintf(colWinsPrefix, opponentId)
}

func Test(test_val string) error {

	put := hbase.CreateNewPut([]byte("test1"))
	put.AddStringValue("info", "test_qual", test_val)
	res, err := client.Put("test", put)

	if err != nil {
		return err
	}

	if !res {
		return errors.New("No Put Result")
	}
	blah.Info.Println("Completed put")

	get := hbase.CreateNewGet([]byte("test1"))
	result, err := client.Get("test", get)

	if err != nil {
		panic(err)
	}

	if !bytes.Equal(result.Row, []byte("test1")) {
		return errors.New("No Row")
	}

	if !bytes.Equal(result.Columns["info:test_qual"].Value, []byte(test_val)) {
		return errors.New("Value doesn't match")
	}

	blah.Info.Println("Completed get")

	results, err := client.Gets("test", []*hbase.Get{get})

	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", results)
	return nil
}

func ResetUser(golferId string) error {

	SetGolferAverage(golferId, "0")
	SetGolferNumRounds(golferId, "0")

	rowId := getRowId(golferId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringFamily(cfScores)
	result, _ := client.Get(tableName, get)

	var dels []*hbase.Delete
	for columnName, _ := range result.Columns {
		del := hbase.CreateNewDelete([]byte(rowId))
		del.AddStringColumn(cfScores, getColumnId(columnName))
		dels = append(dels, del)
	}
	client.Deletes(tableName, dels)
	return nil
}