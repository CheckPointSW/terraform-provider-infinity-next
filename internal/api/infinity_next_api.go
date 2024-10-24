package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

const (
	appIDClaim  = "appId"
	wafAppID    = "64488de9-f813-42a7-93e7-f3fe25dd9011"
	policyAppID = "f47b536c-a990-42fb-9ab2-ec38f8c2dcff"
	wafPath     = "/app/waf/graphql/V1"
	policyPath  = "/app/i2/graphql/V1"
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

func (c *Client) InfinityPortalAuthentication(clientId string, accessKey string) error {
	timeout := clientTimeout * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	formData := url.Values{
		"clientId":  {clientId},
		"accessKey": {accessKey},
	}

	for retryCount := 1; retryCount <= maxNumOfRetries; retryCount++ {
		resp, err := client.PostForm(c.host+"/auth/external", formData)
		if err != nil {
			if retryCount == maxNumOfRetries {
				return err
			}

			time.Sleep(2 * time.Second * time.Duration(retryCount))
			continue
		}

		defer resp.Body.Close()

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			if retryCount == maxNumOfRetries {
				return err
			}

			time.Sleep(2 * time.Second * time.Duration(retryCount))
			continue
		}

		data, err := json.Marshal(result["data"])
		if err != nil {
			if retryCount == maxNumOfRetries {
				return err
			}

			time.Sleep(2 * time.Second * time.Duration(retryCount))
			continue
		}

		jsonStr := string(data)
		var datamap map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &datamap); err != nil {
			if retryCount == maxNumOfRetries {
				return err
			}

			time.Sleep(2 * time.Second * time.Duration(retryCount))
			continue
		}

		tokenInterface, ok := datamap["token"]
		if !ok {
			if retryCount == maxNumOfRetries {
				return fmt.Errorf("missing token in response %#v", result)
			}

			time.Sleep(2 * time.Second * time.Duration(retryCount))
			continue
		}

		c.token = tokenInterface.(string)
		token, _, err := jwt.NewParser().ParseUnverified(c.token, jwt.MapClaims{})
		if err != nil {
			return fmt.Errorf("failed to parse token: %w", err)
		}

		tokenMapClaims := token.Claims.(jwt.MapClaims)
		if appID, ok := tokenMapClaims[appIDClaim]; ok {
			switch appID.(string) {
			case wafAppID:
				c.SetEndpoint(wafPath)
			case policyAppID:
				c.SetEndpoint(policyPath)
			}
		}

		break

	}

	return nil
}

func (c *Client) MakeGraphQLRequest(ctx context.Context, gql, responseKey string, vars ...map[string]any) (any, error) {
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
	for retryCount := 1; retryCount <= maxNumOfRetries; retryCount++ {
		httpRequest, err := http.NewRequest(http.MethodPost, c.host+c.endpoint, bytes.NewReader(graphQlRequestBytes))
		if err != nil {
			return nil, err
		}

		bearer := "Bearer " + c.token
		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Authorization", bearer)

		res, err = client.Do(httpRequest)
		if err != nil {
			if retryCount == maxNumOfRetries {
				return nil, err
			}

			res.Body.Close()
			fmt.Println("[WARN] GraphQL request failed with error " + err.Error() + ", retrying...")
			time.Sleep(time.Second * 2 * time.Duration(retryCount))
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			defer res.Body.Close()
			graphResponse, err := parseGraphQLResponse(res)
			if err != nil {
				return nil, err
			}

			if len(graphResponse.Errors) > 0 {
				if retryCount == maxNumOfRetries {
					return nil, fmt.Errorf("GraphQL response contains errors: %s, ReferenceID: %s. Body: %+v", graphResponse.Errors[0].Message, graphResponse.Errors[0].Extensions.ReferenceID, graphResponse.Data)
				}

				res.Body.Close()
				fmt.Printf("[WARN] GraphQL request failed with status %s and errors %+v, retrying...\n", res.Status, graphResponse.Errors)
				time.Sleep(time.Second * 2 * time.Duration(retryCount))
				continue
			}

			graphResponseMap, ok := graphResponse.Data.(map[string]any)
			if !ok {
				err := fmt.Errorf("invalid response, should be of type map[string]any but got %#v", graphResponse.Data)
				if retryCount == maxNumOfRetries {
					return nil, err
				}

				res.Body.Close()
				fmt.Printf("[WARN] GraphQL request failed with error %v, retrying...\n", err)
				time.Sleep(time.Second * 2 * time.Duration(retryCount))
				continue
			}

			ret, ok := graphResponseMap[responseKey]
			if !ok {
				err := fmt.Errorf("invalid response field: %s. Full response: %+v", responseKey, graphResponseMap)
				if retryCount == maxNumOfRetries {
					return nil, err
				}

				res.Body.Close()
				fmt.Printf("[WARN] GraphQL request failed with error %v, retrying...\n", err)
				time.Sleep(time.Second * 2 * time.Duration(retryCount))
				continue

			}

			if ret == nil {
				// We need to retry only if it's expected to find the resource
				// This is only used for test, because we ensure a resource is destroyed after a test using Read.
				if v := ctx.Value(utils.ExpectResourceNotFound); v != nil && !v.(bool) {
					err := fmt.Errorf("%s - ReferenceID: %s", ErrorNotFound.Error(), getReferenceIDFromHeaders(res.Header))
					if retryCount == maxNumOfRetries {
						return nil, err
					}

					res.Body.Close()
					fmt.Printf("[WARN] GraphQL request failed with error %v, retrying...\n", err)
					time.Sleep(time.Second * 2 * time.Duration(retryCount))
					continue
				}
			}

			return ret, nil
		case http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway, http.StatusRequestTimeout:
			res.Body.Close()
			fmt.Println("[WARN] GraphQL request failed with status " + res.Status + ", retrying...")
			time.Sleep(time.Second * 2 * time.Duration(retryCount))
			continue
		default:
			defer res.Body.Close()
			graphResponse, err := parseGraphQLResponse(res)
			if err != nil {
				return nil, err
			}

			if len(graphResponse.Errors) > 0 {
				if retryCount == maxNumOfRetries {
					return nil, fmt.Errorf("GraphQL response contains errors: %s, ReferenceID: %s. Body: %+v", graphResponse.Errors[0].Message, graphResponse.Errors[0].Extensions.ReferenceID, graphResponse.Data)
				}

				res.Body.Close()
				fmt.Printf("[WARN] GraphQL request failed with status %s and errors %+v, retrying...\n", res.Status, graphResponse.Errors)
				time.Sleep(time.Second * 2 * time.Duration(retryCount))
				continue
			}

			return nil, fmt.Errorf("Non-OK http code (%d) - Body: %+v", res.StatusCode, graphResponse.Data)
		}
	}

	return nil, nil
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
		return nil, fmt.Errorf("failed to read response body: %+v. ReferenceID: %s: %w", response.Body, referenceID, err)
	}

	bResponse := buf.Bytes()
	if err := json.Unmarshal(bResponse, &graphResponsePointer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body. Body: %s. ReferenceID: %s: %w", string(bResponse), referenceID, err)
	}

	return graphResponsePointer, nil
}

func (c *Client) PublishChanges() (bool, error) {
	return true, nil
}

func (c *Client) DiscardChanges() (bool, error) {
	discardChanges, err := c.MakeGraphQLRequest(context.Background(), `
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
