package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

//Auto code generated with help of mysql table schema
// table : Menu

// Insert the Menu to the database.
func (wm Menu) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(wm, abstractimpl.MenuTable, wm.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( Menu) from the database.
func (m Menu) Delete(id string) rest_errors.RestErr {
	err := abstractimpl.Delete(abstractimpl.MenuTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( Menu).
func (m Menu) GetAll(query string) ([]Menu, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.MenuTable)
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[Menu](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( Menu) .
func (m Menu) Get(id string) (*Menu, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.MenuTable, id)
	res, err := abstractimpl.Get[Menu](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
