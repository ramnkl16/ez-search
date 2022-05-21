package models

type Status int8

const (
	StatusDisable Status = iota + 1
	StatusActive
	StatusDeleted
)

// type BaseInterface interface {
// 	CreateOrUpdate() rest_errors.RestErr
// 	Delete(id string) rest_errors.RestErr
// 	GetAllFrom(query string)
// 	Get(id string)
// }

type Base struct {
	ID        string `json:"id"`        // code
	IsActive  string `json:"isActive"`  // isActive
	UpdatedBy string `json:"updatedBy"` // updated_by
	UpdatedAt string `json:"updatedAt"` // updated_at
	CreatedAt string `json:"createdAt"` // createdAt
}
