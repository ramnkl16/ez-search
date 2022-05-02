package models

//collection
type EventQueues []EventQueue

//Auto code generated with help of xml schema
// table : EventQueue

type EventQueue struct {
	ID          string `json:"id"`          // id
	EventTypeID string `json:"eventTypeId"` // eventTypeId
	EventData   string `json:"eventData"`   // eventData
	Status      int    `json:"status"`      // status
	StartAt     string `json:"startAt"`     // startAt
	RetryCount  int    `json:"retryCount"`  // retryCount
	Message     string `json:"message"`     // Message
	IsActive    bool   `json:"isActive"`    // isActive
	CreatedAt   string `json:"createdAt"`   // createdAt
	UpdatedAt   string `json:"updatedAt"`   // updatedAt
}
