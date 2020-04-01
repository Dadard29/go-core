package connectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiUrlSendMessage = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=%s"
	ParseModeMarkdown = "Markdown"
	ParseModeHTML = "HTML"
)

type SendMessageResponseError struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type TelegramConnector struct {
	botToken          string
	httpClient        *http.Client
	chatId string
}

func NewStandardTelegramConnector(targetUserId string) (*TelegramConnector, error) {

	botTokenKey, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramStdBotTokenKey)
	if err != nil {
		return nil, err
	}

	botToken := api.Api.Config.GetEnv(botTokenKey)
	return newTelegramConnector(botToken, targetUserId)
}

func NewCITelegramConnector() (*TelegramConnector, error) {
	ciChatId, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramContinuousIntegrationChatId)
	if err != nil {
		return nil, err
	}
	botTokenKey, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramCiBotTokenKey)
	if err != nil {
		return nil, err
	}

	botToken := api.Api.Config.GetEnv(botTokenKey)
	return newTelegramConnector(botToken, ciChatId)
}

func newTelegramConnector(botToken string, chatId string) (*TelegramConnector, error) {

	if botToken == "" {
		return nil, errors.New("invalid bot token")
	}

	if chatId == "" {
		return nil, errors.New("invalid chat id")
	}

	return &TelegramConnector{
		botToken:          botToken,
		httpClient:        &http.Client{},
		chatId: chatId,
	}, nil
}

func (c *TelegramConnector) SendMessage(message string, parseMode string) error {
	// message MUST NOT BE url encoded, as it is encoded in the function

	messageUrlencoded := url.QueryEscape(message)
	urlFormat := fmt.Sprintf(apiUrlSendMessage, c.botToken, c.chatId, messageUrlencoded, parseMode)

	req, err := http.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s\n", urlFormat)
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
