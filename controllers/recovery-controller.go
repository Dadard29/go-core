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

	if msg, err := managers.CheckUsernamePassword(username, password); err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, msg, w)
		return
	}

	// todo
}

func RecoverUpdate(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	if msg, err := managers.CheckUsernamePassword(username, password); err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, msg, w)
		return
	}

	// todo

}

func RecoverDelete(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	if msg, err := managers.CheckUsernamePassword(username, password); err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, msg, w)
		return
	}

	// todo
}

func LostPassword(w http.ResponseWriter, r *http.Request) {
	// todo
}

func LostPasswordConfirmCode(w http.ResponseWriter, r *http.Request) {
	// todo
}