package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func HTTPRequestUnmarshal(c *http.Client, req *http.Request, responseBody interface{}) (*http.Response, error) {
	defer req.Body.Close()
	response, err := c.Do(req)
	if err != nil {
		return response, err
	}

	bResp, err := io.ReadAll(response.Body)
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(bResp, responseBody); err != nil {
		return response, err
	}

	return response, nil
}

func Map[T, U any](s []T, f func(T) U) []U {
	newSlice := make([]U, len(s))
	for i := range s {
		newSlice[i] = f(s[i])
	}

	return newSlice
}
