package database
import (
	"fmt"
	"bytes"
	"errors"
	"strings"
	"github.com/ghawk1ns/golf/util"
	"github.com/lazyshot/go-hbase"
	"github.com/ghawk1ns/golf/model"
	"github.com/ghawk1ns/golf/logger"
)

var tableName string

// The hbase table must define these column families
var cfScores = "scores"
var cfStats = "stats"
var cfWins = "wins"

var colRounds = "rds"
var colAverage = "avg"

var rowIdPrefix = "golfer:%s"
var colWinsPrefix = "wins.over.%s"
var client *hbase.Client

func InitHBase(hbaseConfig util.HBaseConfig) {
	logger.Info.Println("Initiallizing hbase")
	client = hbase.NewClient([]string{hbaseConfig.Host}, hbaseConfig.Root)
	tableName = hbaseConfig.Table
	logger.Info.Println("hbase ready -> ", tableName, cfScores, cfStats)
}

func PutRound(round model.Round) error {
	date := round.Date
	var puts []*hbase.Put
	for _, score := range round.Scores {
		golferId := score.GolferId
		score := util.IntToStr(score.Score)
		rowId := getRowId(golferId)

		put := hbase.CreateNewPut([]byte(rowId))
		colQual := getColQualifier(date)
		put.AddStringValue(cfScores, colQual, score)

		puts = append(puts, put)
		logger.Info.Printf("%s: putting %s into scores:%s", rowId, score, colQual)
	}

	res, err := client.Puts(tableName, puts)

	if err != nil {
		logger.Error.Printf(err.Error())
		return err
	}

	if !res {
		logger.Error.Printf("No Results saved")
		return errors.New("No results saved")
	}

	logger.Info.Println("Completed put")
	return nil;
}

func GetScoresForGolfer(golferId string) (map[string]int, error) {
	scores := make(map[string]int)
	rowId := getRowId(golferId)
	logger.Info.Printf("%s: getting scores", rowId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringFamily(cfScores)
	result, err := client.Get(tableName, get)

	if err != nil {
		logger.Error.Printf(err.Error())
		return scores, err
	}

	for columnName, scoreColumn := range result.Columns {
		logger.Info.Printf("%s: columnName: %s", rowId, columnName)
		datePlayed := getColumnId(columnName)
		logger.Info.Printf("%s: datePlayed: %s", rowId, datePlayed)
		encodedScore := scoreColumn.Value

		score, err := util.StrToInt(encodedScore.String())
		if err != nil {
			logger.Error.Println(err.Error())
			return nil, err
		}

		scores[datePlayed] = score
		logger.Info.Printf("%s: %s:%s\n", rowId, datePlayed, encodedScore.String())
	}

	return scores, nil
}


func SetGolferNumRounds(golferId string, rounds int) error {
	rowId := getRowId(golferId)
	put := hbase.CreateNewPut([]byte(rowId))
	put.AddStringValue(cfStats, colRounds, util.IntToStr(rounds))
	res, err := client.Put(tableName, put)

	if err != nil {
		return err
	} else if !res {
		return errors.New("No results saved")
	}

	return nil;
}

func GetGolferNumRounds(golferId string) (int, error) {
	rowId := getRowId(golferId)

	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colRounds)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return -1, err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colRounds)]
	var result int
	if rowColResult == nil {
		// If the value isn't stored, then this golfer hasn't played any rounds yet
		result = 0
	} else {
		result, err = util.StrToInt(rowColResult.Value.String())
		if err != nil {
			logger.Error.Println(err.Error())
			return -1, err
		}

	}
	return result, nil
}


func SetGolferAverage(golferId string, average float64) error {
	rowId := getRowId(golferId)
	put := hbase.CreateNewPut([]byte(rowId))
	put.AddStringValue(cfStats, colAverage, util.FloatToStr(average))
	res, err := client.Put(tableName, put)

	if err != nil {
		return err
	} else if !res {
		return errors.New("No results saved")
	}

	return nil;
}

func GetGolferAverage(golferId string) (float64, error) {
	rowId := getRowId(golferId)

	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colAverage)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return -1, err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colAverage)]
	var result float64
	if rowColResult == nil {
		// If the value isn't stored, then this golfer hasn't played any rounds yet
		result = 0
	} else {
		result, err = util.StrToFloat(rowColResult.Value.String())
		if err != nil {
			logger.Error.Println(err.Error())
			return -1, err
		}
	}
	return result, nil
}

func SetGolferWins(golferId string, wins map[string]int) error {
	rowId := getRowId(golferId)

	var puts []*hbase.Put
	for opponentId, wins := range wins {
		put := hbase.CreateNewPut([]byte(rowId))
		put.AddStringValue(cfWins, getWinColumnId(opponentId), util.IntToStr(wins))
		puts = append(puts, put)
	}

	res, err := client.Puts(tableName, puts)

	if err != nil {
		return err
	} else if !res {
		logger.Error.Println(err.Error())
		return errors.New("No results saved")
	}

	return nil;
}

func GetGolferWins(golferId string) (map[string]int, error) {
	rowId := getRowId(golferId)

	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringFamily(cfWins)
	rowResult, err := client.Get(tableName, get)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}

	wins := make(map[string]int)
	for columnName, columnValue := range rowResult.Columns {
		logger.Info.Printf("%s: columnName: %s", rowId, columnName)
		opponentId := getGolferIdFromWinColumn(getColumnId(columnName))
		numWinsEncoded := columnValue.Value
		numWins, err := util.StrToInt(numWinsEncoded.String())
		if err != nil {
			logger.Error.Println(err.Error())
			return nil, err
		}

		wins[opponentId] = numWins
		logger.Info.Printf("%s: %s:%d\n", rowId, opponentId, numWins)
	}
	return wins, nil
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

func getWinColumnId(opponentId string) string {
	return fmt.Sprintf(colWinsPrefix, opponentId)
}

func getGolferIdFromWinColumn(winColumnId string) string {
	columnSlice := strings.SplitN(winColumnId, ".", 3)
	return columnSlice[2]
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
	logger.Info.Println("Completed put")

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

	logger.Info.Println("Completed get")

	results, err := client.Gets("test", []*hbase.Get{get})

	if err != nil {
		return err
	}

	logger.Info.Printf("hbase test success: %#v\n", results)
	return nil
}

func ResetUser(golferId string) error {

	SetGolferAverage(golferId, 0)
	SetGolferNumRounds(golferId, 0)

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