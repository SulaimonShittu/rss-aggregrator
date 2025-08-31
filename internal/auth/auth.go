package auth

import (
	"errors"
	"net/http"
	"strings"
)

// getAPIKey extracts the API key from
// the header of an HTTP request
// Example:
// Authorization: ApiKey {insert apikey here}

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of the header")
	}
	return vals[1], nil
}
