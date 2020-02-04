package repositories

import (
	"github.com/Dadard29/go-core/models"
)

func NotificationBotWebookRepository(webhook models.GitlabWebhookPipeline) (bool, string) {
	return true, "ok"
}
