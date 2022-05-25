package ezcsv

import (
	"encoding/csv"
	"os"

	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
	"go.uber.org/zap/zapcore"
)

func GetJsonFromCsv(fileName string, uniqueColIndex int) (map[string]interface{}, rest_errors.RestErr) {
	records, err := readData(fileName)

	if err != nil {
		logger.Error("Failed while open file", err, zapcore.Field{Key: "fileName", Type: zapcore.StringType, String: fileName})
		return nil, rest_errors.NewInternalServerError("Failed while open File|getJsonFromCsv", err)
	}
	h := records[0] //assuming first row will be header
	indexes := make(map[string]interface{})
	for rowIdex, record := range records {
		if rowIdex == 0 {
			continue
		}
		indexDoc := make(map[string]interface{})
		for idx, c := range h {
			indexDoc[c] = record[idx]
		}
		var id string
		if uniqueColIndex == -1 {
			id = uid_utils.GetUid("cs", false)
		} else {
			id = record[uniqueColIndex]
		}
		indexes[id] = indexDoc
	}
	return indexes, nil
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
