package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

// Insert the EventQueue to the database.
func (m UserMenu) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(m, abstractimpl.UserMenuTable, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( EventQueue) from the database.
func (m UserMenu) Delete(id string) rest_errors.RestErr {
	err := abstractimpl.Delete(abstractimpl.UserMenuTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( EventQueue).
func (eq UserMenu) GetAll(query string) ([]UserMenu, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.UserMenuTable)
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[UserMenu](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( EventQueue) .
func (m UserMenu) Get(id string) (*UserMenu, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.UserMenuTable, id)
	res, err := abstractimpl.Get[UserMenu](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
