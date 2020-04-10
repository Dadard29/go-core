package config

import "time"

// misc
var PrivateKeyFile = "private.pem"
var JwtValidityDuration = 24 * time.Hour * 7 // valid 7 days
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
	Profile               = "profile"
	ProfileJwt            = "jwt"
	ProfileJwtSecretKey   = "secretKey"
	ProfileQuota          = "quota"
	ProfileQuotaSilver    = "silver"
	ProfileQuotaNotSilver = "notSilver"
	ProfileCreation = "creation"
	ProfileCreationConfirmationExpirationDuration = "confirmationExpirationDuration"
)

const (
	Session                 = "session"
	SessionAuth             = "auth"
	SessionAuthUsernameKey  = "usernameKey"
	SessionAuthUPasswordKey = "passwordKey"
)

const (
	Connectors                                    = "connectors"
	ConnectorsTelegram                            = "telegram"

	ConnectorsTelegramContinuousIntegrationChatId = "chat-ci"
	ConnectorsTelegramMonitoringChatId = "chat-monitoring"
	ConnectorsTelegramPrivateChatId = "chat-private"

	ConnectorsTelegramCiBotTokenKey = "bot-ci"
	ConnectorsTelegramStdBotTokenKey = "bot-std"

	ConnectorsEmail = "mail"
)
