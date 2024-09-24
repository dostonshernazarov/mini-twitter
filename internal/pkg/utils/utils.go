package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/token"
)

type QueryParam struct {
	Filters  map[string]string
	Limit    uint64
	Page     uint64
	Ordering []string
	Search   string
}

func ParseQueryParam(queryParams map[string][]string) (*QueryParam, []string) {
	params := QueryParam{
		Filters:  make(map[string]string),
		Limit:    10,
		Page:     1,
		Ordering: []string{},
		Search:   "",
	}
	var errStr []string
	var err error

	for key, value := range queryParams {
		if key == "page" {
			params.Page, err = strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "invalid `page` param")
			}
			continue
		}

		if key == "limit" {
			params.Limit, err = strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "invalid `limit` param")
			}
			continue
		}

		if key == "search" {
			params.Search = value[0]
			continue
		}
		if key == "ordering" {
			params.Ordering = strings.Split(value[0], ",")
			continue
		}
		params.Filters[key] = value[0]
	}
	return &params, errStr
}

func GetClaimsFromToken(request *http.Request, cfg *config.Config) (map[string]interface{}, error) {
	token := request.Header.Get("Authorization")

	if token == "" {
		return map[string]interface{}{
			"role": "unauthorized",
			"sub":  nil,
			"exp":  nil,
			"iat":  nil,
		}, nil
	}

	if strings.Contains(token, "Bearer ") {
		token = strings.Split(token, "Bearer ")[1]
	}

	claims, err := tokens.ExtractClaim(token, []byte(cfg.SigningKey))
	if err != nil {
		return nil, err
	}

	return claims, nil
}
