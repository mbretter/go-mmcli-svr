package api

import (
	"fmt"
	"github.com/mbretter/go-mmcli-svr/backend"
	"regexp"
	"strings"
)

type Api struct {
	backend backend.Backend
}

func Provide(be backend.Backend) *Api {
	api := Api{backend: be}
	return &api
}

type Response struct {
	Message string `json:"message"`
}

// ValidatePathIndex unescapes and validates a Modem ID string by replacing "%2f" and "%2F" with "/"
func ValidatePathIndex(id string) (string, error) {
	if len(id) == 0 {
		return "", nil
	}

	id = strings.ReplaceAll(id, "%2f", "/")
	id = strings.ReplaceAll(id, "%2F", "/")
	regEx, err := regexp.Compile("^[a-zA-Z0-9/]+$")
	if err != nil {
		return "", err
	}

	if !regEx.MatchString(id) {
		return "", fmt.Errorf("invalid ID")
	}
	return id, nil
}
