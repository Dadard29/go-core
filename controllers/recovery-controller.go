package controllers

import (
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

// POST
// Authorization: 	JWT
// Params: 			recover_by, contact
// Body: 			None
func RecoverySet(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	recoverBy := r.URL.Query().Get("recover_by")
	contact := r.URL.Query().Get("contact")
	if contact == "" || recoverBy == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	if msg, err := managers.RecoveryManagerUpdate(username, password, recoverBy, contact); err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "recovery settings set", nil, w)
}

// PUT
// Authorization: 	Basic
// Params: 			recover_by, contact
// Body: 			None
func RecoverUpdate(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	recoverBy := r.URL.Query().Get("recover_by")
	contact := r.URL.Query().Get("contact")
	if contact == "" || recoverBy == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	if msg, err := managers.RecoveryManagerUpdate(username, password, recoverBy, contact); err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "recovery settings updated", nil, w)
}

// DELETE
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func RecoverDelete(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	if msg, err := managers.RecoveryManagerDelete(username, password); err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "recover settings removed", nil, w)
}

// GET
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func RecoverTestGet(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	if msg, err := managers.RecoveryManagerTest(username, password); err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "test sent", nil, w)
}

// GET
// Authorization: 	None
// Params: 			username
// Body: 			None
func LostPassword(w http.ResponseWriter, r *http.Request) {
	// WARNING
	// NO AUTH - SO BE CAREFUL OF WHAT DATA IS EXPOSED

	username := r.URL.Query().Get("username")
	if username == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	if msg, err := managers.RecoverySendCode(username); err != nil {
		logger.Error(err.Error())
		logger.Error(msg)
	}

	// the only output given is this one
	api.Api.BuildJsonResponse(true,
		"if you have valid recovery settings, a notification has been sent",
		nil, w)
}

// POST
// Authorization: 	Basic
// Params: 			confirmation_code
// Body: 			None
func LostPasswordConfirmCode(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		return
	}

	confirmationCode := r.URL.Query().Get("confirmation_code")
	if msg, err := managers.RecoveryConfirmCode(username, password, confirmationCode); err != nil {
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "profile's password changed", nil, w)
}