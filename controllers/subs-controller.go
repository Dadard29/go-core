package controllers

import (
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
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
// Authorization: 	JWT in header Authorization (if request with apiName parameter only)
// Params: 			accessToken or apiName
// Body: 			None
func SubsCheckExists(w http.ResponseWriter, r *http.Request) {
	var err error

	subToken := r.URL.Query().Get("accessToken")
	apiName := r.URL.Query().Get("apiName")

	var s models.SubscriptionJson
	var message string

	if subToken == "" && apiName == "" {
		// no param given
		api.Api.BuildMissingParameter(w)
		return
	} else if subToken != "" && apiName != "" {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "args overload", w)
		return
	} else if subToken != "" {
		// record token usage
		s, message, err = managers.SubsManagerGetFromToken(subToken)
		if err != nil {
			logger.Error(err.Error())
			err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
			logger.CheckErr(err)
			return
		}

	} else if apiName != "" {
		// informative way
		profileKey, err := getProfileKey(r)
		if err != nil {
			api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
			return
		}

		s, message, err = managers.SubsManagerGetFromApiName(apiName, profileKey)
		if err != nil {
			logger.Error(err.Error())
			err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
			logger.CheckErr(err)
			return
		}

	} else {
		api.Api.BuildMissingParameter(w)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, s, w)
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

// PUT
// Authorization: 	JWT In header Authorization
// Params: 			apiName
// Body: 			None
func SubRegenerate(w http.ResponseWriter, r *http.Request) {
	// auth
	profileKey, err := getProfileKey(r)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	apiName := r.URL.Query().Get("apiName")
	if apiName == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	subModified, message, err := managers.SubsManagerUpdate(profileKey, apiName)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		return
	}

	api.Api.BuildJsonResponse(true, message, subModified, w)
}

// DELETE
// Authorization: 	protected token
// Params: 			None
// Body: 			None

// warning: sensitive endpoint
func SubResetAll(w http.ResponseWriter, r *http.Request) {
	authToken := auth.ParseApiKey(r, "Authorization", true)

	protectedTokenKey, err := api.Api.Config.GetValueFromFile(
		config.Endpoints,
		config.EndpointsProtected,
		config.EndpointsProtectedKey)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error getting config", w)
		return
	}
	protectedToken := api.Api.Config.GetEnv(protectedTokenKey)
	if authToken != protectedToken {
		api.Api.BuildErrorResponse(http.StatusForbidden, "bad token", w)
		return
	}

	if msg, err := managers.SubsManagerResetRequestCount(); err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "requests count reset", nil, w)
}

// DELETE
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None

// frozen controller: see https://git.dadard.fr/go-dadard/go-core/issues/6
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

	subDeleted, message, err := managers.SubsManagerDelete(profileKey, apiName)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, subDeleted, w)
	logger.CheckErr(err)
}
