package exceptions

import (
	"fmt"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
)

type ExceptionObjectInput struct {
	Match   string   `json:"match"`
	Actions []string `json:"actions"`
	Comment string   `json:"comment,omitempty"`
}

type SchemaExceptionObject struct {
	ID       string            `json:"id,omitempty"`
	Match    map[string]string `json:"match"`
	Action   string            `json:"action"`
	ActionID string            `json:"action_id,omitempty"`
	Comment  string            `json:"comment,omitempty"`
}

type CreateExceptionBehaviorInput struct {
	Name       string                 `json:"name,omitempty"`
	Visibility string                 `json:"visibility,omitempty"`
	Exceptions []ExceptionObjectInput `json:"exceptions,omitempty"`
}

type CreateExceptionBehaviorResponse struct {
	ID string `json:"id"`
}

type Match struct {
	Type     string   `json:"type"`
	Operator string   `json:"op"`
	Items    []Match  `json:"items,omitempty"`
	Key      string   `json:"key,omitempty"`
	Value    []string `json:"value,omitempty"`
}

type Action struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExceptionObjectAction struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

type ExceptionObject struct {
	ID      string                  `json:"id"`
	Match   string                  `json:"match"`
	Actions []ExceptionObjectAction `json:"actions"`
	Comment string                  `json:"comment,omitempty"`
}

type ExceptionBehavior struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Exceptions []ExceptionObject `json:"exceptions"`
}

type SchemaExceptionBehavior struct {
	ID       string         `json:"id"`
	Match    map[string]any `json:"match"`
	Action   string         `json:"action"`
	ActionID string         `json:"action_id,omitempty"`
	Comment  string         `json:"comment,omitempty"`
}

type ExceptionBehaviorUpdateInput struct {
	Name             string                 `json:"name,omitempty"`
	AddExceptions    []ExceptionObjectInput `json:"addExceptions,omitempty"`
	RemoveExceptions []string               `json:"removeExceptions,omitempty"`
}

type UpdateExceptionBehaviorInput struct {
	ID            string                       `json:"id"`
	BehaviorInput ExceptionBehaviorUpdateInput `json:"behaviorInput"`
}

func NewAction(action string) Action {
	return Action{Key: "action", Value: action}
}

func AndMatchFromMap(m map[string]string) Match {
	matches := make([]Match, 0, len(m))
	for key, value := range m {
		match := Match{
			Type:     "condition",
			Operator: "equals",
			Key:      key,
			Value:    []string{value},
		}

		matches = append(matches, match)
	}

	// no need for an "and" operator if only one match
	if len(matches) == 1 {
		return matches[0]
	}

	// else - create an "and" operator between all matches
	var rootMatch Match
	rootMatch.Type = "operator"
	rootMatch.Operator = "and"
	rootMatch.Items = matches
	return rootMatch
}

func (m Match) ToSchemaMap() map[string]interface{} {
	if len(m.Items) == 0 {
		return map[string]interface{}{
			m.Key: m.Value[0],
		}
	}

	result := make(map[string]interface{}, len(m.Items))
	for _, match := range m.Items {
		if len(match.Value) > 0 {
			result[match.Key] = match.Value[0]
		}
	}

	return result
}

func (m Match) String() string {
	return fmt.Sprintf("%#v", m)
}

// ExceptionsDiff returns the diff of old and new
func ExceptionsDiff(old, new []SchemaExceptionBehavior) (added, removed []string) {
	oldIDs := make([]string, 0, len(old))
	oldValues := make([]string, 0, len(old))
	for _, exception := range old {
		oldIDs = append(oldIDs, exception.ID)
		oldValues = append(oldValues, fmt.Sprintf("%#v", exception))
	}

	newIDs := make([]string, 0, len(new))
	newValues := make([]string, 0, len(new))

	for _, exception := range new {
		newIDs = append(newIDs, exception.ID)
		newValues = append(newValues, fmt.Sprintf("%#v", exception))
	}

	added = utils.Added(oldValues, newValues)
	removed = utils.Removed(oldIDs, newIDs)

	return
}

// ExceptionsUpdate returns the exceptions update input for update operation,
// it deletes all the old exceptions and adds the new ones
func ExceptionsUpdate(oldIDs, newValues []string) (added, removed []string) {
	added = newValues
	removed = oldIDs

	return
}
