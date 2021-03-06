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
var colRoundsFmt = "rds.%s"

var colAverage = "avg"
var colAverageFmt = "avg.%s"


var rowIdPrefix = "golfer:%s"
// courseId.date.timestamp
var colScoreQual = "%s.%s.%s"
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
		colQual := getColQualifier(round.CourseId, date)
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

func GetScoresForGolfer(golferId string) (map[string]map[string]int, error) {
	allScores := make(map[string]map[string]int)
	rowId := getRowId(golferId)
	logger.Info.Printf("%s: getting scores", rowId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringFamily(cfScores)
	result, err := client.Get(tableName, get)

	if err != nil {
		logger.Error.Printf(err.Error())
		return allScores, err
	}

	for columnName, scoreColumn := range result.Columns {
		logger.Info.Printf("%s: columnName: %s", rowId, columnName)
		courseId, datePlayed := GetScoreInfoFromColId(getColumnId(columnName))
		logger.Info.Printf("%s: datePlayed: %s", rowId, datePlayed)
		encodedScore := scoreColumn.Value

		score, err := util.StrToInt(encodedScore.String())
		if err != nil {
			logger.Error.Println(err.Error())
			return nil, err
		}

		courseScores, ok := allScores[courseId]
		if !ok {
			logger.Info.Println("Creating new map")
			courseScores = make(map[string]int)
			allScores[courseId] = courseScores
		} else {
			logger.Info.Println("map exists")
		}
		courseScores[datePlayed] = score
		logger.Info.Printf("%s: %s-%s:%d\n", rowId, courseId, datePlayed, score)
	}

	return allScores, nil
}


func SetGolferNumRounds(golferId string, courseId string, rounds int) error {
	rowId := getRowId(golferId)
	put := hbase.CreateNewPut([]byte(rowId))
	put.AddStringValue(cfStats, getRndColId(courseId), util.IntToStr(rounds))
	res, err := client.Put(tableName, put)

	if err != nil {
		return err
	} else if !res {
		return errors.New("No results saved")
	}

	return nil;
}

func IncGolferTotalRounds(golferId string) error {
	rowId := getRowId(golferId)
	colId := getRndColId("total")
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colId)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colId)]
	var result int
	if rowColResult == nil {
		// If the value isn't stored, then this golfer hasn't played any rounds yet
		result = 0
	} else {
		result, err = util.StrToInt(rowColResult.Value.String())
		if err != nil {
			logger.Error.Println(err.Error())
			return err
		}

	}

	return SetGolferNumRounds(golferId, "total", result + 1)
}

func GetGolferNumRounds(golferId string, courseId string) (int, error) {
	rowId := getRowId(golferId)
	colId := getRndColId(courseId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colId)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return -1, err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colId)]
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


func SetGolferAverage(golferId string, courseId string, average float64) error {
	rowId := getRowId(golferId)
	colId := getAvgColId(courseId)
	put := hbase.CreateNewPut([]byte(rowId))
	put.AddStringValue(cfStats, colId, util.FloatToStr(average))
	res, err := client.Put(tableName, put)

	if err != nil {
		return err
	} else if !res {
		return errors.New("No results saved")
	}

	return nil;
}

func GetGolferAverage(golferId string, courseId string) (float64, error) {
	rowId := getRowId(golferId)
	colId := getAvgColId(courseId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringColumn(cfStats, colId)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return -1, err
	}
	// this result could be empty
	rowColResult := rowResult.Columns[getFullColumnName(cfStats, colId)]
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


func GetAllRoundsForGolfer(golferId string) (map[string]int, error) {
	result, err := getStatsForGolfer(golferId, colRounds)
	if err != nil {
		return nil, err
	}
	convertedResult := make(map[string]int)
	for key,val := range result {
		newVal, err := util.StrToInt(val)
		if err != nil {
			return nil, err
		} else {
			convertedResult[key] = newVal
		}
	}
	return convertedResult, nil
}

func GetAllAveragesForGolfer(golferId string) (map[string]float64, error) {
	result, err := getStatsForGolfer(golferId, colAverage)
	if err != nil {
		return nil, err
	}
	convertedResult := make(map[string]float64)
	for key,val := range result {
		newVal, err := util.StrToFloat(val)
		if err != nil {
			return nil, err
		} else {
			convertedResult[key] = newVal
		}
	}
	return convertedResult, nil
}

func getStatsForGolfer(golferId string, preferredStatType string) (map[string]string, error) {
	rowId := getRowId(golferId)
	get := hbase.CreateNewGet([]byte(rowId))
	get.AddStringFamily(cfStats)
	rowResult, err := client.Get(tableName, get)

	if err != nil {
		return nil, err
	}
	stats := make(map[string]string)
	for _, resultRowCol := range rowResult.Columns {
		logger.Info.Printf("%s: columnName: %s-%s\n", rowId, resultRowCol.Family.String(), resultRowCol.Qualifier.String())
		statType, courseId := GetStatInfoFromColId(resultRowCol.Qualifier.String())
		if statType == preferredStatType {
			encodedStat := resultRowCol.Value
			stats[courseId] = encodedStat.String()
			logger.Info.Printf("%s: %s-%s:%s\n", rowId, courseId, statType, encodedStat.String())
		}
	}
	return stats, nil
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

// colScoreQual -> courseId.date.timestamp
func GetScoreInfoFromColId(columnName string) (string, string) {
	logger.Info.Println("ScoreInfo: %s", columnName)
	columnSlice := strings.SplitN(columnName, ".", 3)
	logger.Info.Printf("ScoreInfo: %s - %s\n", columnSlice[0], columnSlice[1])
	return columnSlice[0],columnSlice[1]
}

// stat col qual -> type.courseId
// type -> rds or avg
func GetStatInfoFromColId(columnName string) (string, string) {
	logger.Info.Printf("StatInfoColdId: %s\n", columnName)
	columnSlice := strings.SplitN(columnName, ".", 2)
	logger.Info.Printf("StatInfo from Col: %s - %s\n", columnSlice[0], columnSlice[1])
	return columnSlice[0],columnSlice[1]
}

func getFullColumnName(colFamily string, colQualifier string) string {
	return fmt.Sprintf("%s:%s", colFamily, colQualifier)
}

func getAvgColId(courseId string) string {
	return fmt.Sprintf(colAverageFmt, courseId)
}

func getRndColId(courseId string) string {
	return fmt.Sprintf(colRoundsFmt, courseId)
}

func getRowId(golferId string) string {
	return fmt.Sprintf(rowIdPrefix, golferId)
}

func getColQualifier(courseId, date string) string {
	return fmt.Sprintf(colScoreQual, courseId, date, util.MakeTimestamp())
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