package models

type UserGroup struct {
	Base
	Name        string `json:"name"`        // name
	Description string `json:"desc"`        // description
	NamespaceID string `json:"namespaceId"` // namespace_id
	Level       int8   `json:"level"`       // level
}
