package controllers

import (
	"encoding/json"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/managers"
	"github.com/Dadard29/go-core/models"
	"io/ioutil"
	"net/http"
)

const (
	profileKeyParam = "profileKey"
	recoverByParam = "recover_by"
	contactParam = "contact"

	tokenHeader = "X-Echo-Slam"
	tokenKey = "ECHO_SLAM_TOKEN"
)

func checkEchoSlamToken(r *http.Request) bool {
	return api.Api.Config.GetEnv(tokenKey) == r.Header.Get(tokenHeader)
}

// GET
// Authorization: 	Basic + echo-slam header
// Params: 			recoverByParam, contactParam
// Body: 			None
func SignUpFromEchoSlam(w http.ResponseWriter, r *http.Request) {
	if !checkEchoSlamToken(r) {
		api.Api.BuildErrorResponse(http.StatusUnauthorized, "unauthorized", w)
		return
	}

	recoverBy := r.URL.Query().Get(recoverByParam)
	contact := r.URL.Query().Get(contactParam)
	if recoverBy == "" || contact == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	var username, password string
	var ok bool
	if username, password, ok = r.BasicAuth(); !ok {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "wrong auth format", w)
		return
	}

	p, msg, err := managers.ProfileManagerCreate(username, password, recoverBy, contact)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
		return
	}

	api.Api.BuildJsonResponse(true, msg, p, w)
}


// POST
// Authorization: 	echo-slam header
// Params: 			profileKeyParam
// Body: 			None
func SubscribeFromEchoSlam(w http.ResponseWriter, r *http.Request) {
	if !checkEchoSlamToken(r) {
		api.Api.BuildErrorResponse(http.StatusUnauthorized, "unauthorized", w)
		return
	}

	pk := r.URL.Query().Get(profileKeyParam)
	if pk == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	// get the list of APIs to subscribe to
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "invalid body", w)
		return
	}

	var body []string
	err = json.Unmarshal(data, &body)
	if err != nil {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "invalid json body", w)
		return
	}

	if len(body) == 0 {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "empty list", w)
		return
	}

	var subs []models.SubscriptionJson
	for _, a := range  body {
		s, msg, err := managers.SubsManagerCreate(pk, a, true)
		if err != nil {
			logger.Error(err.Error())
			api.Api.BuildErrorResponse(http.StatusInternalServerError, msg, w)
			return
		}

		subs = append(subs, s)
	}
	api.Api.BuildJsonResponse(true, "profile setup", subs, w)
}
