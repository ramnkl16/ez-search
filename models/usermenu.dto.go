package models

import "github.com/ramnkl16/ez-search/global"

//collection
type UserMenus []UserMenu

//Auto code generated with help of mysql table schema
// table : UserMenu

type UserMenu struct {
	ID                string               `json:"id"`
	NamespaceID       string               `json:"namespaceId"`       // namespace_id
	MenuID            string               `json:"menuId"`            // menu_id
	ReferenceID       string               `json:"refId"`             // reference_id
	ReferenceType     global.ReferenceType `json:"refType"`           // reference_type
	Privilege         int8                 `json:"privilege"`         // permission
	CustomData        string               `json:"cd"`                // custom data
	PermissionPlus    int8                 `json:"permissionPlus"`    // permission_plus
	MenuExceptionFlag int8                 `json:"menuExceptionFlag"` // menu_exception_flag
	IsActive          bool                 `json:"isActive"`          // active_flag
	UpdatedBy         string               `json:"updatedBy"`         // updated_by
	UpdatedAt         string               `json:"updatedAt"`         // updated_at
	CreatedAt         string               `json:"createdAt"`
}

type PermissionParam struct {
	MenuId         string `json:"menuId"`
	Permission     int8   `json:"permission"`
	ReferenceType  string `json:"referenceType"`
	ReferenceID    string `json:"referenceId"`
	RemoveOverride int8   `json:"removeOverride"`
	Namespace      string `json:"namespace"`
}
