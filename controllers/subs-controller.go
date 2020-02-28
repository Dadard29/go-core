package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
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

	pl := jwt.Infos.(map[string]interface{})
	profileKey := pl["ProfileKey"].(string)

	if r.Method == http.MethodGet {
		SubsList(w, profileKey)
	} else if r.Method == http.MethodPost {
		Subscribe(w, r, profileKey)
	} else if r.Method == http.MethodDelete {
		Unsubscribe(w, r, profileKey)
	} else if r.Method == http.MethodOptions {
		SubsCheckExists(w, r)
	} else {
		err := API.BuildMethodNotAllowedResponse(w)
		logger.CheckErr(err)
	}
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			None
// Body: 			None
func SubsList(w http.ResponseWriter, profileKey string) {
	subList, message, err := managers.SubsManagerList(profileKey)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, subList, w)
	logger.CheckErr(err)
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			accessToken
// Body: 			None
func SubsCheckExists(w http.ResponseWriter, r *http.Request) {
	subToken := r.URL.Query().Get("accessToken")
	if subToken == "" {
		err := API.BuildMissingParameter(w)
		logger.CheckErr(err)
		return
	}

	status, message, err := managers.SubsManagerExists(subToken)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(status, message, "", w)
	logger.CheckErr(err)
}

// POST
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func Subscribe(w http.ResponseWriter, r *http.Request, profileKey string) {
	apiName := r.URL.Query().Get("apiName")
	if apiName == "" {
		err := API.BuildMissingParameter(w)
		logger.CheckErr(err)
		return
	}

	subList, message, err := managers.SubsManagerCreate(profileKey, apiName)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, subList, w)
	logger.CheckErr(err)
}

// DELETE
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func Unsubscribe(w http.ResponseWriter, r *http.Request, profileKey string) {
	apiName := r.URL.Query().Get("apiName")
	if apiName == "" {
		err := API.BuildMissingParameter(w)
		logger.CheckErr(err)
		return
	}

	subList, message, err := managers.SubsManagerDelete(profileKey, apiName)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, subList, w)
	logger.CheckErr(err)
}
