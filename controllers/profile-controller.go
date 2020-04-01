package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
	"net/http"
)

const (
	profileErrorMsg = "bad username/password"

	signUpWayKey = "confirm_by"
	signUpWayTelegram = "telegram"
	signUpWayEmail = "email"

	emailKey = "email"
	telegramIdKey = "telegram_id"
)

// POST
// Authorization: 	Basic
// Params: 			confirm_by, <address or id...>
// Body: 			None

// ask for account creation, confirmation done by using confirm_by specified way
func ProfileSignUp(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		return
	}

	way := r.URL.Query().Get(signUpWayKey)
	if way == signUpWayEmail {
		// send ask for confirmation by mail
		email := r.URL.Query().Get(emailKey)
		msg, err := managers.ProfileManagerSendConfirmationMail(username, password, email)
		if err != nil {
			logger.Error(err.Error())
			api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
			return
		}

	} else if way == signUpWayTelegram {
		// send ask for confirmation with telegram
		telegramId := r.URL.Query().Get(telegramIdKey)
		msg, err := managers.ProfileManagerSendConfirmationTelegram(username, password, telegramId)
		if err != nil {
			logger.Error(err.Error())
			api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
			return
		}

	} else {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "confirmation way not found", w)
		return
	}
}

// GET
// Authorization: 	Basic
// Params: 			confirmation_code
// Body: 			None
func ProfileSignUpConfirm(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		return
	}

	confirmationCode := r.URL.Query().Get("confirmation_code")
	if confirmationCode == "" {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "missing parameter", w)
		return
	}

	// confirm the code

	profile, message, err := managers.ProfileManagerSignUp(username, password)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, profileErrorMsg, w)
		logger.CheckErr(err)
	} else {
		err := api.Api.BuildJsonResponse(true, message, profile, w)
		logger.CheckErr(err)
	}
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			None
// Body: 			None
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	// auth
	profileKey, err := getProfileKey(r)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
		return
	}

	profile, message, err := managers.ProfileManagerGet(profileKey)
	if err != nil {
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, profileErrorMsg, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, profile, w)
	logger.CheckErr(err)
}

// PUT
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func ProfileChangePassword(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	var body models.ProfileChangePassword
	if err := API.ParseJsonBody(r, &body); err != nil || body.NewPassword == "" {
		err := api.Api.BuildErrorResponse(http.StatusBadRequest, "wrong body format", w)
		logger.CheckErr(err)
		return
	}

	profile, message, err := managers.ProfileManagerChangePassword(username, password, body.NewPassword)
	if err != nil {
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, profile, w)
	logger.CheckErr(err)
}

// DELETE
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func ProfileDelete(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	profileDeleted, message, err := managers.ProfileManagerDelete(username, password)
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, profileErrorMsg, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, profileDeleted, w)
	logger.CheckErr(err)
}
