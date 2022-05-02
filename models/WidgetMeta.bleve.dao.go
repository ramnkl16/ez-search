package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"

	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

//Auto code generated with help of mysql table schema
// table : WidgetMeta

// Insert the WidgetMeta to the database.
func (m WidgetMeta) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(m, abstractimpl.QueryMetaTable, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( MetaQuery) from the database.
func (m WidgetMeta) Delete(id string) rest_errors.RestErr {
	err := abstractimpl.Delete(abstractimpl.QueryMetaTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( MetaQuery).
func (m WidgetMeta) GetAll(query string) ([]WidgetMeta, rest_errors.RestErr) {

	if len(query) == 0 {
		query = fmt.Sprintf("select * from %s limit 0, 50", abstractimpl.QueryMetaTable)
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[WidgetMeta](query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Get a record  from ( MetaQuery) .
func (m WidgetMeta) Get(id string) (*WidgetMeta, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.QueryMetaTable, id)
	res, err := abstractimpl.Get[WidgetMeta](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
