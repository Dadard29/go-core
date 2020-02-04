package managers

import (
	"encoding/json"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"io"
	"io/ioutil"
)

func NotificationBotWebookManager(body io.ReadCloser) (bool, string) {
	var webhook models.GitlabWebhookPipeline

	data, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Error(err.Error())
		return false, "error decoding input body"
	}

	err = json.Unmarshal(data, &webhook)
	if err != nil {
		logger.Error(err.Error())
		return false, "error decoding json body"
	}

	return repositories.NotificationBotWebookRepository(webhook)
}
