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

func LostPassword(w http.ResponseWriter, r *http.Request) {
	// todo
}

func LostPasswordConfirmCode(w http.ResponseWriter, r *http.Request) {
	// todo
}