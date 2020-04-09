package managers

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
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

func RecoverySendCode(username string) {

}

func RecoveryConfirmCode(username string, code string) {

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


