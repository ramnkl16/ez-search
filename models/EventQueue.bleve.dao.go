package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

// Insert the EventQueue to the database.
func (eq EventQueue) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(eq, abstractimpl.EventQueueTable, eq.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( EventQueue) from the database.
func (eq EventQueue) Delete(id string) rest_errors.RestErr {
	err := abstractimpl.Delete(abstractimpl.EventQueueTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( EventQueue).
func (eq EventQueue) GetAll(query string) ([]EventQueue, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.EventQueueTable)
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[EventQueue](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( EventQueue) .
func (eq EventQueue) Get(id string) (*EventQueue, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.EventQueueTable, id)
	res, err := abstractimpl.Get[EventQueue](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
