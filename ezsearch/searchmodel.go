package ezsearch

import (
	"time"
)

type SearchLogModelClient struct {
	Ref1          string    `json:"ref1"`
	Ref2          string    `json:"ref2"`
	Level         string    `json:"level"`
	Msg           string    `json:"msg"`
	CreatedAt     time.Time `json:"createdAt"`
	EntityType    string    `json:"entityType"`
	ReqUrl        string    `json:"reqUrl"`
	ReqMethod     string    `json:"reqMethod"`
	ReqBody       string    `json:"reqBody"`
	ResStatus     string    `json:"resStatus"`
	BytesReceived int       `json:"bytesReceived"`
	TimeTaken     int       `json:"timeTaken"`
	LangCode      string    `json:"langCode"`
}

type SearchIntModelClient struct {
	ProductSku     string `json:"productSku"`
	LanguageCode   string `json:"languageCode"`
	ExternalIntId  string `json:"externalIntId"`
	ExternlRefId   string `json:"externlRefId"`
	CreatedDate    string `json:"createdDate"`
	ModifiedDate   string `json:"modifiedDate"`
	Classification string `json:"class"`
	Category       string `json:"category"`
}

type SearchEventModelClient struct {
	Id         string `json:"id"`
	EventType  string `json:"eventType"`
	EventData  string `json:"eventData"`
	StartAt    string `json:"startAt"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Status     string `json:"status"`
	RetryCount int    `json:"retryCount"`
	Message    string `json:"message"`
	IsActive   bool   `json:"isActive"`
}

type AggregateProductClient struct {
	Sku                 string `json:"sku"`
	LanguageCode        string `json:"langCode"`
	CatalogId           string `json:"catalogId"`
	Classification      string `json:"class"`
	Category            string `json:"category"`
	Warrantycode        string `json:"warrantycode"`
	Status              string `json:"status"`
	DisplayName         string `json:"displayName"`
	Discontinueddate    string `json:"discontinuedDt"`
	Launcheddate        string `json:"launchedDt"`
	CMUpdatedDate       string `json:"cmUpdatedDt"`
	CDUpdatedDate       string `json:"cdUpdatedDt"`
	OccUpdatedDate      string `json:"occUpdatedDt"`
	ContentModifiedDate string `json:"contentModifiedDt"`
	HasImage            bool   `json:"hasImage"`
	ErrorMsg            string `json:"errMsg1"`
	ErrorMsg1           string `json:"errMsg2"`
	Key                 string `json:"key"`
	HasError            string `json:"hasError"`
}
