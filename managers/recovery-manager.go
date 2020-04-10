package managers

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"time"
)

const (
	recoverByEmail = "email"
	recoverByTelegram = "telegram"
)

type RecoverByMapping struct {
	test func(p models.Profile) error
	sendCode func(p models.Profile, code string) error
	sendNotification func(p models.Profile, text string) error
}

var recoveryByMapping = map[string]RecoverByMapping {
	"email": {
		test:             repositories.RecoverySendTestMail,
		sendCode:         repositories.RecoverySendCodeMail,
		sendNotification: repositories.RecoverySendNotificationMail,
	},
	"telegram": {
		test:             repositories.RecoverySendTestTelegram,
		sendCode:         repositories.RecoverySendCodeTelegram,
		sendNotification: repositories.RecoverySendNotificationTelegram,
	},
}

func checkUsernamePassword(profile models.Profile, password string) (string, error) {
	if !models.ComparePassword(password, profile.PasswordEncrypt) {
		msg := "bad password"
		return msg, errors.New(msg)
	}

	return "password checked", nil
}

func RecoveryManagerUpdate(username string, password string, recoverBy string, contact string) (string, error) {
	p, msg, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return msg, err
	}

	if msg, err := checkUsernamePassword(p, password); err != nil {
		return msg, err
	}

	if _, check := recoveryByMapping[recoverBy]; !check {
		m := "unknown recovery channel"
		return m, errors.New(m)
	}

	repositories.RecoveryUpdate(p, recoverBy, contact)

	return "recovery settings all set", nil
}

func RecoveryManagerDelete(username string, password string) (string, error) {
	p, msg, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return msg, err
	}

	if msg, err := checkUsernamePassword(p, password); err != nil {
		return msg, err
	}

	repositories.RecoveryDelete(p)

	return "recovery setting removed", nil
}

func RecoverySendCode(username string) (string, error) {
	// WARNING
	// NO AUTH - SO BE CAREFUL OF WHAT DATA IS EXPOSED

	profileDb, msg, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return msg, err
	}

	if profileDb.RecoverBy == "" || profileDb.Contact == "" {
		msg := "your recovery settings are invalid: you are screwed"
		return msg, errors.New(msg)
	}

	if tmpP, err := repositories.TempProfileGet(username); err == nil {
		if tmpP.ExpirationTime.After(time.Now()) {
			msg := "temp profile already exists, what the hell"
			return msg, errors.New(msg)
		}

		repositories.TempProfileDelete(username)
		logger.Debug("deleted temp profile")
	}

	code := generateConfirmationCode()
	expirationTime := getExpirationDuration()
	if err := repositories.TempProfileCreate(models.TempProfile{
				ConfirmationCode: code,
				Username:         profileDb.Username,
				PasswordEncrypt:  "",
				ExpirationTime:   expirationTime,
			}); err != nil {
		return "error creating temp profile", err
	}

	if handler, check := recoveryByMapping[profileDb.RecoverBy]; check {
		err = handler.sendCode(profileDb, code)
		if err != nil {
			return "error sending code", err
		}
	} else {
		msg := "your recovery settings are invalid: you are screwed"
		return msg, errors.New(msg)
	}

	return "code sent", nil
}

func RecoveryConfirmCode(username string, password string, code string) (string, error) {
	// check code
	tempProfile, err := repositories.TempProfileGet(username)
	if err != nil {
		return "this user has not requested an account recovery, what the hell is going on", err
	}

	if tempProfile.ExpirationTime.Before(time.Now()) {
		repositories.TempProfileDelete(username)
		msg := "code expired, ask for a new code"
		return msg, errors.New(msg)
	}

	if code != tempProfile.ConfirmationCode {
		msg := "bad code"
		return msg, errors.New(msg)
	}

	profileDb, msg, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return msg, err
	}

	newPasswordEncrypted, err := models.HashPassword(password)
	if err != nil {
		return "error hashing password", err
	}

	profileDb.PasswordEncrypt = newPasswordEncrypted

	_, msg, err = repositories.ProfileUpdate(profileDb)
	if err != nil {
		return msg, err
	}

	if _, err := repositories.TempProfileDelete(username); err != nil {
		return "error deleting temp profile", err
	}

	return "profile recovered", nil
}

func RecoveryManagerTest(username string, password string) (string, error) {
	var err error
	p, msg, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return msg, err
	}

	if msg, err := checkUsernamePassword(p, password); err != nil {
		return msg, err
	}

	if p.RecoverBy == "" || p.Contact == "" {
		msg := "recovery settings not set"
		return msg, errors.New(msg)
	}

	if handler, check := recoveryByMapping[p.RecoverBy]; check {
		err = handler.test(p)
	} else {
		return "bad recovery settings", errors.New(
			fmt.Sprintf("bad recovery settings for %s: %s", p.Username, p.RecoverBy))
	}

	if err != nil {
		return "error sending test", err
	}


	return "test sent", nil
}


