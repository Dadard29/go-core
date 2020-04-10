package controllers

import (
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
	"net/http"
)

// POST
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func JwtCreate(w http.ResponseWriter, r *http.Request) {
	var err error

	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	profile, message, err := managers.ProfileManagerSignIn(username, password)
	if err != nil {
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, profileErrorMsg, w)
		logger.CheckErr(err)
		return
	}

	logger.Debug(message)
	secret := managers.GetJwtSecret()

	pl := models.JwtPayload{
		ProfileKey: profile.ProfileKey,
	}

	// create token
	token, err := auth.NewJwtHS256(
		secret,
		"core", "", []string{"https://dadard.fr"}, config.JwtValidityDuration,
		pl)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError,
			"error forging token",
			w)
		logger.CheckErr(err)
		return
	}

	// cipher token
	cipheredToken, err := auth.CipherJwtWithJwe(config.PrivateKeyFile, token)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError,
			"error ciphering token",
			w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, "token forged and ciphered", string(cipheredToken), w)
	logger.CheckErr(err)
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			None
// Body: 			None
func JwtValidate(w http.ResponseWriter, r *http.Request) {
	if status := managers.ValidateJwtCiphered(
		r.Header.Get(config.AuthorizationHeader)); status == nil {
		err := api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		logger.CheckErr(err)
	} else {
		err := api.Api.BuildJsonResponse(true, "valid token", nil, w)
		logger.CheckErr(err)
	}
}
