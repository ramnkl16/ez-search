package models

//collection
type WidgetMetas []WidgetMeta

//Auto code generated with help of xml schema
// table : WidgetMeta

type WidgetMeta struct {
	ID        string `json:"id"`   // WidgetName
	Name      string `json:"name"` // WidgetName
	Division  string `json:"division"`
	Module    string `json:"module"`
	Page      string `json:"page"`      // pageName
	Data      string `json:"cd"`        // datasource
	IsActive  string `json:"isActive"`  // isActive
	CreatedAt string `json:"createdAt"` // createdAt
	UpdatedAt string `json:"updatedAt"` // updatedAt
}

// {Name: "Name", Type: "text"},
// 		{Name: "division", Type: "text"},
// 		{Name: "module", Type: "text"},
// 		{Name: "page", Type: "text"},
// 		{Name: "data", Type: "text"},
// 		{Name: "isActive", Type: "bool"},
// 		{Name: "createdAt", Type: "date"},
// 		{Name: "updatedAt", Type: "date"}}
