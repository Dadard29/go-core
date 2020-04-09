package connectors

import (
	"encoding/json"
	"errors"
	"fmt"
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
}

func NewTelegramConnector(botToken string) (*TelegramConnector, error) {

	if botToken == "" {
		return nil, errors.New("invalid bot token")
	}

	return &TelegramConnector{
		botToken:          botToken,
		httpClient:        &http.Client{},
	}, nil
}

func (c *TelegramConnector) SendMessage(message string, chatId string, parseMode string) error {
	// message MUST NOT BE url encoded, as it is encoded in the function

	if chatId == "" {
		return errors.New("invalid chat ID")
	}

	messageUrlencoded := url.QueryEscape(message)
	urlFormat := fmt.Sprintf(apiUrlSendMessage, c.botToken, chatId, messageUrlencoded, parseMode)

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
