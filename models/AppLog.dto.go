package models
	

	 
	//collection
type  AppLogs []AppLog

	//Auto code generated with help of xml schema 
	// table : AppLog
	
type AppLog struct { 
		ID string `json:"id"` // id 
		Ref1 string `json:"ref1"` // ref1 
		Ref2 string `json:"ref2"` // ref2 
		Level string `json:"level"` // level 
		Message string `json:"message"` // message 
		CreatedAt string `json:"createdAt"` // createdAt 
		UpdatedAt string `json:"updatedAt"` // updatedAt	
}
