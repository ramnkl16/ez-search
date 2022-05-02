package services

import (
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

//Auto code generated with help of mysql table schema
// table : UserMenu

//EventQueue service as variable
var (
	UserMenuService userMenuServiceInterface = &userMenuService{}
)

type userMenuService struct{}

type userMenuServiceInterface interface {
	Get(string) (*models.UserMenu, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string) (models.UserMenus, rest_errors.RestErr)
	Save(models.UserMenu) rest_errors.RestErr
	//WebuiSearch(string, string) (models.UserMenus, rest_errors.RestErr)
}

func (srv *userMenuService) Save(m models.UserMenu) rest_errors.RestErr {
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

func (srv *userMenuService) Get(id string) (*models.UserMenu, rest_errors.RestErr) {
	var m *models.UserMenu
	m, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (srv *userMenuService) Delete(id string) rest_errors.RestErr {
	var m *models.UserMenu
	m, _ = m.Get(id)
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	m.IsActive = false
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}
func (srv *userMenuService) Search(query string) (models.UserMenus, rest_errors.RestErr) {
	m := models.UserMenu{}
	return m.GetAll(query)
}
