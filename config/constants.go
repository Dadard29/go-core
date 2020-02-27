package config

import "time"

// misc
var PrivateKeyFile = "private.pem"
var JwtValidityDuration = 24 * time.Hour
var AuthorizationHeader = "Authorization"
var InvalidToken = "invalid token"

const (
	Api = "api"
)

const (
	Notification                    = "notification"
	NotificationBot                 = "bot"
	NotificationBotWebhookSecretKey = "webhookSecretKey"
)

const (
	Profile             = "profile"
	ProfileJwt          = "jwt"
	ProfileJwtSecretKey = "secretKey"
)

const (
	Connectors                                    = "connectors"
	ConnectorsTelegram                            = "telegram"
	ConnectorsTelegramBotTokenKey                 = "botTokenKey"
	ConnectorsTelegramApiUrlSendMessage           = "apiUrlSendMessage"
	ConnectorsTelegramParseModeMarkdown           = "parseModeMarkdown"
	ConnectorsTelegramParseModeHTML               = "parseModeHTML"
	ConnectorsTelegramMonitoringChatId            = "monitoringChatId"
	ConnectorsTelegramContinuousIntegrationChatId = "continuousIntegrationChatId"
	ConnectorsTelegramPrivateChatId               = "privateChatId"
)
