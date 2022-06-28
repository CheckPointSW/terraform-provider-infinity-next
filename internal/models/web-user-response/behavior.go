package models

type WebUserResponseBehavior struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	MessageTitle     string `json:"messageTitle"`
	MessageBody      string `json:"messageBody"`
	HTTPResponseCode int    `json:"httpResponseCode"`
	RedirectURL      string `json:"redirectURL"`
	XEventID         bool   `json:"xEventId"`
}
