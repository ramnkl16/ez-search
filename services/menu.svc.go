package services

import (
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

//Auto code generated with help of mysql table schema
// table : Menu

//EventQueue service as variable
var (
	MenuService menuServiceInterface = &menuService{}
)

type menuService struct{}

type menuServiceInterface interface {
	Get(string) (*models.Menu, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string) (models.Menus, rest_errors.RestErr)
	Save(models.Menu) rest_errors.RestErr
	//WebuiSearch(string, string) (models.Menus, rest_errors.RestErr)
}

func (srv *menuService) Save(m models.Menu) rest_errors.RestErr {
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	if len(m.ID) == 0 { //consider insert
		m.ID = uid_utils.GetUid("mn", false)
		m.IsActive = true
	}
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil

}

func (srv *menuService) Get(id string) (*models.Menu, rest_errors.RestErr) {
	var m *models.Menu
	m, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (srv *menuService) Delete(id string) rest_errors.RestErr {
	var m *models.Menu
	m, _ = m.Get(id)
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	m.IsActive = false
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}
func (srv *menuService) Search(query string) (models.Menus, rest_errors.RestErr) {
	m := models.Menu{}
	return m.GetAll(query)
}
