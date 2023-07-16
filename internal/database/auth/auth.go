package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Get API key from headers of HTTP request
// Authorization: ApiKey {apikey here}

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("incorrect auth header format")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("incorrect auth header index[0]")
	}
	return vals[1], nil
}
