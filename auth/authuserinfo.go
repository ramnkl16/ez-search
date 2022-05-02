package auth
type AuthUserInfo struct {
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	NamespaceId string
	UserRoleID  string
}
