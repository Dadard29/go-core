package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/connectors"
	"github.com/Dadard29/go-core/models"
	"strconv"
)

func ProfileCreate(p models.Profile) (models.Profile, string, error) {
	if profileExists(p) {
		message := "existing profile with same username"
		return models.Profile{}, message, errors.New(message)
	}

	api.Api.Database.Orm.Create(&p)

	if !profileExists(p) {
		message := "error creating profile"
		return models.Profile{}, message, errors.New(message)
	}

	return p, "profile created", nil
}

// create a new temporary profile in database
func TempProfileCreate(p models.TempProfile) error {

	api.Api.Database.Orm.Create(&p)
	return nil
}

func TempProfileGet(username string) (models.TempProfile, error) {
	var out models.TempProfile
	api.Api.Database.Orm.Where(&models.TempProfile{
		Username: username,
	}).Find(&out)

	if out.PasswordEncrypt == "" || out.Username == "" {
		return models.TempProfile{}, errors.New("temp profile not found")
	}

	return out, nil
}

func TempProfileDelete(username string) (models.TempProfile, error) {
	p, err := TempProfileGet(username)
	if err != nil {
		return models.TempProfile{}, err
	}

	api.Api.Database.Orm.Delete(&p)

	return p, nil
}

func SendConfirmationMail(email string, tmpProfile models.TempProfile) error {
	mailConfig, err := api.Api.Config.GetSubcategoryFromFile("connectors", "mail")
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(mailConfig["port"])
	if err != nil {
		return err
	}

	password := api.Api.Config.GetEnv(mailConfig["passwordKey"])

	c := connectors.NewMailConnector(mailConfig["user"],
		connectors.SmtpConfig{
			Host:     mailConfig["host"],
			Port:     port,
			User:     mailConfig["user"],
			Password: password,
		})

	formattedExpirationTime := tmpProfile.ExpirationTime.Format("15:04:05")

	html := "<h3>Identify confirmation</h3>\n" +
		"<p>" +
		"You have requested an account creation in" +
		" <a href=\"https://dadard.fr\">dadard-website</a>." +
		" In order to confirm you identity, please use the" +
		" following confirmation code:" +
		"</p>" +
		"<h1>" + tmpProfile.ConfirmationCode + "</h1>" +
		"<p>" +
		"<i>Note: this code will expire at " + formattedExpirationTime + "</i>" +
		"</p>"

	err = c.SendMail(email, "Identify confirmation", html)

	return err
}

func SendConfirmationTelegram(telegramUserId string, tmpProfile models.TempProfile) error {
	c, err := connectors.NewStandardTelegramConnector(telegramUserId)
	if err != nil {
		return err
	}

	formattedExpirationTime := tmpProfile.ExpirationTime.Format("15:04:05")

	msg := fmt.Sprintf(
		"*Identity confirmation*\n" +
			"You have requested an account creation " +
			"in [dadard-website](https://dadard.fr)\n\n" +
			"In order to confirm your identify, please use " +
			"this confirmation code:\n\n" +
			"*%s*\n\n" +
			"(this code will expire at %v)\n",
			tmpProfile.ConfirmationCode,
			formattedExpirationTime)

	return c.SendMessage(msg, connectors.ParseModeMarkdown)
}

