package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

// Insert the EventQueue to the database.
func (m Namespace) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(m, abstractimpl.NamespaceTable, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( EventQueue) from the database.
func (m Namespace) Delete(id string) rest_errors.RestErr {
	err := abstractimpl.Delete(abstractimpl.EventQueueTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( EventQueue).
func (eq Namespace) GetAll(query string) ([]Namespace, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.NamespaceTable)
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[Namespace](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( EventQueue) .
func (m Namespace) Get(id string) (*Namespace, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.NamespaceTable, id)
	res, err := abstractimpl.Get[Namespace](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
