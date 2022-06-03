package models

//collection
type EventQueues []EventQueue

//Auto code generated with help of xml schema
// table : EventQueue

type EventQueue struct {
	ID                      string `json:"id"`         // id
	EventType               string `json:"eventType"`  // eventTypeId
	EventData               string `json:"eventData"`  // eventData
	Status                  int    `json:"status"`     // status
	StartAt                 string `json:"startAt"`    // startAt
	RetryCount              int    `json:"retryCount"` // retryCount
	RetryMax                int    `json:"retryMax"`
	RetryDuraitionInSeconds int    `json:"retryDuraition"`
	Message                 string `json:"message"`   // Message
	IsActive                string `json:"isActive"`  // isActive
	CreatedAt               string `json:"createdAt"` // createdAt
	UpdatedAt               string `json:"updatedAt"` // updatedAt
	RecurringInSeconds      int    `json:"RecurringInSeconds"`
	LastSyncAt              string `json:"lastSyncAt"`
	TimeTaken               int    `json:"timeTaken"`
}
