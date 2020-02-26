package controllers

import (
	"errors"
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/managers"
	"net/http"
)


// POST
// Authorization: 	header X-Gitlab-Token (set from gitlab webapp)
// Params: 			None
// Body:			https://gitlab.com/help/user/project/integrations/webhooks (pipeline object)
func NotificationBotWebookRoute(w http.ResponseWriter, r *http.Request) {
	var err error

	secretToken := r.Header.Get("X-Gitlab-Token")

	keyWebhookToken, err := api.Api.Config.GetValueFromFile(
		config.Notification,
		config.NotificationBot,
		config.NotificationBotWebhookSecretKey)

	logger.CheckErr(err)
	webhookToken := api.Api.Config.GetEnv(keyWebhookToken)
	if webhookToken == "" {
		logger.CheckErrFatal(errors.New("no configured bot webhook secret"))
	}

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

