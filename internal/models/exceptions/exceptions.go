package models

import (
	"encoding/json"
	"fmt"
)

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

type ExceptionsObjects []ExceptionObject

// ExceptionBehavior represents an exception behavior object as it is returned from the API
type ExceptionBehavior struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Visibility string            `json:"visibility"`
	Exceptions ExceptionsObjects `json:"exceptions"`
}

func (exceptions ExceptionsObjects) ToSchema() []SchemaExceptionObject {
	var ret []SchemaExceptionObject
	for _, exception := range exceptions {
		ret = append(ret, exception.ToSchema())
	}

	return ret
}

func (exception ExceptionObject) ToSchema() SchemaExceptionObject {
	var ret SchemaExceptionObject

	ret.ID = exception.ID
	ret.Comment = exception.Comment
	var action Action
	if err := json.Unmarshal([]byte(exception.Actions[0].Action), &action); err != nil {
		fmt.Printf("failed to unmarshal action string to struct: %+v", err)
	}

	ret.Action = action.Value
	ret.ActionID = exception.Actions[0].ID
	var match Match
	if err := json.Unmarshal([]byte(exception.Match), &match); err != nil {
		fmt.Printf("failed to unmarshal match string to struct: %+v", err)
	}

	ret.Match = []SchemaMatchExpression{MatchToSchema(match)}

	return ret
}

func MatchToSchema(match Match) SchemaMatchExpression {
	var ret SchemaMatchExpression

	ret.Operator = match.Operator
	if len(match.Items) == 0 {
		ret.Key = match.Key
		ret.Value = match.Value
		return ret
	}

	for _, item := range match.Items {
		ret.Operands = append(ret.Operands, MatchToSchema(item))
	}

	return ret
}
