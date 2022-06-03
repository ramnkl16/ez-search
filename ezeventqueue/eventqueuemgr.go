package ezeventqueue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/ezcsv"
	"github.com/ramnkl16/ez-search/ezmssqlconn"
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
		fmt.Println("No Event found")
		return err
	}
	quevedEvents := make(map[string]interface{})
	logger.Info("ProcessEventqueue|25")
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
		abstractimpl.Delete(abstractimpl.EventQueueTable, e.ID)
		fmt.Println("after deleted")
		e.Status = int(global.STATUS_COMPLETED)
		err = updateEventQueueHis(&e)
		if err != nil {
			logger.Error("Failed while updateEventQueueHis", err)
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
			dur = 5
		} else {
			dur = e.RecurringInSeconds / 60
		}
		e.Status = int(global.STATUS_ACTIVE)
		e.StartAt = date_utils.GetNextScheduleDateByMins(e.StartAt, time.Duration(dur))
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
		e.StartAt = date_utils.GetNextScheduleDateByMins(e.StartAt, 2)
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
