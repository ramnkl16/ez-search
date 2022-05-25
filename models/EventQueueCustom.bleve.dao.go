package models

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
)

func (eq *EventQueue) UpdateRetryCount() (int, rest_errors.RestErr) {
	m, err := eq.Get(eq.ID)
	if err != nil {
		return 0, err
	}
	//row := productsyncmysql.Client.QueryRow("call EventqueueUpdateRetryCount(?,?,?,?)", eq.ID, eq.Status, eq.StartAt, eq.UpdatedAt)
	m.RetryCount = m.RetryCount + 1
	m.UpdatedAt = eq.UpdatedAt
	m.StartAt = eq.StartAt
	m.Status = eq.Status
	t, err := abstractimpl.GetTable(abstractimpl.EventQueueTable)
	if err != nil {
		return 0, err
	}
	err1 := t.Index(m.ID, m)
	if err1 != nil {
		return 0, rest_errors.NewInternalServerError("Failed while update retry count", err)
	}

	return m.RetryCount, nil
}

func (eq *EventQueue) UpdateStatusIndex() rest_errors.RestErr {
	m, err := eq.Get(eq.ID)
	if err != nil {
		return err
	}
	//row := productsyncmysql.Client.QueryRow("call EventqueueUpdateRetryCount(?,?,?,?)", eq.ID, eq.Status, eq.StartAt, eq.UpdatedAt)
	m.Status = eq.Status
	m.UpdatedAt = eq.UpdatedAt
	t, err := abstractimpl.GetTable(abstractimpl.EventQueueTable)
	if err != nil {
		return err
	}
	err1 := t.Index(m.ID, m)
	if err1 != nil {
		return rest_errors.NewInternalServerError("Failed while update retry count", err)
	}
	return nil
}

func (eq *EventQueue) HardDeleteFromIndex() rest_errors.RestErr {
	m, err := eq.Get(eq.ID)
	if err != nil {
		return err
	}
	t, err := abstractimpl.GetTable(abstractimpl.EventQueueTable)
	if err != nil {
		return err
	}
	err1 := t.Delete(m.ID)
	if err1 != nil {
		return rest_errors.NewInternalServerError("Failed while update retry count", err)
	}

	return nil
}

// search  from ( EventQueue).
func (eq *EventQueue) GetActiveEventQueuesFromIndex() (EventQueues, rest_errors.RestErr) {

	currDt := date_utils.GetNowSearchFormat()
	query := fmt.Sprintf("select * from %s where status:%d, startAt:>=%s limit 0, 10000", abstractimpl.EventQueueTable, global.STATUS_ACTIVE, currDt)
	return eq.GetAll(query)
	//fmt.Println("GetActiveEventQueues currDt", currDt)

}
