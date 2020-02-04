package controllers

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)

func NotificationBotWebookRoute(w http.ResponseWriter, r *http.Request) {
	var err error

	secretToken := r.Header.Get("X-Gitlab-Token")

	keyWebhookToken, err := api.Api.Config.GetValueFromFile("notification", "bot", "webhookSecretTokenKey")
	logger.CheckErr(err)

	webhookToken := api.Api.Config.GetEnv(keyWebhookToken)

	if secretToken != webhookToken {
		err := API.BuildErrorResponse(http.StatusUnauthorized, "wrong webhook token", w)
		logger.CheckErr(err)
		return
	}

	status, message := managers.NotificationBotWebookManager(r.Body)
	if status == false {
		err = API.BuildErrorResponse(http.StatusInternalServerError, message, w)
	} else {
		err = API.BuildJsonResponse(true, message, "", w)
	}

	logger.CheckErr(err)
}

