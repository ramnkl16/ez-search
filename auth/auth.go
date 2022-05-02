package auth

import (
	"net/http"

	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
)

const (
	headerXPublic = "x-public"
	headerXauth   = "x-auth"
	headerXDebug  = "x-debug"
	headerNS      = "x-ns"
)

func GetNamespace(request *http.Request) string {
	if request == nil {
		return ""
	}
	return request.Header.Get(headerNS)
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func IsDebug(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXDebug) == "true"
}

func GetXauthToken(request *http.Request) string {
	if request == nil {
		return ""
	}
	return request.Header.Get(headerXauth)
}

func GetAuthInfo(authToken string) (AuthUserInfo, rest_errors.RestErr) {
	var userInfo AuthUserInfo
	userObj, b := cache_utils.Cache.Get(authToken)
	if b != nil {
		return userInfo, rest_errors.NewRestError("", http.StatusNotFound, "token already expired or not exist", nil)
	}
	userInfo = userObj.(AuthUserInfo)

	if len(userInfo.UserId) == 0 {
		return userInfo, rest_errors.NewRestError("", http.StatusNotFound, "User info is empty from cache", nil)
	}
	return userInfo, nil
}
func AuthenticateRequest(request *http.Request) (AuthUserInfo, rest_errors.RestErr) {
	var userInfo AuthUserInfo
	if request == nil {
		return userInfo, nil
	}
	userObj, b := cache_utils.Cache.Get(GetXauthToken(request))
	if b != nil {
		return userInfo, rest_errors.NewRestError("", http.StatusNotFound, "token already expired or not exist", nil)
	}
	userInfo = userObj.(AuthUserInfo)

	if len(userInfo.UserId) == 0 {
		return userInfo, rest_errors.NewRestError("", http.StatusNotFound, "User info is empty from cache", nil)
	}
	return userInfo, nil
}
