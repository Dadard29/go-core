package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ProfileGet(w, r)
	} else if r.Method == http.MethodPost {
		ProfileSignUp(w, r)
	} else if r.Method == http.MethodPut {
		ProfileChangePassword(w, r)
	} else if r.Method == http.MethodDelete {
		ProfileDelete(w, r)
	} else {
		API.BuildMethodNotAllowedResponse(w)
	}
}

// POST
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func ProfileSignUp(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	pk, message, err := managers.ProfileManagerSignUp(username, password)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
	} else {
		err := API.BuildJsonResponse(true, message, pk, w)
		logger.CheckErr(err)
	}
}

// GET
// Authorization: 	Basic
// Params: 			None
// Body: 			None
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	profile, message, err := managers.ProfileManagerSignIn(username, password)
	if err != nil {
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, models.ProfileJson{
		ProfileKey:  profile.ProfileKey,
		Username:    profile.Username,
		DateCreated: profile.DateCreated,
	}, w)
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
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	var body models.ProfileChangePassword
	if err := API.ParseJsonBody(r, &body); err != nil || body.NewPassword == "" {
		err := API.BuildErrorResponse(http.StatusBadRequest, "wrong body format", w)
		logger.CheckErr(err)
		return
	}

	profile, message, err := managers.ProfileManagerChangePassword(username, password, body.NewPassword)
	if err != nil {
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, models.ProfileJson{
		ProfileKey:  profile.ProfileKey,
		Username:    profile.Username,
		DateCreated: profile.DateCreated,
	}, w)
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
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	profileDeleted, message, err := managers.ProfileManagerDelete(username, password)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, models.ProfileJson{
		ProfileKey:  profileDeleted.ProfileKey,
		Username:    profileDeleted.Username,
		DateCreated: profileDeleted.DateCreated,
	}, w)
	logger.CheckErr(err)
}
