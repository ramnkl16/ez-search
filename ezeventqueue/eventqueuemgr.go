package ezeventqueue

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		logger.Info(fmt.Sprintf("event Type:%s", e.EventType))
		switch e.EventType {
		case global.EVENT_TYPE_INDEXFROMCSV:
			logger.Info(global.EVENT_TYPE_INDEXFROMCSV)
			//				e.StartAt= date_utils.GetNextScheduleDateByMins(e.StartAt, 2)
			//				e.StartAt= date_utils.GetNextScheduleDateByMins(e.StartAt, 2)
			e.Status = int(global.STATUS_INPROGRESS)
			abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
			err := executeCSVImport(&e)

			if err != nil {
				logger.Error("Failed while execute", err)
				e.Status = int(global.STATUS_ERROR)
				e.Message = err.Error()
				abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
				continue
			}

			abstractimpl.Delete(abstractimpl.EventQueueTable, e.ID)
			fmt.Println("after deleted")
			e.Status = int(global.STATUS_COMPLETED)
			err = updateEventQueueHis(&e)
			if err != nil {
				logger.Error("Failed while updateEventQueueHis", err)
			}
			break
		case global.EVENT_TYPE_MSSQL_SYNC:
			logger.Info(fmt.Sprintf("case EVENT_TYPE_MSSQL_SYNC %s", e.EventType))
			e.Status = int(global.STATUS_INPROGRESS)
			abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
			err := ezmssqlconn.ExecuteMsSqlScript(&e)

			if err != nil {
				logger.Error("Failed while execute", err)
				e.Status = int(global.STATUS_ERROR)
				e.Message = err.Error()
				abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
				continue
			}

			abstractimpl.Delete(abstractimpl.EventQueueTable, e.ID)
			fmt.Println("after deleted")
			e.Status = int(global.STATUS_COMPLETED)
			err = updateEventQueueHis(&e)
			if err != nil {
				logger.Error("Failed while updateEventQueueHis", err)
			}
			break
		default:
			logger.Info(fmt.Sprintf("event type is not implemented %s", e.EventType))
		}
	}
	return nil
}
func handleEventQueueRetry(e models.EventQueue) {
	e.RetryCount = e.RetryCount + 1
	if e.RetryCount > 5 {
		e.Status = int(global.STATUS_SUSPEND)
		e.Message = "Reached max retry count"
	} else {
		e.Status = int(global.STATUS_ACTIVE)
		e.StartAt = date_utils.GetNextScheduleDateByMins(e.StartAt, 5)
	}
	abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
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
	eqh := models.EventQueueHistory{EventQueueID: eq.ID, EventTypeID: eq.EventType, EventData: eq.EventData, Status: eq.Status, RetryCount: eq.RetryCount, Message: eq.Message, IsActive: eq.IsActive}
	eqh.ID = uid_utils.GetUid("eh", true)
	return abstractimpl.CreateOrUpdate(eqh, abstractimpl.EventQueueHisTable, eqh.ID)
}

type CSVImportCustomData struct {
	FileName            string `json:"fileName"`    //consider relative path if it is not start with / char
	IgnoreEmpty         bool   `json:"ignoreEmpty"` //ingore field when empty
	IndexName           string `json:"indexName"`   //with relative path
	UniqueIndexColIndex int    `json:"uniqueIndexColIndex"`
}
