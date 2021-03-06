package repositories

import (
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/connectors"
	"github.com/Dadard29/go-core/models"
)

func RecoveryUpdate(profile models.Profile, recoverBy string, contact string) {
	profile.RecoverBy = recoverBy
	profile.Contact = contact

	api.Api.Database.Orm.Save(&profile)
}

func RecoveryDelete(profile models.Profile) {
	profile.RecoverBy = ""
	profile.Contact = ""

	api.Api.Database.Orm.Save(&profile)
}

// mail
func RecoverySendTestMail(p models.Profile) error {
	html := "<h3>Test recovery settings</h3>\n" +
		"<p>" +
		"Hello my dude," +
		"</p>" +
		"<p>" +
		"If you see this message, it means your recovery settings are correctly set\n\n" +
		"</p>" +
		"<p>" +
		"Have a pleasant day" +
		"</p>" +
		"<p>" +
		"-- dadard" +
		"</p>"

	return emailConnector.SendMail(p.Contact, "Test recovery settings", html)
}

func RecoverySendCodeMail(p models.Profile, code string) error {
	html := "<h3>Recovering account</h3>" +
		"<p>" +
		"Hello my dude," +
		"</p>" +
		"<p>" +
		"Seems you are trying to recover your account " +
		"in <a href=\"https://dadard.fr\">dadard-website</a>." +
		"<p>" +
		"<p>" +
		"To do so, please use this confirmation code:" +
		"</p>" +
		"<h2>" + code + "</h2>" +
		"<p>Have a pleasant day</p>"

	return emailConnector.SendMail(p.Contact, "Recovering account", html)
}

func RecoverySendNotificationMail(p models.Profile, text string) error {

	html := "<p>" +
		text +
		"</p>" +
		"<p>Have a pleasant day</p>" +
		"<p>NB: <i>you receive these messages because you " +
		"activated the notifications on <a href=\"https://dadard.fr/profile\">dadard-website</a>." +
		"You can deactivate this setting at any time</i>"

	return emailConnector.SendMail(p.Contact, "Notification", html)
}

// telegram
func RecoverySendTestTelegram(profile models.Profile) error {
	msg := fmt.Sprintf(
		"*Test recovery settings*\n\n" +
			"Hello my dude,\n\n" +
			"If you see this message, it means your recovery settings are correctly set\n\n" +
			"Have a pleasant day\n")

	return telegramConnector.SendMessage(msg, profile.Contact, connectors.ParseModeMarkdown)
}

func RecoverySendCodeTelegram(p models.Profile, code string) error {
	msg := fmt.Sprintf(
		"*Recovering account*\n\n"+
			"Seems you are trying to recover your account "+
			"in [dadard-website](https://dadard.fr)\n\n"+
			"To do so, please use this confirmation code:\n\n"+
			"*%s*\n\n"+
			"Have a pleasant day\n",
		code)

	return telegramConnector.SendMessage(msg, p.Contact, connectors.ParseModeMarkdown)
}

func RecoverySendNotificationTelegram(p models.Profile, text string) error {

	msg := fmt.Sprintf(
		"*Notifications*\n\n"+
			"%s\n\n"+
			"Have a pleasant day\n\n"+
			"NB: _you receive these messages because you "+
			"activated the notifications on_ [dadard-website](https://dadard.fr/profile). "+
			"_You can deactivate this setting at any time_",
		text)

	return telegramConnector.SendMessage(msg, p.Contact, connectors.ParseModeMarkdown)
}

// temp profile

// I use the same table for the lost passwords
// in case of recovery, the password field is null
