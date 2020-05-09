package controllers

import (
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

// GET
// Authorization: 	JWT
// Params: 			None
// Body: 			None
func ApiListGet(w http.ResponseWriter, r *http.Request) {
	if !checkJwt(r) {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	apiList, message, err := managers.ApiManagerList()
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, apiList, w)
	logger.CheckErr(err)
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func ApiGet(w http.ResponseWriter, r *http.Request) {
	if !checkJwt(r) {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	apiName := r.URL.Query().Get("apiName")
	if apiName == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	apiList, message, err := managers.ApiManagerGet(apiName)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, apiList, w)
	logger.CheckErr(err)
}