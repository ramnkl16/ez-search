package ezcsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/global"
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
		fmt.Println(rowIdex, record)
		if rowIdex == 0 {
			continue

		}
		indexDoc := make(map[string]interface{})
		for idx, c := range h {
			//fmt.Println("col", idx, c, record[idx])
			if len(record[idx]) > 0 {
				indexDoc[c] = record[idx]
			}
		}
		var id string
		if uniqueColIndex == -1 {
			id = uid_utils.GetUid("cs", false)
		} else {
			id = record[uniqueColIndex]
		}
		indexes[id] = indexDoc
		//fmt.Println("index", id, indexDoc)
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
	// if _, err := r.Read(); err != nil {
	// 	return [][]string{}, err
	// }

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

///column name must be comma separated
///it would generate schema for bleve index search
func GenerateIndexSchema(columnNames string, enabledShortName bool) []common.BleveFieldDef {
	schemaList := make([]common.BleveFieldDef, 0)
	names := make(map[string]bool, 0)
	var colCount int = 1
	for _, v := range strings.Split(columnNames, ",") {
		words := getDisplayName(v)
		var name string
		//get the first letter of each word
		if enabledShortName {
			for _, w := range words {
				name = name + w[0:1]
			}
			name = strings.ToLower(name)
		} else {
			name = strings.ToLower(words[0])
			if len(words) > 1 {
				name = name + strings.Join(words[1:], "")
			}
		}
		exist, _ := names[name]
		if exist {
			name = name + string(rune((colCount)))
			colCount = colCount + 1
		}

		names[name] = true
		bd := common.BleveFieldDef{Name: name, DisplayName: strings.Join(words, " "), Type: "text"}
		schemaList = append(schemaList, bd)
	}
	return schemaList
}

func getDisplayName(name string) []string {
	grp := global.RegexParseHasCapitalLetter.FindAllStringIndex(name, -1)
	list := make([]string, 0)

	curIdx := 0
	//maxLen := len(name)
	//fmt.Println("name", name)
	for grpIndex, r := range grp {
		//fmt.Println("r", r)
		if grpIndex != 0 {
			word := strings.Title(name[curIdx : r[1]-1])
			list = append(list, strings.Split(word, "_")...)
		}
		curIdx = r[0]

	}
	r := grp[len(grp)-1]
	word := strings.Title(name[r[0]:])
	list = append(list, strings.Split(word, "_")...)
	//print(list)

	//fmt.Println("query parser", kw)
	return list

}
