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

func NotificationActivate(username string, password string, beNotified bool) (string, error) {

	p, msg, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return msg, err
	}

	if msg, err := checkUsernamePassword(p, password); err != nil {
		return msg, err
	}

	p.BeNotified = beNotified

	if _, msg, err := repositories.ProfileUpdate(p); err != nil {
		return msg, err
	}

	return "notifications settings updated", nil
}
