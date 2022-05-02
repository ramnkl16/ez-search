package services

import (
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

//Auto code generated with help of mysql table schema
// table : UserGroup

//EventQueue service as variable
var (
	UserGroupService userGroupServiceInterface = &userGroupService{}
)

type userGroupService struct{}

type userGroupServiceInterface interface {
	Get(string) (*models.UserGroup, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string) ([]models.UserGroup, rest_errors.RestErr)
	Save(models.UserGroup) rest_errors.RestErr
	//WebuiSearch(string, string) (models.userGroups, rest_errors.RestErr)
}

func (srv *userGroupService) Save(m models.UserGroup) rest_errors.RestErr {
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	if len(m.ID) == 0 { //consider insert
		m.ID = uid_utils.GetUid("um", false)
		m.IsActive = true
	}
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}

func (srv *userGroupService) Get(id string) (*models.UserGroup, rest_errors.RestErr) {
	var m *models.UserGroup
	m, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (srv *userGroupService) Delete(id string) rest_errors.RestErr {
	var m *models.UserGroup
	m, _ = m.Get(id)
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	m.IsActive = false
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}
func (srv *userGroupService) Search(query string) ([]models.UserGroup, rest_errors.RestErr) {
	m := models.UserGroup{}
	return m.GetAll(query)
}
