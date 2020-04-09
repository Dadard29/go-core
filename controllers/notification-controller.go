package controllers

import (
	"errors"
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
		err := api.Api.BuildErrorResponse(http.StatusUnauthorized, "wrong webhook token", w)
		logger.CheckErr(err)
		return
	}

	err, message := managers.NotificationBotWebookManager(r.Body)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, message, w)
		return
	}

	api.Api.BuildJsonResponse(true, message, "", w)
}
