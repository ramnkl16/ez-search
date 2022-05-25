package ezeventqueue

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/ezcsv"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
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
			var ed CSVImportCustomData
			logger.Debug(fmt.Sprintf("eventqueue|%s", e.EventData))

			err1 := json.Unmarshal([]byte(e.EventData), &ed)
			if err1 != nil {
				logger.Error("Failed unmarshal", err1)
				e.Status = int(global.STATUS_ERROR)
				e.Message = fmt.Sprintf("Failed unmarshal %v", err1.Error())
				abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, e.ID)
				continue
			}
			indexDocs, err := ezcsv.GetJsonFromCsv(ed.FileName, -1)
			if err != nil {
				logger.Error("Failed", err)
				continue
			}
			err2 := abstractimpl.BatchCreateOrUpdate(ed.IndexName, indexDocs)
			if err2 != nil {
				logger.Error("Failed", err2)
			}
			break
		}
	}
	return nil
}

type CSVImportCustomData struct {
	FileName            string `json:"fileName"`    //consider relative path if it is not start with / char
	IgnoreEmpty         bool   `json:"ignoreEmpty"` //ingore field when empty
	IndexName           string `json:"indexName"`   //with relative path
	UniqueIndexColIndex int    `json:"uniqueIndexColIndex"`
}
