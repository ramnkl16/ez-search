package ezsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
	"go.uber.org/zap/zapcore"
)

// func CreateNewIndex(indexName string, dateFields []string, isTable bool) rest_errors.RestErr {
// 	index, err := GetIndex(indexName)

// 	if err != nil || index == nil {
// 		logger.Debug("Creating  new index ... ", zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
// 		// create a mapping
// 		indexMapping, err := buildDefultIndexMapping(indexName, dateFields)
// 		if err != nil {
// 			log.Fatal(err)
// 			return rest_errors.NewBadRequestError(err.Error())
// 		}
// 		var indexPath string
// 		if !isTable {
// 			indexPath = fmt.Sprintf("%s%c%s", Conf.IndexBasePath, '/', indexName)
// 		} else {
// 			indexPath = fmt.Sprintf("%s%c%s", Conf.IndexTablesPath, '/', indexName)
// 		}
// 		logger.Debug("index path", zapcore.Field{String: indexPath, Key: "p1", Type: zapcore.StringType})
// 		index, err := bleve.New(indexPath, indexMapping)
// 		if err != nil {
// 			logger.Error("Failed product index field mapping", err)
// 			return rest_errors.NewBadRequestError(err.Error())
// 		}
// 		cache_utils.AddOrUpdateCache(indexName, index)
// 	}
// 	if index == nil {
// 		logger.Info("index is not created", zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
// 	}
// 	return nil

// }

///

func GetBleveIndexSchema(indexName string) ([]common.BleveFieldDef, error) {
	key := fmt.Sprintf("%s.schema", indexName)
	logger.Debug("GetBleveIndexSchema|schema key", zapcore.Field{Type: zapcore.StringType, Key: key})
	schemaByte, err := coredb.GetValue(coredb.Defaultbucket, key)
	if err != nil {
		errStr := fmt.Sprintf(`GetBleveIndexSchema|%s schema is not found in core db. Please try after create schema first. \n%s\n`, indexName, err.Error())
		logger.Error(errStr, errors.New(""))
		return nil, errors.New(errStr)
	}
	if len(schemaByte) == 0 {
		errStr := fmt.Sprintf(`GetBleveIndexSchema|%s schema is not found in core db. Please try after create schema first.`, indexName)
		logger.Error(errStr, errors.New(""))
		return nil, errors.New(errStr)
	}

	var schemaDefs []common.BleveFieldDef
	json.Unmarshal(schemaByte, &schemaDefs)
	return schemaDefs, nil
}
func DeleteIndexDocs(indexNames []string) error {
	for _, indexName := range indexNames {
		index, err := GetIndex(indexName)
		if err != nil {
			return err
		}
		cache_utils.Cache.Remove(indexName)
		key := fmt.Sprintf("%s.schema", indexName)
		coredb.Delete(coredb.Defaultbucket, key)

		index.Close()
		fullpath := filepath.Join(global.WorkingDir, indexName)
		err1 := os.RemoveAll(fullpath)
		if err1 != nil {
			logger.Error("DeleteIndexDocs|Failed to delete index folder", err1, zapcore.Field{String: fullpath, Key: "p1", Type: zapcore.StringType})
			continue
		}
		logger.Info("DeleteIndexDocs|Index deleted", zapcore.Field{String: indexName, Key: "p1", Type: zapcore.StringType})
	}
	return nil
}
