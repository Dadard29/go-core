package controllers

import (
	"fmt"
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

	confirmationCodeKey = "confirmation_code"
	contactKey = "contact"
)

// POST
// Authorization: 	Basic
// Params: 			confirm_by, contact
// Body: 			None

// ask for account creation, confirmation done by using confirm_by specified way
func ProfileSignUp(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		return
	}

	// checks parameters
	way := r.URL.Query().Get(signUpWayKey)
	if way == "" {
		api.Api.BuildMissingParameter(w)
		return
	}
	if !managers.ConfirmationWayIsValid(way) {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "confirmation way not found", w)
		return
	}

	contact := r.URL.Query().Get(contactKey)
	if contact == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	tmpP, msg, err := managers.ProfileManagerCreateTemp(username, password)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}
	logger.Debug(fmt.Sprintf("temp profile %s created", username))

	if msg, err := managers.ProfileManagerSendConfirmationCode(way, contact, tmpP); err != nil {
		logger.Error(err.Error())
		managers.ProfileManagerDeleteTemp(username)
		logger.Debug(fmt.Sprintf("temp profile %s deleted", username))
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, "notification sent", nil, w)
}

// POST
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

	confirmationCode := r.URL.Query().Get(confirmationCodeKey)
	if confirmationCode == "" {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "missing parameter", w)
		return
	}

	// confirm the code
	msg, err := managers.ProfileManagerConfirmCode(username, password, confirmationCode)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	profile, message, err := managers.ProfileManagerCreate(username, password, "", "")
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, profileErrorMsg, w)
	} else {
		api.Api.BuildJsonResponse(true, message, profile, w)
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
