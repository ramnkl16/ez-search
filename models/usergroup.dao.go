package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

// Insert the EventQueue to the database.
func (m UserGroup) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(m, abstractimpl.UserGroupTable, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( EventQueue) from the database.
func (m UserGroup) Delete(id string) rest_errors.RestErr {
	err := abstractimpl.Delete(abstractimpl.EventQueueTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( EventQueue).
func (eq UserGroup) GetAll(query string) ([]UserGroup, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.UserGroupTable)
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[UserGroup](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( EventQueue) .
func (m UserGroup) Get(id string) (*UserGroup, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.UserGroupTable, id)
	res, err := abstractimpl.Get[UserGroup](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
