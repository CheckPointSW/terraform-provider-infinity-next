package models

type CreateWebUserResponseBehaviorInput struct {
	Name             string `json:"name"`
	Visibility       string `json:"visibility"`
	Mode             string `json:"mode"`
	MessageTitle     string `json:"messageTitle"`
	MessageBody      string `json:"messageBody"`
	HTTPResponseCode int    `json:"httpResponseCode"`
	RedirectURL      string `json:"redirectURL"`
	XEventID         bool   `json:"xEventId"`
}
