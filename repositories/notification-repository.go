package repositories

import (
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/connectors"
	"github.com/Dadard29/go-core/models"
)

func NotificationBotWebookRepository(webhook models.GitlabWebhookPipeline) error {

	projectName := webhook.Project.Name
	projectUrl := webhook.Project.WebURL
	pipelineStatus := webhook.ObjectAttributes.Status
	createdAt := webhook.ObjectAttributes.CreatedAt
	user := webhook.User.Username
	pipelineDuration := webhook.ObjectAttributes.Duration

	msg := fmt.Sprintf("*%s* ([see on gitlab.com](%s/pipelines))\n"+
		"- build status:\t *%s*\n"+
		"- created at:\t *%s*\n"+
		"- started by:\t *%s*\n"+
		"- duration:\t *%ds*\n",
		projectName, projectUrl, pipelineStatus, createdAt, user, pipelineDuration)

	ciChatId, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramContinuousIntegrationChatId)
	if err != nil {
		return err
	}

	return telegramCiConnector.SendMessage(msg, ciChatId, connectors.ParseModeMarkdown)
}
