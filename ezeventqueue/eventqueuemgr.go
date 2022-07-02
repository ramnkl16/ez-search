package ezeventqueue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/ezcsv"
	"github.com/ramnkl16/ez-search/ezmssqlconn"
	"github.com/ramnkl16/ez-search/ezsearch"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

func ProcessEventqueue() rest_errors.RestErr {

	events, err := services.EventQueueCustomService.GetActiveEventQueues() //fetch all events queue
	logger.Info("ProcessEventqueue")
	if err != nil && err.Status() == http.StatusNotFound {
		//fmt.Println("No Event found")
		logger.Error("Failed while process EventQueue", err)
		return err
	}
	quevedEvents := make(map[string]interface{})
	logger.Info("ProcessEventqueue")
	for _, e := range events {
		e.Status = int(global.STATUS_QUEUED)
		quevedEvents[e.ID] = e
	}
	abstractimpl.BatchCreateOrUpdate(abstractimpl.EventQueueTable, quevedEvents)
	for _, e := range events {
		start := time.Now()
		logger.Info(fmt.Sprintf("event Type:%s", e.EventType))
		e.Status = int(global.STATUS_INPROGRESS)
		updateEventQueue(&e)
		var err rest_errors.RestErr
		switch e.EventType {

		case global.EVENT_TYPE_INDEXFROMCSV:
			logger.Info(global.EVENT_TYPE_INDEXFROMCSV)
			//				e.StartAt= date_utils.GetNextScheduleDateByMins(e.StartAt, 2)
			//				e.StartAt= date_utils.GetNextScheduleDateByMins(e.StartAt, 2)
			err = executeCSVImport(&e)
			break
		case global.EVENT_TYPE_MSSQL_SYNC:
			logger.Info(fmt.Sprintf("case EVENT_TYPE_MSSQL_SYNC %s", e.EventType))
			err = ezmssqlconn.ExecuteMsSqlScript(&e)
			break
		case global.EVENT_TYPE_DETETE_LOG:
			logger.Info(fmt.Sprintf("case EVENT_TYPE_DETETE_LOG %s", e.EventType))
			err = executeDeleteIndexDocs(&e)
		default:
			logger.Info(fmt.Sprintf("event type is not implemented %s", e.EventType))
		}
		e.TimeTaken = int(time.Since(start))
		if err != nil {
			logger.Debug("Eventqueue failed section")
			if e.RecurringInSeconds == 0 {
				logger.Error("Failed while execute", err)
				e.Status = int(global.STATUS_ERROR)
				e.Message = err.Error()
				abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
			} else {
				logger.Debug("Eventqueue failed retry enabled")
				handleEventQueueRetry(&e)
			}
			continue
		}
		logger.Debug("Eventqueue completed section")
		//fmt.Println(e.RecurringInSeconds)
		e.Status = int(global.STATUS_COMPLETED)
		err = updateEventQueueHis(&e)
		if err != nil {
			logger.Error("Failed while updateEventQueueHis", err)
		}
		if e.RecurringInSeconds == 0 {
			abstractimpl.Delete(abstractimpl.EventQueueTable, e.ID)
		} else {
			logger.Debug("Eventqueue completed retry enabled")
			e.Status = int(global.STATUS_ACTIVE)
			e.StartAt = date_utils.GetNextScheduleDateBySeconds(e.StartAt, time.Duration(e.RecurringInSeconds))
			err = updateEventQueue(&e)
			if err != nil {
				logger.Error("Failed while updateEventQueue", err)
			}
		}

	}
	return nil
}
func handleEventQueueRetry(e *models.EventQueue) {
	rm := e.RetryMax
	if rm == 0 {
		rm = 5
	}
	if e.RetryCount > rm {
		e.Status = int(global.STATUS_SUSPEND)
		e.Message = "Reached max retry count"
	} else {
		e.RetryCount = e.RetryCount + 1
		e.Status = int(global.STATUS_ACTIVE)
		dur := e.RecurringInSeconds
		if dur <= 0 {
			dur = 5 * 60
		} else {
			dur = e.RecurringInSeconds
		}
		e.Status = int(global.STATUS_ACTIVE)
		e.StartAt = date_utils.GetNextScheduleDateBySeconds(e.StartAt, time.Duration(dur))
	}
	//abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
	updateEventQueue(e)
	updateEventQueueHis(e)
}

