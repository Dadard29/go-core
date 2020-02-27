package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ApiGet(w, r)
	} else {
		API.BuildMethodNotAllowedResponse(w)
	}
}

// GET
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func ApiGet(w http.ResponseWriter, r *http.Request) {
	if status, _, msg := managers.ValidateJwtBody(
		r.Header.Get(config.AuthorizationHeader)); !status {
		err := API.BuildErrorResponse(http.StatusForbidden, msg, w)
		logger.CheckErr(err)
		return
	}

	api, message, err := managers.ApiManagerGet(r.URL.Query().Get("apiName"))
	if err != nil {
		logger.Error(err.Error())
		err := API.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = API.BuildJsonResponse(true, message, api, w)
	logger.CheckErr(err)
}

// PRIVILEGE SYSTEM
// a privileged session could be created via an endpoint (called by a client like
// a telegram bot, THE CHANNEL MUST BE SECURED). Then, the session could be available
// with a token generated. The session would have an expired time. The session creation
// would be protected by password

// POST
// Authorization: 	None
// Params: 			None
// Body: 			models.Api

// the Api creation must be reserved to the admin
// this endpoint MUST NOT be activated now
// need to think on a privilege system before
func ApiCreate(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PUT
// Authorization: 	None
// Params: 			None
// Body: 			None

// the Api modification must be reserved to the admin
// this endpoint MUST NOT be activated now
// need to think on a privilege system before
func ApiPut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// DELETE
// Authorization: 	None
// Params: 			None
// Body: 			None

// the Api modification must be reserved to the admin
// this endpoint MUST NOT be activated now
// need to think on a privilege system before
func ApiDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
}
