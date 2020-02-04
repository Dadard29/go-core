package connectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SendMessageResponseError struct {
	Ok bool `json:"ok"`
	Description string `json:"description"`
}

type TelegramConnector struct {
	botToken string
	apiUrlSendMessage string
	httpClient *http.Client
}

func NewTelegramConnector() (*TelegramConnector, error) {
	telegramConfig, err := api.Api.Config.GetSubcategoryFromFile(
		config.Connectors,
		config.ConnectorsTelegram)

	if err != nil {
		return nil, err
	}

	botTokenKey := telegramConfig[config.ConnectorsTelegramBotTokenKey]
	botTokenValue := api.Api.Config.GetEnv(botTokenKey)
	if botTokenValue == "" {
		return nil, errors.New("no configured bot token")
	}

	apiUrlSendMessage := telegramConfig[config.ConnectorsTelegramApiUrlSendMessage]

	return &TelegramConnector{
		botToken: botTokenValue,
		apiUrlSendMessage: apiUrlSendMessage,
		httpClient: &http.Client{},
	}, nil
}

func (c *TelegramConnector) SendMessage(message string, chatId string, parseMode string) error {
	urlFormat := fmt.Sprintf(c.apiUrlSendMessage, c.botToken, chatId, message, parseMode)

	req, err := http.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Status != strconv.Itoa(http.StatusOK) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var res SendMessageResponseError
		err = json.Unmarshal(body, &res)
		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("error while send message: %s", res.Description))
	} else {
		return nil
	}
}