func executeCSVImport(e *models.EventQueue) rest_errors.RestErr {
	var ed CSVImportCustomData
	logger.Debug(fmt.Sprintf("eventqueue|%s", e.EventData))

	err1 := json.Unmarshal([]byte(e.EventData), &ed)
	if err1 != nil {
		logger.Error("Failed unmarshal", err1)
		e.Status = int(global.STATUS_ERROR)
		e.Message = fmt.Sprintf("Failed unmarshal %v", err1.Error())

		//abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
		return rest_errors.NewInternalServerError("Failed unmarshal", err1)
	}
	indexDocs, err := ezcsv.GetJsonFromCsv(ed.FileName, -1)
	if err != nil {
		logger.Error("Failed", err)
		e.Status = int(global.STATUS_ERROR)
		e.Message = fmt.Sprintf("Failed unmarshal %v", err1.Error())

		//abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
		return rest_errors.NewInternalServerError("Failed|getJsonFrom", err1)
	}
	err2 := abstractimpl.BatchCreateOrUpdate(ed.IndexName, indexDocs)
	if err2 != nil {
		logger.Error("Failed", err2)
		e.Status = int(global.STATUS_ERROR)
		e.Message = fmt.Sprintf("Failed unmarshal %v", err1.Error())
		e.StartAt = date_utils.GetNextScheduleDateBySeconds(e.StartAt, 2*60)
		e.Status = int(global.STATUS_ACTIVE)
		//abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
		updateEventQueueHis(e)
		return rest_errors.NewInternalServerError("Failed|getJsonFrom", err1)
	}
	return nil
} //func saveEventQueueHistory (EventQueue eq )
func updateEventQueueHis(eq *models.EventQueue) rest_errors.RestErr {
	eqh := models.EventQueueHistory{EventQueueID: eq.ID,
		EventType: eq.EventType, EventData: eq.EventData, Status: eq.Status, RetryCount: eq.RetryCount,
		RetryDuraitionInSeconds: eq.RecurringInSeconds, StartAt: eq.StartAt, LastSyncAt: eq.LastSyncAt, TimeTaken: eq.TimeTaken,
		Message: eq.Message, IsActive: eq.IsActive}
	eqh.ID = uid_utils.GetUid("eh", true)
	return abstractimpl.CreateOrUpdate(eqh, abstractimpl.EventQueueHisTable, eqh.ID)
}

func updateEventQueue(eq *models.EventQueue) rest_errors.RestErr {
	eq.UpdatedAt = date_utils.GetNowSearchFormat()
	return abstractimpl.CreateOrUpdate(eq, abstractimpl.EventQueueTable, eq.ID)
}

type CSVImportCustomData struct {
	FileName            string `json:"fileName"`    //consider relative path if it is not start with / char
	IgnoreEmpty         bool   `json:"ignoreEmpty"` //ingore field when empty
	IndexName           string `json:"indexName"`   //with relative path
	UniqueIndexColIndex int    `json:"uniqueIndexColIndex"`
}

type deleteIndexDocsCustomData struct {
	NoofDaysPersist int    `json:"noDays"`       //no of days to persist
	IndexNameKey    string `json:"indexNameKey"` //with relative path
}

func executeDeleteIndexDocs(e *models.EventQueue) rest_errors.RestErr {
	var ed deleteIndexDocsCustomData
	// ed1 := deleteIndexDocsCustomData{NoofDaysPersist: -1, IndexNameKey: "test.key"}
	// ss, _ := json.Marshal(ed1)
	// logger.Debug(fmt.Sprintf("eventqueues|%s", string(ss)))
	logger.Debug(fmt.Sprintf("eventqueue|%s", e.EventData))

	err1 := json.Unmarshal([]byte(e.EventData), &ed)
	if err1 != nil {
		logger.Error("Failed unmarshal", err1)
		e.Status = int(global.STATUS_ERROR)
		e.Message = fmt.Sprintf("Failed unmarshal %v", err1.Error())

		//abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
		return rest_errors.NewInternalServerError("Failed unmarshal", err1)
	}
	indexdocNames, err := coredb.GetValue(coredb.Defaultbucket, ed.IndexNameKey)
	if err != nil {
		logger.Error("Failed", err)
		return rest_errors.NewInternalServerError(fmt.Sprintf("Failed|while coredb key=%s", ed.IndexNameKey), err)
	}
	lastDt := time.Now().AddDate(0, 0, -ed.NoofDaysPersist).Format(date_utils.DateyyyymmddLayout)
	listToDelete := make([]string, 0) //list of index to delete
	indexes := common.GetAllIndexes()
	//fmt.Println("coredbvalue", string(indexdocNames))
	for _, indexName := range strings.Split(string(indexdocNames), ",") {
		//fmt.Println("indexName", indexName)
		if strings.Contains(indexName, "{") {
			patternIndexName, _ := common.GetPatternIndexName(indexName, lastDt)
			for k, _ := range indexes {
				//fmt.Println("k", k, "patternIndexName", patternIndexName)
				if k < patternIndexName {
					listToDelete = append(listToDelete, k)
				}
			}
		} else {
			listToDelete = append(listToDelete, indexName)
		}
	}
	//fmt.Println("executeDelteIndexDocs|listToDelete", listToDelete)
	err = ezsearch.DeleteIndexDocs(listToDelete)
	if err != nil {
		logger.Error("Failed", err)
		e.Status = int(global.STATUS_ERROR)
		e.Message = fmt.Sprintf("Failed unmarshal %v", err.Error())
		e.StartAt = date_utils.GetNextScheduleDateBySeconds(e.StartAt, 2*60)
		e.Status = int(global.STATUS_ACTIVE)
		updateEventQueueHis(e)
		return rest_errors.NewInternalServerError("Failed|getJsonFrom", err1)
	}
	return nil
}
