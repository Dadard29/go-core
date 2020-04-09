package repositories

import (
	"github.com/Dadard29/go-api-utils/log"
	"github.com/Dadard29/go-api-utils/log/logLevel"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/connectors"
	"strconv"
)

var logger = log.NewLogger("REPOSITORY", logLevel.DEBUG)

// connectors
// cuz fuck it

var telegramConnector *connectors.TelegramConnector
var telegramCiConnector *connectors.TelegramConnector
var emailConnector *connectors.MailConnector

func SetTelegramConnectors() error {
	var err error
	stdBotTokenKey, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramStdBotTokenKey)
	if err != nil {
		return err
	}

	stdBotToken := api.Api.Config.GetEnv(stdBotTokenKey)
	telegramConnector, err = connectors.NewTelegramConnector(stdBotToken)
	if err != nil {
		return err
	}

	ciBotTokenKey, err := api.Api.Config.GetValueFromFile(
		config.Connectors,
		config.ConnectorsTelegram,
		config.ConnectorsTelegramCiBotTokenKey)
	if err != nil {
		return err
	}

	ciBotToken := api.Api.Config.GetEnv(ciBotTokenKey)
	telegramCiConnector, err = connectors.NewTelegramConnector(ciBotToken)
	if err != nil {
		return err
	}

	return nil
}

func SetEmailConnectors() error {
	var err error
	mailConfig, err := api.Api.Config.GetSubcategoryFromFile(
		config.Connectors,
		config.ConnectorsEmail)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(mailConfig["port"])
	if err != nil {
		return err
	}

	password := api.Api.Config.GetEnv(mailConfig["passwordKey"])

	emailConnector = connectors.NewMailConnector(mailConfig["user"],
		connectors.SmtpConfig{
			Host:     mailConfig["host"],
			Port:     port,
			User:     mailConfig["user"],
			Password: password,
		})

	return nil
}
