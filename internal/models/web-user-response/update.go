package models

type UpdateWebUserResponseBehaviorInput struct {
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	MessageTitle     string `json:"messageTitle,omitempty"`
	MessageBody      string `json:"messageBody,omitempty"`
	Visibility       string `json:"visibility,omitempty"`
	HTTPResponseCode int    `json:"httpResponseCode"`
	RedirectURL      string `json:"redirectURL,omitempty"`
	XEventID         bool   `json:"xEventId"`
}

type DisplayObject struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	SubType      string `json:"subType,omitempty"`
	ObjectStatus string `json:"objectStatus,omitempty"`
}

type DisplayObjects []DisplayObject
