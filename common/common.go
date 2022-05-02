package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
)

type BleveFieldDef struct {
	Name string `json:"name"`
	Type string `json:"type"` //possible values [bool|text|date|numeric|geopoint]
}

///Index name date pattern should match with date value
func GetPatternIndexName(indexName string, dt string) (string, string) {
	var errMsg string
	grp := global.RegexParseDate.FindAllSubmatch([]byte(indexName), -1)
	if grp != nil {

		if len(dt) == 0 {
			errMsg = "When use index name using date pattern, you must provide date value."
			return "", errMsg
		}
		dtFormat := string(grp[0][1])

		dtVal, err := time.Parse(dtFormat, dt)

		if err != nil {
			errMsg = fmt.Sprintf("Failed while parse the date %s err %s", dt, err.Error())
			return "", errMsg
		}
		dtStr := dtVal.Format(dtFormat)
		indexName = strings.Replace(indexName, fmt.Sprintf("{%s}", dtFormat), dtStr, -1)
		//fmt.Println("formated index Name", indexName, dtVal, dtStr)

	}
	return indexName, errMsg
}

func GetIndexes(isTable bool) map[string]bool {
	var files []string
	var err error

	err = filepath.Walk(".", getIndexFolders(&files))
	if err != nil {
		logger.Error("failed while read file ", err)
		return nil
	}
	list := make(map[string]bool)
	for _, v := range files {
		list[v] = true
	}
	//fmt.Println("list of indexes", list)
	return list
}
func getIndexFolders(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
		}
		if info != nil && info.IsDir() {
			p := fmt.Sprintf("%s/%s", path, "index_meta.json")
			_, err := os.Stat(p)
			if !os.IsNotExist(err) {
				//fmt.Println("fullpath", path, info.Name())
				*files = append(*files, strings.ReplaceAll(path, `\`, "/"))
			}
		}

		return nil
	}
}
func GetAllIndexes() map[string]bool {
	var files []string
	var err error
	if err != nil {
		logger.Error("Failed get current directory", err)
	}
	err = filepath.Walk(".", getIndexFolders(&files))

	if err != nil {
		logger.Error("Failed while read file from index folder", err)
		return nil
	}
	list := make(map[string]bool)
	for _, v := range files {
		list[v] = true
	}
	return list
}
