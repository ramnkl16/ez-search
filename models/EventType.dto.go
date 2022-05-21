package models

//collection
type EventTypes []EventType

//Auto code generated with help of xml schema
// table : EventType

type EventType struct {
	ID        string `json:"id"`        // id
	Name      string `json:"name"`      // name
	Hint      string `json:"Hint"`      // Hint
	IsActive  string `json:"isActive"`  // isActive
	CreatedAt string `json:"createdAt"` // createdAt
	UpdatedAt string `json:"updatedAt"` // updatedAt
}
