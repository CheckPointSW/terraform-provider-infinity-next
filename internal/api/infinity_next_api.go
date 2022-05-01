package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphErrorExtension struct {
	Code          string `json:"code,omitempty"`
	ReferenceID   string `json:"referenceId,omitempty"`
	MessageParams any    `json:"messageParams,omitempty"`
}

type GraphError struct {
	Message    string               `json:"message,omitempty"`
	Path       []string             `json:"path,omitempty"`
	Extensions *GraphErrorExtension `json:"extensions,omitempty"`
}

type GraphResponse struct {
	Data   any          `json:"data"`
	Errors []GraphError `json:"errors,omitempty"`
}

type RetryErrorResponse struct {
	Stop          any `json:"stop"`
	RC            any `json:"rc"`
	ReqDidTimeout any `json:"reqDidTimeout"`
}

func (c *Client) InfinityPortalAuthentication(clientId string, secretKey string) error {
	timeout := clientTimeout * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	formData := url.Values{
		"clientId":  {clientId},
		"accessKey": {secretKey},
	}

	resp, err := client.PostForm(c.host+"/auth/external", formData)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	data, err := json.Marshal(result["data"])
	if err != nil {
		return err
	}

	jsonStr := string(data)
	var datamap map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &datamap); err != nil {
		return err
	}

	tokenInterface, ok := datamap["token"]
	if !ok {
		return fmt.Errorf("missing token in response %#v", result)
	}

	c.token = tokenInterface.(string)

	return nil
}

func (c *Client) MakeGraphQLRequest(gql, responseKey string, vars ...map[string]any) (any, error) {
	variables := make(map[string]any)
	for _, varMap := range vars {
		for k, v := range varMap {
			variables[k] = v
		}
	}

	graphQlRequest := GraphQLRequest{
		Query:     gql,
		Variables: variables,
	}

	graphQlRequestBytes, err := json.Marshal(graphQlRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request. Error: %+v", err)
	}

	timeout := clientTimeout * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	var res *http.Response
	for i := 0; i < rateLimitingNumOfRetries; i++ {
		httpRequest, err := http.NewRequest(http.MethodPost, c.host+c.endpoint, bytes.NewReader(graphQlRequestBytes))
		if err != nil {
			return nil, err
		}

		bearer := "Bearer " + c.token
		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Authorization", bearer)

		res, err = client.Do(httpRequest)
		if err != nil {
			return nil, err
		}

		switch res.StatusCode {
		case http.StatusOK:
			// breaks out of switch statemant
		case http.StatusTooManyRequests:
			res.Body.Close()
			fmt.Println("[WARN] GraphQL request failed due to rate limmiting, retrying..")
			time.Sleep(time.Second * 2)
			continue
		default:
			defer res.Body.Close()
			graphResponse, err := parseGraphQLResponse(res)
			if err != nil {
				return nil, err
			}

			if len(graphResponse.Errors) > 0 {
				return nil, fmt.Errorf("GraphQL response contains errors: %s, ReferenceID: %s. Body: %+v", graphResponse.Errors[0].Message, graphResponse.Errors[0].Extensions.ReferenceID, graphResponse.Data)
			}

			return nil, fmt.Errorf("Non-OK http code (%d) - Body: %+v", res.StatusCode, graphResponse.Data)
		}

		// if we are here then status code == 200 ok
		// break out of loop
		break
	}

	defer res.Body.Close()
	graphResponse, err := parseGraphQLResponse(res)
	if err != nil {
		return nil, err
	}

	if len(graphResponse.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL response contains errors: %s. ReferenceID: %+v. Body: %+v", graphResponse.Errors[0].Message, graphResponse.Errors[0].Extensions.ReferenceID, graphResponse.Data)
	}

	graphResponseMap, ok := graphResponse.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response, should be of type map[string]any but got %#v", graphResponse.Data)
	}

	ret, ok := graphResponseMap[responseKey]
	if !ok {
		return nil, fmt.Errorf("invalid response field: %s. Full response: %+v", responseKey, graphResponseMap)
	}

	if ret == nil {
		return nil, fmt.Errorf("%s - ReferenceID: %s", ErrorNotFound.Error(), getReferenceIDFromHeaders(res.Header))
	}

	return ret, nil
}

func getReferenceIDFromHeaders(headers http.Header) string {
	for k, v := range headers {
		if k == "Logger-Token" {
			return v[0]
		}
	}

	return "not found in headers"
}

func parseGraphQLResponse(response *http.Response) (*GraphResponse, error) {
	var graphResponse GraphResponse
	graphResponsePointer := &graphResponse
	referenceID := getReferenceIDFromHeaders(response.Header)
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, response.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %+v. ReferenceID: %s", response.Body, referenceID)
	}

	if err := json.NewDecoder(&buf).Decode(&graphResponsePointer); err != nil {
		return nil, fmt.Errorf("faied to decode response body to struct. Body: %+v. ReferenceID: %s", response.Body, referenceID)
	}

	return graphResponsePointer, nil
}

func (c *Client) PublishChanges() (bool, error) {
	return true, nil
}

func (c *Client) DiscardChanges() (bool, error) {
	discardChanges, err := c.MakeGraphQLRequest(`
	mutation discardChanges{
		discardChanges
	}`, "discardChanges")

	if err != nil {
		return false, fmt.Errorf("failed discarding changes: %w", err)
	}

	isDiscarded, ok := discardChanges.(bool)
	if !ok {
		return false, fmt.Errorf("failed discarding changes: got invalid response %#v", discardChanges)
	}

	return isDiscarded, err
}
