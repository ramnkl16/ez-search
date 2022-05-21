package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/auth"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
	"github.com/ramnkl16/ez-search/utils/crypto_utils"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

var (
	UserService userServiceInterface = &userService{}
)

type userService struct{}

type userServiceInterface interface {
	Save(models.User) (*string, rest_errors.RestErr)
	Get(string, auth.Header) (*models.User, rest_errors.RestErr)
	Search(string, auth.Header) (models.Users, rest_errors.RestErr)
	Login(emailOrMobile string, password string, namespaceId string) (*models.User, rest_errors.RestErr)
	ChangePassword(emailOrMobile string, oldPassword string, password string, namespaceId string) (*models.User, rest_errors.RestErr)
	Logout(authToken string) rest_errors.RestErr
	Delete(string, auth.Header) rest_errors.RestErr
}

func (srv *userService) Save(us models.User) (*string, rest_errors.RestErr) {
	us.UpdatedAt = date_utils.GetNowSearchFormat()
	if len(us.ID) == 0 { //consider insert
		us.ID = uid_utils.GetUid("us", false)
		us.Token = crypto_utils.GetMd5(us.Token) //asumming while create first time token field would hae clean text password
		//us.ActiveFlag= 1 //TODO to be modified as enum type
		//var m models.User
	}
	if err := us.CreateOrUpdate(); err != nil {
		return nil, err
	}

	//cacheLoader.LoadUserData(&user)
	return &us.ID, nil

}

