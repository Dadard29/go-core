package repositories

import (
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/connectors"
	"github.com/Dadard29/go-core/models"
)

func NotificationBotWebookRepository(webhook models.GitlabWebhookPipeline) (bool, string) {
	telegram, err := connectors.NewTelegramConnector()
	logger.CheckErr(err)

	projectName := webhook.Project.Name
	projectUrl := webhook.Project.WebURL
	pipelineStatus := webhook.ObjectAttributes.Status
	createdAt := webhook.ObjectAttributes.CreatedAt
	user := webhook.User.Username
	pipelineDuration := webhook.ObjectAttributes.Duration

	message := "*%s* ([see on gitlab.com](%s))\n" +
		"- build status:\t *%s*\n" +
		"- created at:\t *%s*\n" +
		"- started by:\t *%s*\n" +
		"- duration:\t *%ds*\n"

	messageFormat := fmt.Sprintf(message,
		projectName, projectUrl, pipelineStatus, createdAt, user, pipelineDuration)

	chatId, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramContinuousIntegrationChatId)

	if err != nil {
		return false, err.Error()
	}

	parseMode, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramParseModeMarkdown)

	if err != nil {
		return false, err.Error()
	}

	err = telegram.SendMessage(messageFormat, chatId, parseMode)
	if err != nil {
		return false, err.Error()
	}

	return true, "webhook notification sent"
}
