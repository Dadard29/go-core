package controllers

import (
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

// GET
// Authorization: 	JWT in header Authorization
// Params: 			apiName
// Body: 			None
func ApiGet(w http.ResponseWriter, r *http.Request) {
	if !checkJwt(r) {
		api.Api.BuildErrorResponse(http.StatusForbidden, config.InvalidToken, w)
	}

	apiList, message, err := managers.ApiManagerList()
	if err != nil {
		logger.Error(err.Error())
		err := api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		logger.CheckErr(err)
		return
	}

	err = api.Api.BuildJsonResponse(true, message, apiList, w)
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
	// TODO.md
}

// PUT
// Authorization: 	None
// Params: 			None
// Body: 			None

// the Api modification must be reserved to the admin
// this endpoint MUST NOT be activated now
// need to think on a privilege system before
func ApiPut(w http.ResponseWriter, r *http.Request) {
	// TODO.md
}

// DELETE
// Authorization: 	None
// Params: 			None
// Body: 			None

// the Api modification must be reserved to the admin
// this endpoint MUST NOT be activated now
// need to think on a privilege system before
func ApiDelete(w http.ResponseWriter, r *http.Request) {
	// TODO.md
}
