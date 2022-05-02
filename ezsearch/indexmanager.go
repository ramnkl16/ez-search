package ezsearch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/logger"
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

func GetBleveTableschema(indexName string) ([]common.BleveFieldDef, error) {
	key := fmt.Sprintf("%s.schema", indexName)
	logger.Debug("schema key", zapcore.Field{Type: zapcore.StringType, Key: key})
	schemaByte, err := coredb.GetKey(key)
	if err != nil {
		errStr := fmt.Sprintf(`%s schema is not found in core db. Please try after create schema first. \n%s\n`, indexName, err.Error())
		logger.Error(errStr, errors.New(""))
		return nil, errors.New(errStr)
	}
	if len(schemaByte) == 0 {
		errStr := fmt.Sprintf(`%s schema is not found in core db. Please try after create schema first.`, indexName)
		logger.Error(errStr, errors.New(""))
		return nil, errors.New(errStr)
	}

	var schemaDefs []common.BleveFieldDef
	json.Unmarshal(schemaByte, &schemaDefs)
	return schemaDefs, nil
}

//}
