package models

//collection
type EventQueueHistories []EventQueueHistory

//Auto code generated with help of xml schema
// table : EventQueueHistory

type EventQueueHistory struct {
	ID                      string `json:"id"`           // id
	EventQueueID            string `json:"eventQueueId"` // eventQueueId
	EventType               string `json:"eventType"`    // eventTypeId
	EventData               string `json:"eventData"`    // eventData
	Status                  int    `json:"status"`       // status
	StartAt                 string `json:"startAt"`      // startAt
	RetryCount              int    `json:"retryCount"`   // retryCount
	RetryDuraitionInSeconds int    `json:"retryDuraition"`
	RetryMax                int    `json:"retryMax"`
	Message                 string `json:"message"`   // Message
	IsActive                string `json:"isActive"`  // isActive
	CreatedAt               string `json:"createdAt"` // createdAt
	UpdatedAt               string `json:"updatedAt"` // updatedAt
	RecurringInSeconds      int    `json:"RecurringInSeconds"`
	LastSyncAt              string `json:"lastSyncAt"`
	TimeTaken               int    `json:"timeTaken"`
}
