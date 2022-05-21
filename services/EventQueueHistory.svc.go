package services

import (
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
)

//Auto code generated with help of mysql table schema
// table : EventQueueHistory

//EventQueueHistory service as variable
var (
	EventQueueHistoryService eventQueueHistoryServiceInterface = &eventQueueHistoryService{}
)

type eventQueueHistoryService struct{}

type eventQueueHistoryServiceInterface interface {
	Create(models.EventQueueHistory) rest_errors.RestErr
	Update(models.EventQueueHistory) rest_errors.RestErr
	Get(string) (*models.EventQueueHistory, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string, string) (models.EventQueueHistories, rest_errors.RestErr)
	//WebuiSearch(string, string) (models.EventQueueHistories, rest_errors.RestErr)
}

func (srv *eventQueueHistoryService) Create(eq models.EventQueueHistory) rest_errors.RestErr {

	eq.IsActive = "t"
	eq.CreatedAt = date_utils.GetNowSearchFormat()
	eq.UpdatedAt = date_utils.GetNowSearchFormat()
	// if err := eq.Create(); err != nil {
	// 	return err
	// }
	return nil
}

func (srv *eventQueueHistoryService) Update(eq models.EventQueueHistory) rest_errors.RestErr {
	eq.UpdatedAt = date_utils.GetNowSearchFormat()
	// if err := eq.Update(); err != nil {
	// 	return err
	// }
	return nil
}

func (srv *eventQueueHistoryService) Get(id string) (*models.EventQueueHistory, rest_errors.RestErr) {
	dao := &models.EventQueueHistory{ID: id}
	// if err := dao.Get(); err != nil {
	// 	return nil, err
	// }
	return dao, nil

}

func (srv *eventQueueHistoryService) Delete(id string) rest_errors.RestErr {
	dao := &models.EventQueueHistory{ID: id}
	dao.UpdatedAt = date_utils.GetNowSearchFormat()

	// if err := dao.Delete(); err != nil {
	// 	return err
	// }
	return nil
}
func (srv *eventQueueHistoryService) Search(start string, limit string) (models.EventQueueHistories, rest_errors.RestErr) {
	// dao := &models.EventQueueHistory{}
	// if start == "" {
	// 	start = "0"
	// }
	// if limit == "" {
	// 	limit = "50"
	// }
	// return dao.Searh(start, limit)
	return nil, nil
}
