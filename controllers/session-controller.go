package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	username, password, err := auth.ParseBasicAuth(r)
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong auth format", w)
		logger.CheckErr(err)
		return
	}

	check, err := managers.SessionManagerCheckAuth(username, password)
	if err != nil || ! check {
		if err != nil {
			logger.Error(err.Error())
		}

		err := API.BuildErrorResponse(http.StatusUnauthorized, "bad credentials", w)
		logger.CheckErr(err)
		return
	}


	if r.Method == http.MethodGet {
		SessionGet(w, r)
	} else if r.Method == http.MethodPost {
		SessionCreate(w, r)
	} else {
		err := API.BuildMethodNotAllowedResponse(w)
		logger.CheckErr(err)
		return
	}
}

// GET
// Authorization: 	Basic
// Params: 			None
// Body: 			None
// Check if a session is enabled
func SessionGet(w http.ResponseWriter, r *http.Request) {

}

// GET
// Authorization: 	Basic + check of remote addr
// Params: 			duration
// Body: 			None
func SessionCreate(w http.ResponseWriter, r *http.Request) {

}
