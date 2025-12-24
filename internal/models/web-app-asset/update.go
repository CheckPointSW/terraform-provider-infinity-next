package models

// UpdateSourceIdentifierValue represents the input for updating an existing source identifier value
// in a source identifier field of an existing WebApplicationAsset object
type UpdateSourceIdentifierValue struct {
	ID              string `json:"id"`
	IdentifierValue string `json:"identifierValue"`
}

// UpdateSourceIdentifier represents the input for updating an existing proxy setting field of an existing WebApplicationAsset object
type UpdateSourceIdentifier struct {
	ID               string                        `json:"id"`
	SourceIdentifier string                        `json:"sourceIdentifier"`
	AddValues        []string                      `json:"addValues"`
	RemoveValues     []string                      `json:"removeValues"`
	UpdateValues     []UpdateSourceIdentifierValue `json:"updateValues"`
}

// AddSourceIdentifier represents the input for adding a source identifier field to an existing WebApplicationAsset object
type AddSourceIdentifier struct {
	SourceIdentifier string   `json:"sourceIdentifier"`
	Values           []string `json:"values"`
}

// UpdateURL represents the input for updating an existing url field of an existing WebApplicationAsset object
type UpdateURL struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// UpdateProxySetting represents the input for updating an existing proxy setting field of an existing WebApplicationAsset object
type UpdateProxySetting struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AddProxySetting represents the input for adding a proxy setting field to an existing WebApplicationAsset object
type AddProxySetting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AddPracticeMode represents the api input for adding a practice mode (sub practice) in a practice to add
// to an existing WebApplicationAsset object
type AddPracticeMode struct {
	Mode        string `json:"mode"`
	SubPractice string `json:"subPractice,omitempty"`
}

// AddPracticeWrapper represents the input for adding a practice wrapper field to an existing WebApplicationAsset object
type AddPracticeWrapper struct {
	PracticeID       string            `json:"practiceId"`
	MainMode         string            `json:"mainMode"`
	SubPracticeModes []AddPracticeMode `json:"subPracticeModes,omitempty"`
	Triggers         []string          `json:"triggers"`
}

type AddTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddTags []AddTag

// UpdateWebApplicationAssetInput represents the input for updating an existing WebApplicationAsset object
type UpdateWebApplicationAssetInput struct {
	Name                    string                   `json:"name,omitempty"`
	AddPracticeWrappers     []AddPracticeWrapper     `json:"addPractices,omitempty"`
	RemovePracticeWrappers  []string                 `json:"removePractices,omitempty"`
	AddProfiles             []string                 `json:"addProfiles,omitempty"`
	RemoveProfiles          []string                 `json:"removeProfiles,omitempty"`
	AddBehaviors            []string                 `json:"addBehaviors,omitempty"`
	RemoveBehaviors         []string                 `json:"removeBehaviors,omitempty"`
	AddTags                 AddTags                  `json:"addTags,omitempty"`
	RemoveTags              []string                 `json:"removeTags,omitempty"`
	State                   string                   `json:"state,omitempty"`
	AddProxySetting         []AddProxySetting        `json:"addProxySetting,omitempty"`
	RemoveProxySetting      []string                 `json:"removeProxySetting,omitempty"`
	UpdateProxySetting      []UpdateProxySetting     `json:"updateProxySetting,omitempty"`
	UpstreamURL             string                   `json:"upstreamURL,omitempty"`
	AddURLs                 []string                 `json:"addURLs,omitempty"`
	RemoveURLs              []string                 `json:"removeURLs,omitempty"`
	UpdateURLs              []UpdateURL              `json:"updateURLs,omitempty"`
	AddSourceIdentifiers    []AddSourceIdentifier    `json:"addSourceIdentifiers,omitempty"`
	RemoveSourceIdentifiers []string                 `json:"removeSourceIdentifiers,omitempty"`
	UpdateSourceIdentifiers []UpdateSourceIdentifier `json:"updateSourceIdentifiers,omitempty"`
	IsSharesURLs            *bool                    `json:"isSharesURLs,omitempty"`
}