func (srv *userService) Get(id string, h auth.Header) (*models.User, rest_errors.RestErr) {

	//var al models.AuditLog

	// ai, err := auth.GetAuthInfo(h.HeaderXauth)
	// if err != nil {
	// 	logger.Error("Failed|User", err)
	// 	return nil, err
	// }

	//namespaceID := ai.NamespaceId
	// for _, ns := range global.GetAllGlobalAdminNS() {
	// 	fmt.Println("Indie if", ns)
	// 	// fmt.Println("Indie if", pp.Namespace)
	// 	if ns == ai.NamespaceId {
	// 		// fmt.Println("Indie if if", pp.Namespace)
	// 		namespaceID = h.HeaderNS
	// 		break
	// 	}
	// }

	//return if present in User cache
	// cacheResult := cacheLoader.GetUserData(id, namespaceID)
	// if cacheResult != nil {
	// 	fmt.Println(cacheResult)
	// 	fmt.Println("Found User in Cache")
	// 	if configuration.Config.EnableAuditLog {
	// 		al.Create()
	// 	}
	// 	return cacheResult, nil
	// }
	var m models.User
	m.ID = id
	result, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (srv *userService) Search(query string, h auth.Header) (models.Users, rest_errors.RestErr) {

	// ai, err := auth.GetAuthInfo(h.HeaderXauth)
	// if err != nil {
	// 	logger.Error("Failed|User", err)
	// 	return nil, err
	// }
	//fmt.Println(ai.NamespaceId)
	// namespaceID := ai.NamespaceId
	// for _, ns := range global.GetAllGlobalAdminNS() {
	// 	fmt.Println("Indie if", ns)
	// 	// fmt.Println("Indie if", pp.Namespace)
	// 	if ns == ai.NamespaceId {
	// 		// fmt.Println("Indie if if", pp.Namespace)
	// 		namespaceID = h.HeaderNS
	// 		break
	// 	}
	// }
	m := models.User{}
	return m.GetAll(query)

}
func (srv *userService) Login(emailOrMobile string, password string, nsCode string) (*models.User, rest_errors.RestErr) {

	//dao.Token = password
	var email, mobile string
	if !strings.Contains(emailOrMobile, "@") && len(emailOrMobile) < 10 {
		return nil, rest_errors.NewUnauthorizedError("Invalid email or mobile")
	}
	if len(password) == 0 {
		return nil, rest_errors.NewUnauthorizedError("Invalid password")
	}

	if strings.Contains(emailOrMobile, "@") {
		email = emailOrMobile
	} else {
		mobile = emailOrMobile
	}
	// sqlnsStr := fmt.Sprintf("SELECT * FROM  %s  WHERE isActive:true,code:%s", abstractimpl.NamespaceTable, nsCode)
	// nslsit, err := NamespaceService.Search(sqlnsStr)
	// if err != nil {
	// 	return nil, err
	// }
	// ns := nslsit[0].Code
	// nsId := nslsit[0].ID
	us := models.User{}
	sqlstr := fmt.Sprintf("SELECT * FROM  %s  WHERE +isActive:t,+namespaceId:%s,+email:%s ", abstractimpl.UserTable, nsCode, email)
	if len(mobile) > 0 && len(email) == 0 {
		sqlstr = fmt.Sprintf("SELECT * FROM  %s  WHERE +isActive:t,+namespaceId:%s,+mobile:%s ", abstractimpl.UserTable, nsCode, mobile)
	}
	//TODO never reach this scenario
	// else if len(mobile) > 0 && len(email) > 0 {
	// 	sqlstr = fmt.Sprintf("SELECT * FROM  %s  WHERE isActive:1, email:%s,mobil:%s", abstractimpl.UserTable, email, mobile)
	// }

	list, err := us.GetAll(sqlstr)
	if err != nil {
		return nil, err
	}
	//fmt.Println("length of user", len(list), list)
	if len(list) == 0 {
		return nil, rest_errors.NewUnauthorizedError("user name or passwrod is not matching")
	}
	if list[0].Token != crypto_utils.GetMd5(password) {
		return nil, rest_errors.NewUnauthorizedError("Invalid Password")
	}
	u := list[0]
	bu := models.UserBase64{UserName: emailOrMobile, Namespace: nsCode, Password: password}
	buStr, err1 := json.Marshal(bu)
	if err1 != nil {
		return nil, rest_errors.NewBadRequestError(err1.Error())
	}
	enToken := b64.StdEncoding.EncodeToString(buStr)
	u.Token = enToken

	fmt.Println("u", u)
	cache_utils.AddOrUpdateCredentialCache(enToken, auth.AuthUserInfo{UserId: u.ID, UserName: emailOrMobile, NamespaceId: u.NamespaceID})
	return &u, nil
}

func (srv *userService) ChangePassword(emailOrMobile string, oldPassword string, newPassword string, namespaceId string) (*models.User, rest_errors.RestErr) {

	//dao.Token = password

	if !global.IsValidEmail(emailOrMobile) && !global.IsValidPhoneNumber(emailOrMobile) {
		return nil, rest_errors.NewUnauthorizedError("Invalid email or mobile")
	}
	if len(newPassword) == 0 {
		return nil, rest_errors.NewUnauthorizedError("Invalid new password")
	}
	var email, mobile string
	if global.IsValidEmail(emailOrMobile) {
		email = emailOrMobile
	} else {
		mobile = emailOrMobile
	}
	us, err := srv.Login(email, mobile, namespaceId)
	if err != nil {
		return nil, err
	}
	if us.Token != crypto_utils.GetMd5(oldPassword) {
		e := rest_errors.NewUnauthorizedError("old Password is not matching")
		//	e.SetErrorCode(rest_errors.OldPasswordIsNotMatching)
		return nil, e
	}
	us.Token = crypto_utils.GetMd5(newPassword)
	us.PasswordUpdatedAt = date_utils.GetNowSearchFormat()
	_, err = srv.Save(*us)
	//err1 := daos.UpdateNewPassword(us.ID, token, date_utils.GetNowSearchFormat())
	if err != nil {
		return nil, err
	}
	return us, nil

}

func (srv *userService) Logout(authToken string) rest_errors.RestErr {

	//Since using crendential as token nothing require but managed from UI by clearing the session

	// if err := daos.UserLogoutUpdate(authToken); err != nil {
	// 	return err
	// }

	return nil
}

func (srv *userService) Delete(id string, h auth.Header) rest_errors.RestErr {
	u, err := srv.Get(id, h)
	if err != nil {
		return err
	}
	u.IsActive = "f"
	_, err = srv.Save(*u)
	if err != nil {
		return err
	}
	return nil
}
