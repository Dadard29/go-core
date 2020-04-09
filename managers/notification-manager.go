package managers

import (
	"encoding/json"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"io"
	"io/ioutil"
)

func NotificationBotWebookManager(body io.ReadCloser) (error, string) {
	var webhook models.GitlabWebhookPipeline

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err, "error decoding input body"
	}

	err = json.Unmarshal(data, &webhook)
	if err != nil {
		return err, "error decoding json body"
	}

	err = repositories.NotificationBotWebookRepository(webhook)
	if err != nil {
		return err, "error sending notif"
	}

	return nil, "notif sent"
}
