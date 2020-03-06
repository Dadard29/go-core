package controllers

import (
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

// GET
// Authorization: 	JWT in header Authorization
// Params: 			None
// Body: 			None
func SubsList(w http.ResponseWriter, r *http.Request) {
	// auth
	profileKey, err := getProfileKey(r)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	subList, message, err := managers.SubsManagerList(profileKey)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, subList, w)
	logger.CheckErr(err)
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			accessToken
// Body: 			None
func SubsCheckExists(w http.ResponseWriter, r *http.Request) {
	// auth
	var err error
	profileKey, err := getProfileKey(r)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	subToken := r.URL.Query().Get("accessToken")
	apiName := r.URL.Query().Get("apiName")

	var status bool
	var message string

	if subToken == "" && apiName == "" {
		// no param given
		api.Api.BuildMissingParameter(w)
		return
	} else if subToken != "" && apiName != "" {
		// bot params given
		api.Api.BuildErrorResponse(http.StatusBadRequest, "args overload", w)
		return
	} else if subToken != "" {
		status, message, err = managers.SubsManagerExistsFromToken(subToken)
	} else if apiName != "" {
		status, message, err = managers.SubsManagerExistsFromApiName(apiName, profileKey)
	}


	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(status, message, "", w)
	logger.CheckErr(err)
}

// POST
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func Subscribe(w http.ResponseWriter, r *http.Request) {
	// auth
	profileKey, err := getProfileKey(r)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	apiName := r.URL.Query().Get("apiName")
	if apiName == "" {
		err := api.Api.BuildMissingParameter(w)
		logger.CheckErr(err)
		return
	}

	subList, message, err := managers.SubsManagerCreate(profileKey, apiName)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, subList, w)
	logger.CheckErr(err)
}

// DELETE
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	// auth
	profileKey, err := getProfileKey(r)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	apiName := r.URL.Query().Get("apiName")
	if apiName == "" {
		err := api.Api.BuildMissingParameter(w)
		logger.CheckErr(err)
		return
	}

	subList, message, err := managers.SubsManagerDelete(profileKey, apiName)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, subList, w)
	logger.CheckErr(err)
}
