package models

type WebUserResponseBehavior struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	MessageTitle     string `json:"messageTitle"`
	MessageBody      string `json:"messageBody"`
	Visibility       string `json:"visibility"`
	HTTPResponseCode int    `json:"httpResponseCode"`
	RedirectURL      string `json:"redirectURL"`
	XEventID         bool   `json:"xEventId"`
}
