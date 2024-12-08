package models

type UpdateWebUserResponseBehaviorInput struct {
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	MessageTitle     string `json:"messageTitle,omitempty"`
	MessageBody      string `json:"messageBody,omitempty"`
	Visibility       string `json:"visibility,omitempty"`
	HTTPResponseCode int    `json:"httpResponseCode,omitempty"`
	RedirectURL      string `json:"redirectURL,omitempty"`
	XEventID         bool   `json:"xEventId"`
}
