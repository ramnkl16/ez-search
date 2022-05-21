package services

import (
	"time"

	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
)

//Auto code generated with help of mysql table schema
// table : EventQueue

//EventQueue service as variable
var (
	EventQueueCustomService eventQueueCustomServiceInterface = &eventQueueCustomService{}
)

type eventQueueCustomService struct{}

type eventQueueCustomServiceInterface interface {
	RetryCountUpdate(eventId string, status int, intervalInMinutes int, errorMsg string) (int, rest_errors.RestErr)
	StatusUpdate(eventId string, status int, erroMsg string) rest_errors.RestErr
	GetActiveEventQueues() (models.EventQueues, rest_errors.RestErr)
	HardDelete(eventqueueId string) rest_errors.RestErr
	CreateWithIndex(models.EventQueue) rest_errors.RestErr
	UpdateWithIndex(models.EventQueue) rest_errors.RestErr
}

func (srv *eventQueueCustomService) CreateWithIndex(eq models.EventQueue) rest_errors.RestErr {
	eq.IsActive = "t"
	eq.CreatedAt = date_utils.GetNowSearchFormat()
	eq.UpdatedAt = date_utils.GetNowSearchFormat()
	if len(eq.StartAt) == 0 {
		eq.StartAt = date_utils.GetNowSearchFormat()
	}
	if err := eq.CreateOrUpdate(); err != nil {
		return err
	}
	//updateEventIndex(&eq)
	return nil
}

// func updateEventIndex(eq *models.EventQueue) {
// 	m := ezsearch.SearchEventModelClient{Id: eq.ID, Status: fmt.Sprintf("%v", global.StatusEnum(eq.Status)), EventType: eq.EventTypeID, EventData: eq.EventData, RetryCount: eq.RetryCount, IsActive: eq.IsActive, Message: eq.Message}
// 	m.StartAt = eq.StartAt
// 	m.CreatedAt = eq.CreatedAt
// 	m.UpdatedAt = eq.UpdatedAt
// 	err := ezsearch.EventIndex.Index(m.Id, m)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Failed while update |updateEventIndex| %s %v", m.EventData, m.CreatedAt), err)
// 	}
// }

// func getEventIndexById(id string) *ezsearch.SearchEventModelClient {
// 	//rm := ezsearch.SearchRequestModel{PharseQueries: []string{fmt.Sprintf(`"%s"`, id)}, Fields: []string{"*"}, IndexName: ezsearch.Conf.EventIndexpath, From: 0, Size: 1}
// 	queryStr := fmt.Sprintf("select * from %s where id:%s limit 0,1", ezsearch.Conf.EventIndexpath, id)
// 	//fmt.Println("getEventIndexById|SearchRequestModel", rm)
// 	resM, err := ezsearch.PostSearchResult(queryStr)
// 	if err != nil {
// 		logger.Error("Failed while fetch indexed value|getEventIndexById ", err)
// 		return nil
// 	}
// 	if resM == nil {
// 		logger.Info(fmt.Sprintf("Failed while fetch indexed event %s", id))
// 		e, err := EventQueueService.Get(id)
// 		if err != nil {
// 			logger.Error(fmt.Sprintf("Failed while fetch event from db %s", id), err)
// 		}
// 		m := ezsearch.SearchEventModelClient{EventData: e.EventData, Id: e.ID,
// 			EventType: e.EventTypeID, StartAt: e.StartAt, CreatedAt: e.CreatedAt,
// 			UpdatedAt: e.UpdatedAt, Status: fmt.Sprintf("%v", global.StatusEnum(e.Status)), Message: e.Message}
// 		m.RetryCount = e.RetryCount
// 		m.IsActive = e.IsActive
// 		return &m
// 	}
// 	//str, _ := json.Marshal(resM)
// 	//fmt.Println("getEventIndexById|string", string(str))
// 	row := resM.ResultRow[0]
// 	m := ezsearch.SearchEventModelClient{EventData: getInterfaceValToStr(row["eventData"]), Id: getInterfaceValToStr(row["id"]),
// 		EventType: getInterfaceValToStr(row["eventType"]), StartAt: getInterfaceValToStr(row["startAt"]), CreatedAt: getInterfaceValToStr(row["cretedAt"]),
// 		UpdatedAt: getInterfaceValToStr(row["updatedAt"]), Status: getInterfaceValToStr(row["status"]), Message: getInterfaceValToStr(row["message"])}
// 	m.RetryCount, _ = strconv.Atoi(getInterfaceValToStr(row["retryCount"]))
// 	m.IsActive, _ = strconv.ParseBool(getInterfaceValToStr(row["isActive"]))
// 	return &m
// }
// func getInterfaceValToStr(val interface{}) string {
// 	return fmt.Sprintf("%v", val)
// }

func (srv *eventQueueCustomService) UpdateWithIndex(eq models.EventQueue) rest_errors.RestErr {
	eq.UpdatedAt = date_utils.GetNowSearchFormat()
	if err := eq.CreateOrUpdate(); err != nil {
		return err
	}
	//updateEventIndex(&eq)
	return nil
}

func (srv *eventQueueCustomService) RetryCountUpdate(id string, status int, intervalInMinutes int, errorMsg string) (int, rest_errors.RestErr) {
	dao := models.EventQueue{ID: id, Status: status, Message: errorMsg, StartAt: date_utils.GetNextScheduleDate(date_utils.GetNowSearchFormat(),
		time.Duration(intervalInMinutes)), UpdatedAt: date_utils.GetNowSearchFormat()}
	i, err := dao.UpdateRetryCount()
	return i, err
}

func (srv *eventQueueCustomService) HardDelete(id string) rest_errors.RestErr {
	dao := models.EventQueue{ID: id}
	// m := getEventIndexById(id)
	// m.Status = fmt.Sprintf("%v", global.StatusEnum(global.STATUS_DELETED))
	// ezsearch.EventIndex.Index(m.Id, m)
	err := dao.HardDeleteFromIndex()
	return err
}

func (srv *eventQueueCustomService) StatusUpdate(id string, status int, msg string) rest_errors.RestErr {
	dao := models.EventQueue{ID: id, Status: status, Message: msg, UpdatedAt: date_utils.GetNowString()}
	err := dao.UpdateStatusIndex()

	return err
}

func (srv *eventQueueCustomService) GetActiveEventQueues() (models.EventQueues, rest_errors.RestErr) {
	dao := models.EventQueue{}
	result, err := dao.GetActiveEventQueuesFromIndex()
	return result, err
}
