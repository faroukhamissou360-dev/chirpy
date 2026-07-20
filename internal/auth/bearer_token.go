package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader != "" {
		tokenBearer := strings.Fields(authorizationHeader)
		return tokenBearer[1], nil
	}
	return "", errors.New("Token string does not exist")

}
