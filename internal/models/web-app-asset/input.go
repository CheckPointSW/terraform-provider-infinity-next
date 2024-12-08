package models

import (
	"strings"
)

// PracticeModeInput represents the api input for creating a practice mode field
// in the practice field of the web application asset
type PracticeModeInput struct {
	Mode        string `json:"mode"`
	SubPractice string `json:"subPractice,omitempty"`
}

// TagInput represents the api input for creating a tag field in the web application asset
type TagInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	ID    string `json:"id,omitempty"`
}

type TagsInputs []TagInput

// PracticeWrapperInput represents the api input for creating a practice field in the web application asset
type PracticeWrapperInput struct {
	PracticeWrapperID string              `json:"practiceWrapperId,omitempty"`
	PracticeID        string              `json:"practiceId"`
	MainMode          string              `json:"mainMode"`
	SubPracticeModes  []PracticeModeInput `json:"subPracticeModes,omitempty"`
	Triggers          []string            `json:"triggers,omitempty"`
}

type PracticeWrappersInputs []PracticeWrapperInput

// ProxySettingInput represents the api input for creating a proxy setting field in the web application asset
type ProxySettingInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	ID    string `json:"id,omitempty"`
}

type ProxySettingInputs []ProxySettingInput

// SourceIdentifierInput represents the api input for creating a source identifier field in the web application asset
type SourceIdentifierInput struct {
	SourceIdentifier string    `json:"sourceIdentifier"`
	Values           []string  `json:"values"`
	ValuesIDs        ValuesIDs `json:"valuesIDs,omitempty"`
	ID               string    `json:"id,omitempty"`
}

type ValuesIDs []string

type SourceIdentifiersInputs []SourceIdentifierInput

// CreateWebApplicationAssetInput represents the api input for creating a web application asset
type CreateWebApplicationAssetInput struct {
	Name              string                  `json:"name"`
	PracticeWrappers  PracticeWrappersInputs  `json:"practices,omitempty"`
	Profiles          []string                `json:"profiles,omitempty"`
	Behaviors         []string                `json:"behaviors,omitempty"`
	Tags              TagsInputs              `json:"tags,omitempty"`
	ProxySettings     ProxySettingInputs      `json:"proxySetting,omitempty"`
	UpstreamURL       string                  `json:"upstreamURL,omitempty"`
	URLs              []string                `json:"URLs,omitempty"`
	SourceIdentifiers SourceIdentifiersInputs `json:"sourceIdentifiers,omitempty"`
	IsSharesURLs      bool                    `json:"isSharesURLs,omitempty"`
}

// ToIndicatorsMap converts a ProxySettingInputs to a map from a proxy setting key to the proxy setting struct itself
func (inputs ProxySettingInputs) ToIndicatorsMap() map[string]ProxySettingInput {
	ret := make(map[string]ProxySettingInput)
	for _, input := range inputs {
		ret[input.Key] = input
	}

	return ret
}

// ToIndicatorsMap converts a SourceIdentifiersInputs to a map from a source identifier field to the source identifier struct itself
func (inputs SourceIdentifiersInputs) ToIndicatorsMap() map[string]SourceIdentifierInput {
	ret := make(map[string]SourceIdentifierInput)
	for _, input := range inputs {
		ret[input.SourceIdentifier] = input
	}

	return ret
}

func (ids ValuesIDs) ToIndicatorsMap() map[string]string {
	ret := make(map[string]string)
	for _, sourceIdentifierValueID := range ids {
		valueAndID := strings.Split(sourceIdentifierValueID, SourceIdentifierValueIDSeparator)
		ret[valueAndID[0]] = valueAndID[1]
	}

	return ret
}

func (inputs TagsInputs) ToIndicatorsMap() map[string]TagInput {
	ret := make(map[string]TagInput)
	for _, input := range inputs {
		ret[input.Key] = input
	}

	return ret
}

func (mtlsInputs FileSchemas) ToIndicatorMap() map[string]FileSchema {
	mTLSs := make(map[string]FileSchema)
	for _, mTLS := range mtlsInputs {
		mTLSs[mTLS.Type] = mTLS
	}

	return mTLSs
}
