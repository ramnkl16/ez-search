package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"

	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
)

//Auto code generated with help of mysql table schema
// table : User

// Insert the User to the database.
func (m *User) CreateOrUpdate() rest_errors.RestErr {
	err := abstractimpl.CreateOrUpdate(m, abstractimpl.UserTable, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a recored ( User) from the database.
func (eq User) Delete(id string) rest_errors.RestErr {

	err := abstractimpl.Delete(abstractimpl.UserTable, id)
	if err != nil {
		return err
	}
	return nil
}

// search  from ( User).
func (m User) GetAll(query string) ([]User, rest_errors.RestErr) {

	if len(query) == 0 {
		return nil, rest_errors.NewBadRequestError("No search query")
	}
	logger.Info(query)
	res, err := abstractimpl.GetAll[User](query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get a record  from ( User) .
func (m User) Get(id string) (*User, rest_errors.RestErr) {
	query := fmt.Sprintf("select * from %s where id:%s limit 0, 1", abstractimpl.UserTable, id)
	res, err := abstractimpl.Get[User](query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
