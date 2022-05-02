package services

import (
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
)

//Auto code generated with help of mysql table schema
// table : EventQueue

//EventQueue service as variable
var (
	EventQueueService eventQueueServiceInterface = &eventQueueService{}
)

type eventQueueService struct{}

type eventQueueServiceInterface interface {
	Get(string) (*models.EventQueue, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string) (models.EventQueues, rest_errors.RestErr)
	//WebuiSearch(string, string) (models.EventQueues, rest_errors.RestErr)
}

func (srv *eventQueueService) Get(id string) (*models.EventQueue, rest_errors.RestErr) {
	var m *models.EventQueue
	m, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (srv *eventQueueService) Delete(id string) rest_errors.RestErr {
	var m *models.EventQueue
	m, _ = m.Get(id)
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	m.IsActive = false
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}
func (srv *eventQueueService) Search(query string) (models.EventQueues, rest_errors.RestErr) {
	m := models.EventQueue{}
	return m.GetAll(query)
}
