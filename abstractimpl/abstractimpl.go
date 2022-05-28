package abstractimpl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/ezsearch"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/uid_utils"

	"go.uber.org/zap/zapcore"
)

// Insert the EventQueue to the database.
func CreateOrUpdate[T any](model T, tableName string, id string) rest_errors.RestErr {
	t, err := GetTable(tableName)
	if err != nil {
		return err
	}
	//fmt.Println("abstract|CreateOrUpdate", tableName, model)
	if len(id) == 0 {
		uid_utils.GetUid("rt", true)
	}
	err1 := t.Index(id, model)
	if err != nil {
		return rest_errors.NewInternalServerError("abstractimpl|Failed while update", err1)
	}
	return nil
}

// Insert the EventQueue to the database.
func Put[T any](kv map[string]string, tableName string, id string) rest_errors.RestErr {
	t, err := GetTable(tableName)
	if err != nil {
		return err
	}
	oldM, err := Get[T](fmt.Sprintf("abstractimpl|select * from %s where id=%s", tableName, id))
	if err != nil {
		msg := fmt.Sprintf("abstractimpl|Failed while partial update %s", tableName)
		return rest_errors.NewInternalServerError(msg, err)
	}
	fields := reflect.ValueOf(&oldM).Elem()
	shouldReturn, returnValue := updateValueUsingreflection(kv, fields)
	if shouldReturn {
		return returnValue
	}

	err1 := t.Index(id, oldM)
	if err != nil {
		return rest_errors.NewInternalServerError("abstractimpl|Failed while insert record", err1)
	}
	return nil
}

func updateValueUsingreflection(kv map[string]string, fields reflect.Value) (bool, rest_errors.RestErr) {
	for k, v := range kv {
		f := fields.FieldByName(k)
		if f.CanSet() {
			switch f.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				x, err := strconv.Atoi(v)
				if err != nil {
					return true, rest_errors.NewBadRequestError(err.Error())
				}
				f.SetInt(int64(x))
			case reflect.Float32, reflect.Float64:
				x, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return true, rest_errors.NewBadRequestError(err.Error())
				}
				f.SetFloat(x)
			case reflect.Bool:
				x, err := strconv.ParseBool(v)
				if err != nil {
					return true, rest_errors.NewBadRequestError(err.Error())
				}
				f.SetBool(x)
			case reflect.String:
				f.SetString(v)
			default:
				return true, rest_errors.NewBadRequestError(fmt.Sprintf("Struct field type is not [int,float and string] %s %s", k, v))
			}
		}

	}
	return false, nil
}

// Delete a recored ( EventQueue) from the database.
func Delete(tableName string, id string) rest_errors.RestErr {
	t, err := GetTable(tableName)
	if err != nil {
		return err
	}
	err1 := t.Delete(id)
	if err1 != nil {
		return rest_errors.NewInternalServerError("abstractimpl|Failed while delete eventqueue", err1)
	}
	return nil
}

func BatchCreateOrUpdate(indexName string, data map[string]interface{}) rest_errors.RestErr {
	i, err := GetTable(indexName)
	if err != nil {
		logger.Warn("Schema is missing. build schema using {{baseUrl}}/api/createschema?indexName=macindex/new/customer api")
		return rest_errors.NewInternalServerError("Failed missing index schema", err)
		//i, err = BuilddynamicSchema(indexName)
	}
	batchCount := 0
	batch := i.NewBatch()
	count := 0
	for k, v := range data {
		batch.Index(k, v)
		batchCount = batchCount + 1
		count = count + 1
		if batchCount >= global.MaxIndexbatchSize {
			logger.Info("batch executed", zapcore.Field{Key: "executed", Type: zapcore.Int32Type, Integer: int64(count)})
			err := i.Batch(batch)
			if err != nil {
				logger.Error("Failed in the batch index", err)
				return rest_errors.NewInternalServerError("Failed in the batch index", err)
			}
			batchCount = 0
		}
	}
	if batchCount > 0 {
		err := i.Batch(batch)
		if err != nil {
			logger.Error("Failed while last batch ", err)
			return rest_errors.NewInternalServerError("Failed while last batch", err)
		}
	}
	return nil
}

// search  from ( EventQueue).
func GetAll[T any](query string) ([]T, rest_errors.RestErr) {

	logger.Info(fmt.Sprintf("abstractimpl|GetAll|query| %s", query))
	result, err := ezsearch.PostSearchResult(query)
	//s, _ := json.Marshal(result)
	//fmt.Println("abstractimpl|getall|", string(s))
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.ResultRow) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("No record found for %s", query))
	}
	res := getResultObjs[T](result.ResultRow)
	return res, nil
}

func getResultObjs[T any](rows []map[string]interface{}) []T {

	// load results
	res := make([]T, 0)
	//tableFields = ezsearch.GetBleveTableschema(EventQueueTable)
	for _, row := range rows {
		var eq T
		jsonStr, err := json.Marshal(row)
		if err != nil {
			logger.Error("abstractimpl|getResultobjs|Failed while get object", err)
			fmt.Println(err)
		}
		fmt.Println("getResultObjs|row", string(jsonStr))
		// Convert json string to struct
		if err := json.Unmarshal(jsonStr, &eq); err != nil {
			logger.Error("abstractimpl|Failed while Unmarshal ", err, zapcore.Field{String: string(jsonStr), Key: "p1", Type: zapcore.StringType})
			continue
		}
		res = append(res, eq)
	}
	return res
}

// Get a record  from ( EventQueue) .
func Get[T any](query string) (T, rest_errors.RestErr) {
	var m T
	result, err := ezsearch.PostSearchResult(query)
	if err != nil {
		return m, err
	}
	if result == nil || len(result.ResultRow) > 0 {
		return m, rest_errors.NewNotFoundError(fmt.Sprintf("abstractimpl|no record found %s", query))
	}
	res := getResultObjs[T](result.ResultRow)
	return res[0], nil
}

func BuilddynamicSchema(pindexName string) (bleve.Index, rest_errors.RestErr) {
	fdefs, err := ezsearch.GetBleveTableschema(pindexName)
	if err != nil {
		saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("index [%s] schema was not created", pindexName))
		logger.Error(fmt.Sprintf("Failed|BuilddynamicSchema %s", pindexName), saveErr)
	}
	if fdefs == nil || len(fdefs) == 0 {
		fdefs = append(fdefs, common.BleveFieldDef{Name: "timestamp", Type: "date"})
	}
	err = BuildIndexSchema(pindexName, fdefs, pindexName)
	i, err := ezsearch.GetIndex(pindexName)
	if err != nil {
		saveErr := rest_errors.NewBadRequestError(fmt.Sprintf("index [%s] is not created. Please try after the index first", pindexName))
		logger.Error(fmt.Sprintf("Failed|AddOrUpdateIndex  %s", pindexName), saveErr)
		return nil, saveErr

	}
	return i, nil
}
