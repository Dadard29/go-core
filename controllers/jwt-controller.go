package controllers

import (
	"errors"
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
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
		API.BuildMethodNotAllowedResponse(w)
	}
}

func getJwtSecret() string {
	jwtSecretKey, err := api.Api.Config.GetValueFromFile(
		config.Profile,
		config.ProfileJwt,
		config.ProfileJwtSecretKey)

	logger.CheckErrFatal(err)

	secret := api.Api.Config.GetEnv(jwtSecretKey)
	if secret == "" {
		logger.CheckErrFatal(errors.New("no configured jwt secret"))
	}

	return secret
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
	secret := getJwtSecret()

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
		return
	}

	jwtCiphered := body.JwtCiphered
	jwtDeciphered, err := auth.DecipherJwtWithJwe(config.PrivateKeyFile, []byte(jwtCiphered))
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, "error deciphering token", w)
		logger.CheckErr(err)
		return
	}

	secret := getJwtSecret()
	_, err = auth.VerifyJwtHS256(jwtDeciphered, secret)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusBadRequest, "token invalid", w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, "token valid", "", w)
	logger.CheckErr(err)
}
