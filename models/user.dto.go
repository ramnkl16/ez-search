package models

//collection
type Users []User

type User struct {
	ID                string `json:"id"`   // WidgetName
	Name              string `json:"name"` // WidgetName
	UserName          string `json:"userName"`
	Email             string `json:"email"`
	Mobile            string `json:"mobile"`
	NamespaceID       string `json:"namespaceId"`
	UserGroupID       string `json:"userGroupID"`
	Token             string `json:"token"`             // token
	FirstName         string `json:"firstName"`         // pageName
	LastName          string `json:"lastName"`          // datasource
	RoleId            string `json:"roleId"`            // datasourceType
	IsActive          string `json:"isActive"`          // isActive
	CreatedAt         string `json:"createdAt"`         // createdAt
	UpdatedAt         string `json:"updatedAt"`         // updatedAt
	EmailVerified     string `json:"emailVerified"`     // email_verified
	PasswordUpdatedAt string `json:"passwordUpdatedAt"` // password_updated_at
}

type Groups []Group
type Group struct {
	ID        string `json:"id"`   // WidgetName
	Name      string `json:"name"` // WidgetName
	Desc      string `json:"desc"`
	IsActive  string `json:"isActive"`  // isActive
	CreatedAt string `json:"createdAt"` // createdAt
	UpdatedAt string `json:"updatedAt"` // updatedAt
}
type Menus []Menu
type Menu struct {
	ID        string `json:"id"`   // WidgetName
	Name      string `json:"name"` // WidgetName
	Link      string `json:"link"`
	ParentId  string `json:"parentId"`
	IsActive  string `json:"isActive"`  // isActive
	CreatedAt string `json:"createdAt"` // createdAt
	UpdatedAt string `json:"updatedAt"` // updatedAt
}
type UserBase64 struct {
	UserName  string `json:"u"`
	Password  string `json:"p"`
	Namespace string `json:"n"`
}
