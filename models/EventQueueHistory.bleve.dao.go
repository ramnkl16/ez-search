package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

//Auto code generated with help of mysql table schema
// table : EventQueueHistory

// Insert the EventQueueHistory to the database.
func (m EventQueueHistory) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(m, abstractimpl.EventQueueHisTable, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( EventQueueHistory) from the database.
func (eq EventQueueHistory) Delete(id string) rest_errors.RestErr {

	err := abstractimpl.Delete(abstractimpl.EventQueueHisTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( EventQueueHistory).
func (EventQueueHistory) GetAll(query string) ([]EventQueueHistory, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.EventQueueHisTable)
	}

	logger.Info(query)
	res, err := abstractimpl.GetAll[EventQueueHistory](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( EventQueue) .
func (eq EventQueueHistory) Get(id string) (*EventQueueHistory, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.EventQueueHisTable, id)
	res, err := abstractimpl.Get[EventQueueHistory](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
