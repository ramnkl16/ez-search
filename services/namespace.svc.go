package services

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/crypto_utils"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

//Auto code generated with help of mysql table schema
// table : Namespace

//EventQueue service as variable
var (
	NamespaceService namespaceServiceInterface = &namespaceService{}
)

type namespaceService struct{}

type namespaceServiceInterface interface {
	Get(string) (*models.Namespace, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string) ([]models.Namespace, rest_errors.RestErr)
	Save(models.Namespace) rest_errors.RestErr
	New(models.NamespaceParam) rest_errors.RestErr

	//WebuiSearch(string, string) (models.Namespaces, rest_errors.RestErr)
}

func (srv *namespaceService) New(m models.NamespaceParam) rest_errors.RestErr {

	// var na *models.Namespace
	na := &models.Namespace{
		Name:         m.Name,
		Code:         m.Code,
		ContextToken: m.ContextToken,
	}
	na.Name = m.Name
	na.ContextToken = m.ContextToken
	na.ID = uid_utils.GetUid("ns", false)
	na.UpdatedAt = date_utils.GetNowSearchFormat()
	na.CreatedAt = date_utils.GetNowSearchFormat()
	//na.UpdatedBy = ai.UserId
	na.IsActive = "t"

	if !global.IsValidEmail(m.Email) {
		err := rest_errors.NewBadRequestError("Invalid Email")
		logger.Error("Invalid Email Provided for Namespace Creation", err)
		return err
	}

	if !global.IsValidPhoneNumber(m.Mobile) {
		err := rest_errors.NewBadRequestError("Invalid Mobile Number")
		logger.Error("Invalid Email Provided for Namespace Creation", err)
		return err
	}
	ug := models.UserGroup{Name: "admin", Description: "default admin", NamespaceID: na.ID}
	ug.UpdatedAt = date_utils.GetNowSearchFormat()
	//ug.UpdatedBy = ai.UserId
	ug.IsActive = "t"
	ug.ID = uid_utils.GetUid("ug", false)

	//User menu initialization
	um := models.UserMenu{
		NamespaceID:       na.ID,
		MenuID:            global.RootMenuID,
		ReferenceID:       na.ID,
		ReferenceType:     global.RefTypeNamespace,
		Privilege:         global.DefaultMenuPermission,
		PermissionPlus:    0,
		MenuExceptionFlag: 1,
		CreatedAt:         date_utils.GetNowSearchFormat(),
		UpdatedAt:         date_utils.GetNowSearchFormat(),
	}
	um.IsActive = "t"
	hashPwd := crypto_utils.GetMd5(global.DefaultUserPassword)
	um.UpdatedAt = date_utils.GetNowSearchFormat()
	um.ID = uid_utils.GetUid("um", false)

	us := models.User{
		NamespaceID:       na.ID,
		UserGroupID:       ug.ID,
		UserName:          m.Email,
		Token:             hashPwd,
		Email:             m.Email,
		Mobile:            m.Mobile,
		FirstName:         global.DefaultUserFirstName,
		LastName:          global.DefaultLastName,
		RoleId:            "ROLE1",
		EmailVerified:     date_utils.GetNowSearchFormat(),
		PasswordUpdatedAt: date_utils.GetNowSearchFormat(),
		CreatedAt:         date_utils.GetNowSearchFormat(),
		UpdatedAt:         date_utils.GetNowSearchFormat(),
	}
	us.IsActive = "t"
	us.UpdatedAt = date_utils.GetNowSearchFormat()
	us.CreatedAt = us.UpdatedAt
	us.ID = uid_utils.GetUid("ur", false)

	u := models.User{
		NamespaceID:       m.Code,
		UserGroupID:       ug.ID,
		UserName:          m.Email,
		Token:             hashPwd,
		Email:             "intuser@gost.com",
		Mobile:            m.Mobile,
		FirstName:         global.DefaultUserFirstName,
		LastName:          global.DefaultLastName,
		RoleId:            "ROLE1",
		EmailVerified:     date_utils.GetNowSearchFormat(),
		PasswordUpdatedAt: date_utils.GetNowSearchFormat(),
		CreatedAt:         date_utils.GetNowSearchFormat(),
		UpdatedAt:         date_utils.GetNowSearchFormat(),
	}
	u.IsActive = "t"
	u.UpdatedAt = date_utils.GetNowSearchFormat()
	u.CreatedAt = us.UpdatedAt
	u.ID = uid_utils.GetUid("ur", false)
	q := models.WidgetMeta{
		ID:        "q1",
		Name:      fmt.Sprintf("User List %s", m.Code),
		Division:  m.Code,
		Module:    "search",
		Page:      "user",
		Data:      fmt.Sprintf("select * from tables/tables.user"),
		IsActive:  "t",
		CreatedAt: date_utils.GetNowSearchFormat(),
		UpdatedAt: date_utils.GetNowSearchFormat(),
	}
	if m.Code != "platform" {
		q.Data = fmt.Sprintf("select * from tables/tables.user where namespaceId:%s", m.Code)
	}

	r := models.WidgetMeta{ID: uid_utils.GetUid("rp", true), Name: "applogs", Division: m.Code, Module: "mod", Page: "report", Data: "select * from indexes/applogs-{2006-01-02} since t:30 seconds ago sort -t facets l", IsActive: "t", CreatedAt: date_utils.GetNowSearchFormat(), UpdatedAt: date_utils.GetNowSearchFormat()}

	abstractimpl.CreateOrUpdate(na, abstractimpl.NamespaceTable, na.ID)
	abstractimpl.CreateOrUpdate(us, abstractimpl.UserTable, us.ID)
	abstractimpl.CreateOrUpdate(ug, abstractimpl.UserGroupTable, ug.ID)
	abstractimpl.CreateOrUpdate(um, abstractimpl.UserMenuTable, um.ID)
	abstractimpl.CreateOrUpdate(u, abstractimpl.UserTable, u.ID)
	abstractimpl.CreateOrUpdate(q, abstractimpl.QueryMetaTable, q.ID)
	abstractimpl.CreateOrUpdate(r, abstractimpl.QueryMetaTable, r.ID)
	return nil
}

func (srv *namespaceService) Save(m models.Namespace) rest_errors.RestErr {
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	if len(m.ID) == 0 { //consider insert
		m.ID = uid_utils.GetUid("ns", false)
		m.IsActive = "t"
	}
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}

func (srv *namespaceService) Get(id string) (*models.Namespace, rest_errors.RestErr) {
	var m *models.Namespace
	m, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func (srv *namespaceService) Delete(id string) rest_errors.RestErr {
	var m *models.Namespace
	m, _ = m.Get(id)
	m.UpdatedAt = date_utils.GetNowSearchFormat()
	m.IsActive = "f"
	if err := m.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}
func (srv *namespaceService) Search(query string) ([]models.Namespace, rest_errors.RestErr) {
	m := models.Namespace{}
	return m.GetAll(query)
}
