package models

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

const (
	URLIDSeperator = ";;;;"
)

type PracticeMode struct {
	Mode        string `json:"mode"`
	SubPractice string `json:"subPractice,omitempty"`
}

type Trigger struct {
	ID string `json:"id"`
}

type Behavior struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Behaviors []Behavior

type Practice struct {
	ID string `json:"id"`
}

// PracticeWrapper represents a PracticeWrapper object returned from mgmt
type PracticeWrapper struct {
	PracticeWrapperID string         `json:"id"`
	MainMode          string         `json:"mainMode,omitempty"`
	SubPracticeModes  []PracticeMode `json:"subPracticeModes,omitempty"`
	Triggers          []Trigger      `json:"triggers,omitempty"`
	Behaviors         []Behavior     `json:"behaviors,omitempty"`
	Practice          Practice       `json:"practice"`
}

type PracticesWrappers []PracticeWrapper

// SourceIdentifierValue represents a SourceIdentifierValue object returned from mgmt
// IdentifierValue field is defined by the user
// id field is a unique uuid generated by mgmt
type SourceIdentifierValue struct {
	ID              string `json:"id"`
	IdentifierValue string `json:"IdentifierValue"`
}

// SourceIdentifier represents a SourceIdentifier object returned from mgmt
// sourceIdentifier field is defined by the user
// id field is a unique uuid generated by mgmt
type SourceIdentifier struct {
	ID               string                  `json:"id"`
	SourceIdentifier string                  `json:"sourceIdentifier"`
	Values           []SourceIdentifierValue `json:"values"`
}

type SourceIdentifiers []SourceIdentifier

// ProxySetting represents a ProxySetting object returned from mgmt
// key and value fields are defined by the user
// id field is a unique uuid generated by mgmt
type ProxySetting struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProxySettings []ProxySetting

// URL represents an URL object returned from mgmt
// url field is defined by the user
// id field is a unique uuid generated by mgmt
type URL struct {
	ID  string `json:"id"`
	URL string `json:"URL"`
}

type URLs []URL

// Profile represents a profileId associated with the web API asset
type Profile struct {
	ID string `json:"id"`
}

type Profiles []Profile

// WebAPIAsset represents the response from mgmt after creating the asset
type WebAPIAsset struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	AssetType         string            `json:"assetType"`
	Class             string            `json:"class"`
	Category          string            `json:"category"`
	Family            string            `json:"family"`
	Group             string            `json:"group"`
	Order             string            `json:"order"`
	Kind              string            `json:"kind"`
	MainAttributes    string            `json:"mainAttributes"`
	IntelligenceTags  string            `json:"intelligenceTags"`
	State             string            `json:"state,omitempty"`
	UpstreamURL       string            `json:"upstreamURL,omitempty"`
	Sources           string            `json:"sources"`
	URLs              URLs              `json:"URLs"`
	ProxySettings     ProxySettings     `json:"proxySetting"`
	SourceIdentifiers SourceIdentifiers `json:"sourceIdentifiers"`
	Behaviors         Behaviors         `json:"behaviors,omitempty"`
	Profiles          Profiles          `json:"profiles,omitempty"`
	Practices         PracticesWrappers `json:"practices,omitempty"`
	ReadOnly          bool              `json:"readOnly"`
}

// ToSchema returns a slice of profiles IDs to be saved in the state file
func (profiles Profiles) ToSchema() []string {
	mapFunc := func(profile Profile) string {
		return profile.ID
	}

	return utils.Map(profiles, mapFunc)
}

// ToSchema returns a slice of behaviors IDs to be saved in the state file
func (behaviors Behaviors) ToSchema() []string {
	mapFunc := func(behavior Behavior) string {
		return behavior.ID
	}

	return utils.Map(behaviors, mapFunc)
}

// ToSchema converts the URLs as returned from the APi to 2 slices of strings to be saved in the state file:
// 1. IDs slice
// 2. URLs slice
func (urls URLs) ToSchema() ([]string, []string) {
	schemaURLs := make([]string, len(urls))
	schemaURLsIDs := make([]string, len(urls))
	for i, url := range urls {
		schemaURLs[i] = url.URL
		schemaURLsIDs[i] = url.createSchemaID()
	}

	return schemaURLs, schemaURLsIDs
}

// createSchemaID excepts a URL object and returns the following string:
// "<url><seperator><url_id>"
// This string is saved in the state file as an url ID
func (url *URL) createSchemaID() string {
	return fmt.Sprintf("%s%s%s", url.URL, URLIDSeperator, url.ID)
}

// ToSchema converts the SourceIdentifiers field as returned from the API to a slice of
// SchemaSourceIdentifier to be saved in the state file
func (sourceIdentifiers SourceIdentifiers) ToSchema() []SchemaSourceIdentifier {
	mapFunc := func(source SourceIdentifier) SchemaSourceIdentifier {
		return source.ToSchema()
	}

	return utils.Map(sourceIdentifiers, mapFunc)
}

// toSchema converts a single SourceIdentifier as returned from the API to a single
// SchemaSourceIdentifier to be saved in the state file
func (sourceIdentifier SourceIdentifier) ToSchema() SchemaSourceIdentifier {
	values := make([]string, len(sourceIdentifier.Values))
	valuesIDs := make([]string, len(sourceIdentifier.Values))
	for j, value := range sourceIdentifier.Values {
		values[j] = value.IdentifierValue
		valuesIDs[j] = fmt.Sprintf("%s%s%s", value.IdentifierValue, SourceIdentifierValueIDSeparator, value.ID)
	}

	return SchemaSourceIdentifier{
		ID:               sourceIdentifier.ID,
		SourceIdentifier: sourceIdentifier.SourceIdentifier,
		Values:           values,
		ValuesIDs:        valuesIDs,
	}
}

// ToSchema converts the paractices field as returned from the API to a slice of
// SchemaPracticeWrapper to be saved in the state file
func (wrappers PracticesWrappers) ToSchema() []SchemaPracticeWrapper {
	mapFunc := func(wrapper PracticeWrapper) SchemaPracticeWrapper {
		return wrapper.toSchema()
	}

	return utils.Map(wrappers, mapFunc)
}

// toSchema converts a single PracticeWrapper as returned from the API to a single
// SchemaPracticeWrapper to be saved in the state file
func (practiceWrapper PracticeWrapper) toSchema() SchemaPracticeWrapper {
	triggers := make([]string, len(practiceWrapper.Triggers))
	for j, trigger := range practiceWrapper.Triggers {
		triggers[j] = trigger.ID
	}

	behaviors := make([]string, len(practiceWrapper.Behaviors))
	for j, behavior := range practiceWrapper.Behaviors {
		behaviors[j] = behavior.ID
	}

	subPracticeModes := make(map[string]string)
	for _, mode := range practiceWrapper.SubPracticeModes {
		subPracticeModes[mode.SubPractice] = mode.Mode
	}

	return SchemaPracticeWrapper{
		PracticeWrapperID: practiceWrapper.PracticeWrapperID,
		PracticeID:        practiceWrapper.Practice.ID,
		MainMode:          practiceWrapper.MainMode,
		SubPracticeModes:  subPracticeModes,
		Triggers:          triggers,
		Behaviors:         behaviors,
	}
}
