package models
	

	 
	//collection
type  EventQueueHistories []EventQueueHistory

	//Auto code generated with help of xml schema 
	// table : EventQueueHistory
	
type EventQueueHistory struct { 
		ID string `json:"id"` // id 
		EventQueueID string `json:"eventQueueId"` // eventQueueId 
		EventTypeID string `json:"eventTypeId"` // eventTypeId 
		EventData string `json:"eventData"` // eventData 
		Status int `json:"status"` // status 
		RetryCount int `json:"retryCount"` // retryCount 
		Message string `json:"Message"` // Message 
		IsActive bool `json:"isActive"` // isActive 
		CreatedAt string `json:"createdAt"` // createdAt 
		UpdatedAt string `json:"updatedAt"` // updatedAt	
}
