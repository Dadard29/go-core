package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
	"net/http"
)

func JwtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		JwtValidate(w, r)
	} else if r.Method == http.MethodPost {
		JwtCreate(w, r)
	} else {
		err := API.BuildMethodNotAllowedResponse(w)
		logger.CheckErr(err)
	}
}

// POST
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func JwtCreate(w http.ResponseWriter, r *http.Request) {
	var err error

	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	profile, message, err := managers.ProfileManagerSignIn(username, password)
	if err != nil {
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
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
		"core", "", []string{"http://dadard.fr"}, config.JwtValidityDuration,
		pl)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError,
			"error forging token",
			w)
		logger.CheckErr(err)
		return
	}

	// cipher token
	cipheredToken, err := auth.CipherJwtWithJwe(config.PrivateKeyFile, token)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError,
			"error ciphering token",
			w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, "token forged and ciphered", string(cipheredToken), w)
	logger.CheckErr(err)
}

// GET
// Authorization: 	None
// Params: 			None
// Body: 			{"JwtCiphered": "..."}
func JwtValidate(w http.ResponseWriter, r *http.Request) {
	var body models.JwtValidate
	if err := API.ParseJsonBody(r, &body); err != nil || body.JwtCiphered == "" {
		err := API.BuildErrorResponse(http.StatusBadRequest, "wrong body format", w)
		logger.CheckErr(err)
	}

	jwtCiphered := body.JwtCiphered

	if status:= managers.ValidateJwtCiphered(jwtCiphered); status == nil {
		err := API.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		logger.CheckErr(err)
	} else {
		err := API.BuildJsonResponse(true, "valid token", "", w)
		logger.CheckErr(err)
	}
}
