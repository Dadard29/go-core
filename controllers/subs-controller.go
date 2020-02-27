package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
	"net/http"
)

func SubsHandler(w http.ResponseWriter, r *http.Request) {
	jwt := managers.ValidateJwtCiphered(
		r.Header.Get(config.AuthorizationHeader))
	if jwt == nil {
		err := API.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		logger.CheckErr(err)
		return
	}

	pl := jwt.Infos.(models.JwtPayload)

	if r.Method == http.MethodGet {
		SubsList(w, pl.ProfileKey)
	} else if r.Method == http.MethodPost {
		Subscribe(w, r)
	} else if r.Method == http.MethodDelete {
		Unsubscribe(w, r)
	} else {
		API.BuildMethodNotAllowedResponse(w)
	}
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			None
// Body: 			None
func SubsList(w http.ResponseWriter, profileKey string) {
	api, message, err := managers.SubsManagerList(profileKey)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, api, w)
	logger.CheckErr(err)
}

// POST
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func Subscribe(w http.ResponseWriter, r *http.Request) {
	// todo
}

// DELETE
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	// todo
}
