package controllers

import (
	"errors"
	"github.com/Dadard29/go-api-utils/log"
	"github.com/Dadard29/go-api-utils/log/logLevel"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

var logger = log.NewLogger("CONTROLLER", logLevel.DEBUG)

func checkJwt(r *http.Request) bool {
	jwt := managers.ValidateJwtCiphered(
		r.Header.Get(config.AuthorizationHeader))

	return jwt != nil
}

func getProfileKey(r *http.Request) (string, error) {
	jwt := managers.ValidateJwtCiphered(
		r.Header.Get(config.AuthorizationHeader))
	if jwt == nil {
		return "", errors.New(config.InvalidToken)
	}

	pl := jwt.Infos.(map[string]interface{})
	profileKey := pl["ProfileKey"].(string)

	return profileKey, nil
}
