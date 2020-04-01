package repositories

import (
	"fmt"
	"github.com/Dadard29/go-core/connectors"
	"github.com/Dadard29/go-core/models"
)

func NotificationBotWebookRepository(webhook models.GitlabWebhookPipeline) (bool, string) {
	telegram, err := connectors.NewCITelegramConnector()
	logger.CheckErrFatal(err)

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

	if err != nil {
		return false, err.Error()
	}

	parseMode := connectors.ParseModeMarkdown

	err = telegram.SendMessage(messageFormat, parseMode)
	if err != nil {
		return false, err.Error()
	}

	return true, "webhook notification sent"
}
