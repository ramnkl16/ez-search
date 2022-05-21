package models

type Namespace struct {
	ID           string `json:"id"`           // id
	Name         string `json:"name"`         // eventTypeId
	CustomJson   string `json:"customJson"`   // eventData
	Code         string `json:"code"`         // Message
	ContextToken string `json:"contextToken"` // Message
	IsActive     string `json:"isActive"`     // isActive
	CreatedAt    string `json:"createdAt"`    // createdAt
	UpdatedAt    string `json:"updatedAt"`    // updatedAt
}

type NamespaceParam struct {
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	Code         string `json:"code"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	ContextToken string `json:"contextToken"`
	Id           string `json:"id"`
}
